package pocketsmith

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// pageJob is a unit of work sent to a worker: a page number to fetch and the
// pre-built senderRequest for that page.
type pageJob struct {
	pageNum int
	sr      senderRequest
}

// pageResult holds the ordered result of a single page fetch.
type pageResult[T any] struct {
	pageNum int
	items   []T
	err     error
}

// fetchAllPages fetches all pages of a paginated API endpoint using a fixed
// worker pool.
//
// The algorithm:
//  1. Page 1 is fetched first (by the calling goroutine) to discover the total
//     record count via the "Total" and "Per-Page" response headers.
//  2. If only one page exists the result is returned immediately — no workers
//     are started.
//  3. For multi-page responses, `workers` goroutines are started. Each worker
//     reads page jobs from a shared jobs channel, fetches the page from the
//     API, and writes the result (in its own pre-allocated slot) to a results
//     slice. The same worker handles fetch and store for each job, keeping
//     data ownership simple and lock-free on the results slice.
//  4. Results are assembled in page order before returning.
//
// Fallback: if the "Total" or "Per-Page" headers are absent the single page-1
// result is returned as-is (the endpoint does not paginate).
//
// The `sr` argument must describe the page-1 request (queries already set).
// Workers inject `page=N` into a clone of that request for each subsequent page.
//
// Note: Go does not allow generic methods on a named type, so this is a
// package-level generic function that accepts the client explicitly.
func fetchAllPages[T any](
	ctx context.Context,
	c *Client,
	sr senderRequest,
	workers int,
) ([]T, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "fetchAllPages")
	defer span.End()

	// --- page 1: fetch sequentially to learn pagination metadata ---
	var page1 []T
	resp, err := c.sender(newCtx, sr, &page1)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to fetch page 1: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// read pagination metadata from response headers.
	total, perPage, ok := parsePaginationHeaders(resp.Header)
	if !ok || perPage == 0 {
		// endpoint does not return pagination headers — return page 1 as-is.
		span.SetAttributes(attribute.Bool("paginated", false))
		return page1, nil
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	span.SetAttributes(
		attribute.Bool("paginated", true),
		attribute.Int("total_records", total),
		attribute.Int("per_page", perPage),
		attribute.Int("total_pages", totalPages),
		attribute.Int("workers", workers),
	)

	if totalPages <= 1 {
		// everything fit on page 1.
		return page1, nil
	}

	remaining := totalPages - 1 // pages 2..N

	// cap workers to the number of remaining pages — no point spinning up more
	// workers than there is work.
	if workers > remaining {
		workers = remaining
	}

	// pre-allocate one result slot per remaining page so workers write to
	// their own index without any synchronisation on the slice.
	results := make([]pageResult[T], remaining)

	// jobs channel: buffered so we can enqueue all jobs before workers start
	// consuming, avoiding a potential deadlock if workers == 0.
	jobs := make(chan pageJob, remaining)
	for i := 0; i < remaining; i++ {
		pageNum := i + 2 // pages are 1-indexed; we already fetched page 1.
		jobs <- pageJob{
			pageNum: pageNum,
			sr:      cloneSenderRequest(sr, pageNum),
		}
	}
	close(jobs) // no more jobs; workers will drain and exit.

	// start the fixed worker pool.
	// A cancellable context is derived from newCtx so that if any worker
	// encounters an error it can signal the others to stop fetching, avoiding
	// wasted API calls for data that will be discarded.
	workerCtx, cancelWorkers := context.WithCancel(newCtx)
	defer cancelWorkers()

	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				// bail early if another worker already failed.
				if workerCtx.Err() != nil {
					// mark the slot so the assembly loop doesn't treat it as
					// a successful empty page — it will be skipped because the
					// first real error already caused an early return.
					idx := job.pageNum - 2
					results[idx] = pageResult[T]{
						pageNum: job.pageNum,
						err:     workerCtx.Err(),
					}
					continue
				}
				var batch []T
				_, err := c.sender(workerCtx, job.sr, &batch)
				// each job maps to a fixed index: pageNum 2 → idx 0, 3 → idx 1, …
				idx := job.pageNum - 2
				results[idx] = pageResult[T]{
					pageNum: job.pageNum,
					items:   batch,
					err:     err,
				}
				if err != nil {
					// cancel remaining workers — no point fetching more pages.
					cancelWorkers()
				}
			}
		}()
	}
	wg.Wait()

	// assemble: page 1 first, then pages 2..N in order.
	all := make([]T, 0, total)
	all = append(all, page1...)
	for _, r := range results {
		if r.err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to fetch page %d: %v", r.pageNum, r.err))
			span.RecordError(r.err)
			return nil, fmt.Errorf("page %d: %w", r.pageNum, r.err)
		}
		all = append(all, r.items...)
	}

	return all, nil
}

// parsePaginationHeaders reads the "Total" and "Per-Page" headers that the
// PocketSmith API includes on all paginated responses.
// Returns (0, 0, false) if either header is absent or not a valid integer.
func parsePaginationHeaders(h http.Header) (total, perPage int, ok bool) {
	totalStr := h.Get("Total")
	perPageStr := h.Get("Per-Page")
	if totalStr == "" || perPageStr == "" {
		return 0, 0, false
	}
	t, err1 := strconv.Atoi(totalStr)
	pp, err2 := strconv.Atoi(perPageStr)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return t, pp, true
}

// cloneSenderRequest returns a copy of sr with the `page` query parameter set
// to pageNum. All other query parameters from the original request are
// preserved, including any multi-value keys.
func cloneSenderRequest(sr senderRequest, pageNum int) senderRequest {
	clone := senderRequest{
		method: sr.method,
		path:   sr.path,
		body:   sr.body,
	}

	// clone url.Values directly — copying the full []string slice for each key
	// to preserve multi-value params and avoid the round-trip through
	// map[string]string that would silently drop extra values.
	newQ := make(url.Values, len(sr.queries)+1)
	for k, vals := range sr.queries {
		cp := make([]string, len(vals))
		copy(cp, vals)
		newQ[k] = cp
	}
	// inject (or override) the page number.
	newQ.Set("page", strconv.Itoa(pageNum))
	clone.queries = newQ

	return clone
}

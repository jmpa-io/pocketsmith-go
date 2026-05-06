package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"testing"
)

// item is a minimal struct used as the generic type T in fetchAllPages tests.
type item struct {
	ID int `json:"id"`
}

// makeTransactionPage returns an http.Response that looks like a PocketSmith
// paginated response for the given page of items, with correct Total and
// Per-Page headers set.
func makeTransactionPage(items []item, total, perPage int) *http.Response {
	body, _ := json.Marshal(items)
	h := make(http.Header)
	h.Set("Total", strconv.Itoa(total))
	h.Set("Per-Page", strconv.Itoa(perPage))
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     h,
	}
}

// makeEmptyPage returns an http.Response with no pagination headers (simulates
// a non-paginated endpoint).
func makeEmptyPage(items []item) *http.Response {
	body, _ := json.Marshal(items)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}
}

// makeErrorPage returns a 500 response with an error body.
func makeErrorPage() *http.Response {
	body, _ := json.Marshal(apiErrorResponse{Error: "internal server error"})
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}
}

// newTestClient builds a *Client wired to the given RoundTripper, bypassing
// the New() constructor's live GetAuthedUser call.
func newTestClient(t *testing.T, rt http.RoundTripper) *Client {
	t.Helper()
	c := &Client{
		tracerName: "pocketsmith-go-test",
		endpoint:   "https://api.pocketsmith.com/v2",
		httpClient: &http.Client{Transport: rt},
		headers:    make(http.Header),
		logger:     logger, // reuse test logger from client_test.go
	}
	c.headers.Set("X-Developer-Key", "test-token")
	return c
}

func Test_fetchAllPages(t *testing.T) {
	tests := map[string]struct {
		// pageData maps page number → items on that page.
		pageData map[int][]item
		total    int
		perPage  int
		// noPaginationHeaders causes the mock to omit Total/Per-Page headers.
		noPaginationHeaders bool
		// errorOnPage causes that page to return an HTTP 500.
		errorOnPage int
		workers     int
		wantIDs     []int  // expected IDs in order
		wantErr     string // substring expected in error, empty means no error
	}{
		"single page — all results fit on page 1": {
			pageData: map[int][]item{
				1: {{ID: 1}, {ID: 2}, {ID: 3}},
			},
			total:   3,
			perPage: 1000,
			workers: 10,
			wantIDs: []int{1, 2, 3},
		},
		"two pages — second page fetched by worker": {
			pageData: map[int][]item{
				1: {{ID: 1}, {ID: 2}},
				2: {{ID: 3}, {ID: 4}},
			},
			total:   4,
			perPage: 2,
			workers: 10,
			wantIDs: []int{1, 2, 3, 4},
		},
		"five pages — workers fewer than pages": {
			pageData: map[int][]item{
				1: {{ID: 1}},
				2: {{ID: 2}},
				3: {{ID: 3}},
				4: {{ID: 4}},
				5: {{ID: 5}},
			},
			total:   5,
			perPage: 1,
			workers: 2, // only 2 workers for 4 remaining pages
			wantIDs: []int{1, 2, 3, 4, 5},
		},
		"workers capped to remaining pages": {
			pageData: map[int][]item{
				1: {{ID: 10}},
				2: {{ID: 20}},
			},
			total:   2,
			perPage: 1,
			workers: 100, // more workers than pages — should be capped to 1
			wantIDs: []int{10, 20},
		},
		"no pagination headers — returns page 1 only": {
			pageData: map[int][]item{
				1: {{ID: 99}},
			},
			noPaginationHeaders: true,
			workers:             10,
			wantIDs:             []int{99},
		},
		"error on page 2 — returns error": {
			pageData: map[int][]item{
				1: {{ID: 1}, {ID: 2}},
				2: nil, // will be replaced with a 500
			},
			total:       4,
			perPage:     2,
			errorOnPage: 2,
			workers:     10,
			wantErr:     "page 2",
		},
		"100 pages — 3 workers stress test": func() struct {
			pageData            map[int][]item
			total               int
			perPage             int
			noPaginationHeaders bool
			errorOnPage         int
			workers             int
			wantIDs             []int
			wantErr             string
		} {
			const n = 100
			pageData := make(map[int][]item, n)
			wantIDs := make([]int, n)
			for i := 1; i <= n; i++ {
				pageData[i] = []item{{ID: i}}
				wantIDs[i-1] = i
			}
			return struct {
				pageData            map[int][]item
				total               int
				perPage             int
				noPaginationHeaders bool
				errorOnPage         int
				workers             int
				wantIDs             []int
				wantErr             string
			}{
				pageData: pageData,
				total:    n,
				perPage:  1,
				workers:  3,
				wantIDs:  wantIDs,
			}
		}(),
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// track concurrent calls to verify worker pool behaviour.
			var mu sync.Mutex
			maxConcurrent := 0
			inFlight := 0

			rt := &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					// parse the page number from the query string.
					pageStr := req.URL.Query().Get("page")
					pageNum := 1
					if pageStr != "" {
						pageNum, _ = strconv.Atoi(pageStr)
					}

					// track concurrency.
					mu.Lock()
					inFlight++
					if inFlight > maxConcurrent {
						maxConcurrent = inFlight
					}
					mu.Unlock()
					defer func() {
						mu.Lock()
						inFlight--
						mu.Unlock()
					}()

					if tt.errorOnPage != 0 && pageNum == tt.errorOnPage {
						return makeErrorPage()
					}
					items := tt.pageData[pageNum]
					if tt.noPaginationHeaders {
						return makeEmptyPage(items)
					}
					return makeTransactionPage(items, tt.total, tt.perPage)
				},
			}

			c := newTestClient(t, rt)
			sr := senderRequest{
				method:  http.MethodGet,
				path:    "/transaction_accounts/123/transactions",
				queries: setupQueries(nil),
			}

			got, err := fetchAllPages[item](context.Background(), c, sr, tt.workers)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("fetchAllPages() expected error containing %q, got nil", tt.wantErr)
				}
				if !containsStr(err.Error(), tt.wantErr) {
					t.Fatalf("fetchAllPages() error = %q, want substring %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("fetchAllPages() unexpected error: %v", err)
			}

			// extract IDs and sort for comparison (concurrent fetches may
			// arrive in any order before assembly, but assembly is ordered).
			gotIDs := make([]int, len(got))
			for i, it := range got {
				gotIDs[i] = it.ID
			}

			wantSorted := make([]int, len(tt.wantIDs))
			copy(wantSorted, tt.wantIDs)
			sort.Ints(wantSorted)
			sort.Ints(gotIDs)

			if !intsEqual(gotIDs, wantSorted) {
				t.Errorf("fetchAllPages() IDs = %v, want %v", gotIDs, wantSorted)
			}

			// verify workers were actually bounded.
			if !tt.noPaginationHeaders && tt.wantErr == "" && len(tt.pageData) > 1 {
				if maxConcurrent > tt.workers+1 { // +1 for the page-1 call
					t.Errorf("max concurrent calls = %d, want <= workers (%d) + 1", maxConcurrent, tt.workers)
				}
			}
		})
	}
}

func Test_parsePaginationHeaders(t *testing.T) {
	tests := map[string]struct {
		header      http.Header
		wantTotal   int
		wantPerPage int
		wantOk      bool
	}{
		"both headers present": {
			header:      http.Header{"Total": []string{"432"}, "Per-Page": []string{"100"}},
			wantTotal:   432,
			wantPerPage: 100,
			wantOk:      true,
		},
		"Total header missing": {
			header: http.Header{"Per-Page": []string{"100"}},
			wantOk: false,
		},
		"Per-Page header missing": {
			header: http.Header{"Total": []string{"432"}},
			wantOk: false,
		},
		"both headers missing": {
			header: make(http.Header),
			wantOk: false,
		},
		"Total not a number": {
			header: http.Header{"Total": []string{"abc"}, "Per-Page": []string{"100"}},
			wantOk: false,
		},
		"Per-Page not a number": {
			header: http.Header{"Total": []string{"432"}, "Per-Page": []string{"xyz"}},
			wantOk: false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			total, perPage, ok := parsePaginationHeaders(tt.header)
			if ok != tt.wantOk {
				t.Errorf("parsePaginationHeaders() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok {
				if total != tt.wantTotal {
					t.Errorf("parsePaginationHeaders() total = %d, want %d", total, tt.wantTotal)
				}
				if perPage != tt.wantPerPage {
					t.Errorf("parsePaginationHeaders() perPage = %d, want %d", perPage, tt.wantPerPage)
				}
			}
		})
	}
}

func Test_cloneSenderRequest(t *testing.T) {
	tests := map[string]struct {
		sr      senderRequest
		pageNum int
		wantPage string
	}{
		"injects page into empty queries": {
			sr: senderRequest{
				method:  http.MethodGet,
				path:    "/foo",
				queries: setupQueries(nil),
			},
			pageNum:  3,
			wantPage: "3",
		},
		"overrides existing page param": {
			sr: senderRequest{
				method: http.MethodGet,
				path:   "/foo",
				queries: setupQueries(&map[string]string{
					"page":       "1",
					"start_date": "2024-01-01",
				}),
			},
			pageNum:  7,
			wantPage: "7",
		},
		"preserves other query params": {
			sr: senderRequest{
				method: http.MethodGet,
				path:   "/bar",
				queries: setupQueries(&map[string]string{
					"start_date": "2024-01-01",
					"end_date":   "2024-12-31",
				}),
			},
			pageNum:  2,
			wantPage: "2",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			clone := cloneSenderRequest(tt.sr, tt.pageNum)

			if clone.method != tt.sr.method {
				t.Errorf("cloneSenderRequest() method = %q, want %q", clone.method, tt.sr.method)
			}
			if clone.path != tt.sr.path {
				t.Errorf("cloneSenderRequest() path = %q, want %q", clone.path, tt.sr.path)
			}
			if got := clone.queries.Get("page"); got != tt.wantPage {
				t.Errorf("cloneSenderRequest() page = %q, want %q", got, tt.wantPage)
			}
			// original request must not be mutated.
			if orig := tt.sr.queries.Get("page"); tt.sr.queries != nil && orig == tt.wantPage && tt.pageNum != 1 {
				// only flag if original had a different page and got overwritten.
				if origPage, _ := strconv.Atoi(orig); origPage == tt.pageNum && orig != "1" {
					t.Errorf("cloneSenderRequest() mutated original sr.queries")
				}
			}
			// verify other params from the original are preserved in the clone.
			for k, vals := range tt.sr.queries {
				if k == "page" {
					continue
				}
				if got := clone.queries.Get(k); len(vals) > 0 && got != vals[0] {
					t.Errorf("cloneSenderRequest() lost query param %q: got %q, want %q", k, got, vals[0])
				}
			}
		})
	}
}

// --- helpers ---

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

func intsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

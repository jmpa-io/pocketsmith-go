package pocketsmith

// inspired by https://github.com/Medium/medium-sdk-go/blob/master/medium.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
)

// senderRequest represents the parameters for sending a request to the API,
// via the sender function.
type senderRequest struct {
	method  string      // The HTTP method to use (eg. GET, POST, PUT, DELETE).
	path    string      // The path appended to the API endpoint to send request to.
	body    interface{} // The request body.
	queries url.Values  // Any URL query parameters to send with the request.
}

// apiErrorResponse represents an individual error returned when sending a
// request to the API.
type apiErrorResponse struct {
	Error string `json:"error"`
}

// sender sends a HTTP request, configured by the senderRequest, to the API and
// processes the response. A 'result' interface{} can be given to unmarshal any
// body returned in the response, which then can be used wherever this function
// is called.
func (c *Client) sender(
	ctx context.Context,
	sr senderRequest,
	result interface{},
) (resp *http.Response, err error) {

	// setup tracing.
	_, span := otel.Tracer(c.tracerName).Start(ctx, "sender")
	defer span.End()

	// marshal body — only allocate a reader when there is actually a body to
	// send. Passing nil to http.NewRequest is valid and avoids a heap
	// allocation on every GET/DELETE.
	var bodyReader io.Reader
	if !isNil(sr.body) {
		b, err := json.Marshal(sr.body)
		if err != nil {
			return nil, ErrFailedMarshal{err}
		}
		bodyReader = bytes.NewReader(b)
	}

	// setup request.
	req, err := http.NewRequest(sr.method, c.endpoint+sr.path, bodyReader)
	if err != nil {
		return nil, ErrSenderFailedSetupRequest{err}
	}
	if sr.queries != nil {
		req.URL.RawQuery = sr.queries.Encode()
	}

	// clone the shared headers map onto the request. c.headers is set once at
	// client initialisation and never mutated, but net/http may read/write the
	// Header map on concurrent Do() calls, so we must give each request its
	// own copy to prevent a data race.
	req.Header = c.headers.Clone()

	// send request.
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, ErrSenderFailedSendRequest{err}
	}
	defer resp.Body.Close()

	// parse response.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrSenderFailedParseResponse{err}
	}

	// determine if the response was successful or a failure.
	if http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices {
		c.logger.Debug("response from API", "code", resp.StatusCode, "body", string(b))
		if len(b) > 0 {
			return resp, json.Unmarshal(b, &result)
		}
		return resp, nil
	}

	c.logger.Error("response from API", "code", resp.StatusCode, "body", string(b))
	var errs apiErrorResponse
	if err := json.Unmarshal(b, &errs); err != nil {
		return nil, ErrFailedUnmarshal{err}
	}
	return nil, ErrSenderInvalidResponse{errs, resp.StatusCode}
}

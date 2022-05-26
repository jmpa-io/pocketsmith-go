package pocketsmith

// inspired by https://github.com/Medium/medium-sdk-go/blob/master/medium.go.

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log/level"
)

// clientRequest simplifies sending a given request to the API by the sender.
type clientRequest struct {
	method string
	path   string
	data   interface{}
}

// envelope defines errors returned from the API.
type envelope struct {
	Error string `json:"error"`
}

// sender sends the given request to the API.
func (c *Client) sender(cr clientRequest, result interface{}) (*http.Response, error) {

	// marshal body.
	var body []byte
	if !isNil(cr.data) {
		b, err := json.Marshal(cr.data)
		if err != nil {
			return nil, ErrFailedMarshal{err}
		}
		body = b
	}

	// setup request.
	req, err := http.NewRequest(cr.method, c.endpoint+cr.path, bytes.NewReader(body))
	if err != nil {
		return nil, ErrSenderFailedSetupRequest{err}
	}

	// add headers to request.
	req.Header = c.headers

	// send request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrSenderFailedSendRequest{err}
	}
	defer resp.Body.Close()

	// parse response.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrSenderFailedParseResponse{err}
	}
	_ = level.Debug(c.logger).Log("msg", "response from API",
		"code", resp.StatusCode, "body", string(body))

	// was this a valid request to the API?
	if http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices {
		if len(b) > 0 {
			return resp, json.Unmarshal(b, &result)
		}
		return resp, nil
	}

	// since we have an unexpected invalid response, return a generic response.
	var env envelope
	if err := json.Unmarshal(b, &env); err != nil {
		return nil, ErrFailedUnmarshal{err}
	}
	return nil, ErrSenderInvalidResponse{env.Error, resp.StatusCode}
}

package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// A reader to break reading from a *http.Request body.
// from: https://stackoverflow.com/questions/45126312/how-do-i-test-an-error-on-reading-from-a-request-body
type brokenReader struct{}

func (br *brokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed reading")
}

func (br *brokenReader) Close() error {
	return fmt.Errorf("failed closing")
}

// An empty error used in testing for Err* errors.
var emptyErr = fmt.Errorf("")

func Test_sender(t *testing.T) {
	tests := map[string]struct {
		mock    *mockRoundTripper
		request senderRequest
		result  interface{}
		err     string
	}{
		// "catch json marshal error": {
		// 	request: senderRequest{
		// 		body: make(chan int), // can't marshal a channel.
		// 	},
		// 	err: ErrFailedMarshal{emptyErr}.Error(),
		// },
		// "catch failed setup request": {
		// 	mock: &mockRoundTripper{
		// 		MockFunc: func(req *http.Request) *http.Response {
		// 			return &http.Response{}
		// 		},
		// 	},
		// 	request: senderRequest{
		// 		method: "https://", // returns an error on invalid method.
		// 	},
		// 	err: ErrSenderFailedSetupRequest{emptyErr}.Error(),
		// },
		// "catch failed send request": {
		// 	mock: &mockRoundTripper{
		// 		MockFunc: func(req *http.Request) *http.Response {
		// 			return &http.Response{
		// 				StatusCode: http.StatusUnauthorized,
		// 			}
		// 		},
		// 	},
		// 	request: senderRequest{
		// 		method: "error", // not a real http method, this is just to test.
		// 	},
		// 	err: ErrSenderFailedSendRequest{emptyErr}.Error(),
		// },
		// "catch failed parse response": {
		// 	mock: &mockRoundTripper{
		// 		MockFunc: func(req *http.Request) *http.Response {
		// 			return &http.Response{
		// 				Body: &brokenReader{},
		// 			}
		// 		},
		// 	},
		// 	err: ErrSenderFailedParseResponse{emptyErr}.Error(),
		// },
		// "catch json unmarshal error": {
		// 	mock: &mockRoundTripper{
		// 		MockFunc: func(req *http.Request) *http.Response {
		// 			return &http.Response{}
		// 		},
		// 	},
		// 	err: ErrFailedUnmarshal{emptyErr}.Error(),
		// },
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// setup client with mock.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(&http.Client{
				Transport: tt.mock,
			}),
		)

		// run tests.
		t.Run(name, func(t *testing.T) {
			_, err := c.sender(ctx, tt.request, &tt.result)

			// any errors?
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"sender() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("sender() returned an error;\nerror=%v\n", err)
				return
			}

		})
	}
}

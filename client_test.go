package pocketsmith

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var (

	// test logger.
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))

	// test httpClient.
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	// test authed user mock.
	authedUserMock = &mockRoundTripper{
		MockFunc: func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(""))),
				Header:     make(http.Header),
			}
		},
	}
)

func Test_New(t *testing.T) {
	tests := map[string]struct {
		mock    *mockRoundTripper
		token   string
		options []Option
		want    *Client
		err     string
	}{
		"default": {
			mock:  authedUserMock,
			token: "xxxx",
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
			},
		},
		"no token": {
			want: &Client{},
			err:  "the provided token is empty",
		},
		"with log level (debug)": {
			mock:    authedUserMock,
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelDebug)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   -4,
			},
		},
		"with log level (info)": {
			mock:    authedUserMock,
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelInfo)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   0,
			},
		},
		"with log level (warn)": {
			mock:    authedUserMock,
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelWarn)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   4,
			},
		},
		"with log level (error)": {
			mock:    authedUserMock,
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelError)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   8,
			},
		},
		"with logger": {
			mock:    authedUserMock,
			token:   "xxxx",
			options: []Option{WithLogger(logger)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     logger,
			},
		},
		"with http client": {
			mock:    authedUserMock,
			token:   "xxxx",
			options: []Option{WithHttpClient(httpClient)},
			want: &Client{
				httpClient: httpClient,
				logger:     slog.Default(),
			},
		},
	}
	for name, tt := range tests {

		// setup headers.
		headers := make(http.Header)
		if tt.token != "" {
			headers.Add("Authorization", "Bearer "+tt.token)
		}

		// add mock to client.
		if tt.mock != nil {
			tt.options = append(tt.options, WithHttpClient(
				&http.Client{
					Transport: tt.mock,
				},
			))
		}

		// run tests.
		t.Run(name, func(t *testing.T) {
			got, err := New(context.Background(), tt.token, tt.options...)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("New() returned an unexpected error; want=%v, got=%v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("New() returned an error; error=%v", err)
				return
			}
			switch {
			case
				got.headers.Get("Authorization") != headers.Get("Authorization"),
				got.logLevel != tt.want.logLevel,
				(got.logger != slog.Default() && tt.want.logger != slog.Default()) && got.logger != tt.want.logger,
				got.httpClient != tt.want.httpClient:
				t.Errorf(
					"New() returned unexpected configuration; want=%+v, got=%+v\n",
					tt.want,
					got,
				)
				return
			}
		})
	}
}

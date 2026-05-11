package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"testing"
)

// userJSON is the minimal User JSON returned for the GetAuthedUser init call.
// customTime fields (forecast_start_date, forecast_end_date) require YYYY-MM-DD
// format — using a raw map avoids the zero-time marshalling issue.
var userJSON, _ = json.Marshal(map[string]any{
	"id":                   1,
	"login":                "testuser",
	"forecast_start_date":  "2026-01-01",
	"forecast_end_date":    "2026-12-31",
})

// newMethodClient creates a pocketsmith Client whose HTTP transport answers the
// first request (the GetAuthedUser init call from New()) with a valid user
// response, then delegates all subsequent requests to fn.
func newMethodClient(t *testing.T, fn func(*http.Request) *http.Response) *Client {
	t.Helper()
	var count atomic.Int32
	rt := &mockRoundTripper{MockFunc: func(req *http.Request) *http.Response {
		if count.Add(1) == 1 {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(userJSON)),
				Header:     make(http.Header),
			}
		}
		return fn(req)
	}}
	c, err := New(context.Background(), "xxxx",
		WithHttpClient(&http.Client{Transport: rt}),
		WithLogLevel(slog.LevelError),
	)
	if err != nil {
		t.Fatalf("newMethodClient: New() error: %v", err)
	}
	return c
}

func Test_ListAccounts(t *testing.T) {
	tests := map[string]struct {
		mockFn func(*http.Request) *http.Response
		want   Accounts
		err    string
	}{
		"success — returns two accounts": {
			mockFn: func(req *http.Request) *http.Response {
				want := Accounts{
					{ID: 1, Title: "Everyday Account", Type: "bank"},
					{ID: 2, Title: "Savings Account", Type: "bank"},
				}
				b, _ := json.Marshal(want)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: Accounts{
				{ID: 1, Title: "Everyday Account", Type: "bank"},
				{ID: 2, Title: "Savings Account", Type: "bank"},
			},
		},
		"success — empty list": {
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(Accounts{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: Accounts{},
		},
		"api error — returns error": {
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(apiErrorResponse{Error: "unauthorized"})
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "unauthorized",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			c := newMethodClient(t, tt.mockFn)

			got, err := c.ListAccounts(context.Background())
			if tt.err != "" {
				if err == nil {
					t.Fatalf("ListAccounts() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("ListAccounts() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ListAccounts() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ListAccounts() returned %d accounts, want %d", len(got), len(tt.want))
			}
			for i, a := range got {
				if a.ID != tt.want[i].ID {
					t.Errorf("ListAccounts()[%d].ID = %d, want %d", i, a.ID, tt.want[i].ID)
				}
				if a.Title != tt.want[i].Title {
					t.Errorf("ListAccounts()[%d].Title = %q, want %q", i, a.Title, tt.want[i].Title)
				}
			}
		})
	}
}

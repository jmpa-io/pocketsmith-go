package pocketsmith

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

var tdCategories = newTestdata("categories")

func Test_ListCategories(t *testing.T) {
	tests := map[string]struct {
		mockFn func(*http.Request) *http.Response
		want   Categories
		err    string
	}{
		"success — returns two categories": {
			mockFn: func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(tdCategories.content)),
					Header:     make(http.Header),
				}
			},
			want: Categories{
				{ID: 5, Title: "Groceries"},
				{ID: 6, Title: "Fuel"},
			},
		},
		"success — empty list": {
			mockFn: func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("[]"))),
					Header:     make(http.Header),
				}
			},
			want: Categories{},
		},
		"api error — unauthorized": {
			mockFn: func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"error":"unauthorized"}`))),
					Header:     make(http.Header),
				}
			},
			err: "unauthorized",
		},
		"success — category with children": {
			mockFn: func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`[{"id":5,"title":"Groceries","colour":"#4caf50","children":[{"id":50,"title":"Fresh Produce","colour":"#8bc34a"}]}]`))),
					Header:     make(http.Header),
				}
			},
			want: Categories{
				{
					ID:    5,
					Title: "Groceries",
					Children: []*Category{
						{ID: 50, Title: "Fresh Produce"},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			c := newMethodClient(t, tt.mockFn)

			got, err := c.ListCategories(context.Background())
			if tt.err != "" {
				if err == nil {
					t.Fatalf("ListCategories() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("ListCategories() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ListCategories() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ListCategories() returned %d categories, want %d", len(got), len(tt.want))
			}
			for i, cat := range got {
				if cat.ID != tt.want[i].ID {
					t.Errorf("ListCategories()[%d].ID = %d, want %d", i, cat.ID, tt.want[i].ID)
				}
				if cat.Title != tt.want[i].Title {
					t.Errorf("ListCategories()[%d].Title = %q, want %q", i, cat.Title, tt.want[i].Title)
				}
			}
		})
	}
}

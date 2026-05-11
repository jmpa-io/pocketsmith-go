package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func Test_ListCategories(t *testing.T) {
	tests := map[string]struct {
		mockFn func(*http.Request) *http.Response
		want   Categories
		err    string
	}{
		"success — returns two categories": {
			mockFn: func(req *http.Request) *http.Response {
				want := Categories{
					{ID: 5, Title: "Groceries", Colour: "#4caf50"},
					{ID: 6, Title: "Fuel", Colour: "#ff5722"},
				}
				b, _ := json.Marshal(want)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: Categories{
				{ID: 5, Title: "Groceries", Colour: "#4caf50"},
				{ID: 6, Title: "Fuel", Colour: "#ff5722"},
			},
		},
		"success — empty list": {
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(Categories{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: Categories{},
		},
		"api error — unauthorized": {
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
		"success — category with children": {
			mockFn: func(req *http.Request) *http.Response {
				child := &Category{ID: 50, Title: "Fresh Produce", Colour: "#8bc34a"}
				want := Categories{
					{ID: 5, Title: "Groceries", Colour: "#4caf50", Children: []*Category{child}},
				}
				b, _ := json.Marshal(want)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
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

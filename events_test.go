package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func Test_GetBudget(t *testing.T) {
	tests := map[string]struct {
		mockFn func(*http.Request) *http.Response
		want   []BudgetItem
		err    string
	}{
		"success — returns budget items": {
			mockFn: func(req *http.Request) *http.Response {
				items := []BudgetItem{
					{
						Category:   Category{ID: 5, Title: "Groceries"},
						IsTransfer: false,
						Expense: &BudgetEntry{
							StartDate:           "2026-05-01",
							EndDate:             "2026-05-31",
							CurrencyCode:        "AUD",
							TotalActualAmount:   -55.30,
							TotalForecastAmount: -200.00,
						},
					},
					{
						Category:   Category{ID: 6, Title: "Fuel"},
						IsTransfer: false,
						Expense: &BudgetEntry{
							StartDate:           "2026-05-01",
							EndDate:             "2026-05-31",
							CurrencyCode:        "AUD",
							TotalActualAmount:   -80.00,
							TotalForecastAmount: -150.00,
						},
					},
				}
				b, _ := json.Marshal(items)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: []BudgetItem{
				{Category: Category{ID: 5, Title: "Groceries"}},
				{Category: Category{ID: 6, Title: "Fuel"}},
			},
		},
		"success — empty budget": {
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal([]BudgetItem{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: []BudgetItem{},
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

			got, err := c.GetBudget(context.Background())
			if tt.err != "" {
				if err == nil {
					t.Fatalf("GetBudget() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("GetBudget() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetBudget() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("GetBudget() returned %d items, want %d", len(got), len(tt.want))
			}
			for i, item := range got {
				if item.Category.ID != tt.want[i].Category.ID {
					t.Errorf("GetBudget()[%d].Category.ID = %d, want %d", i, item.Category.ID, tt.want[i].Category.ID)
				}
			}
		})
	}
}

func Test_ListEvents(t *testing.T) {
	tests := map[string]struct {
		mockFn  func(*http.Request) *http.Response
		options *ListEventsOptions
		want    Events
		err     string
	}{
		"success — returns two events": {
			options: &ListEventsOptions{StartDate: "2026-05-01", EndDate: "2026-05-31"},
			mockFn: func(req *http.Request) *http.Response {
				events := Events{
					{
						ID:           "evt-001",
						Amount:       -200.00,
						CurrencyCode: "AUD",
						RepeatType:   "monthly",
					},
					{
						ID:           "evt-002",
						Amount:       -150.00,
						CurrencyCode: "AUD",
						RepeatType:   "monthly",
					},
				}
				b, _ := json.Marshal(events)
				h := make(http.Header)
				h.Set("Total", strconv.Itoa(2))
				h.Set("Per-Page", strconv.Itoa(1000))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     h,
				}
			},
			want: Events{
				{ID: "evt-001", Amount: -200.00},
				{ID: "evt-002", Amount: -150.00},
			},
		},
		"success — empty events": {
			options: &ListEventsOptions{StartDate: "2026-05-01", EndDate: "2026-05-31"},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(Events{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: Events{},
		},
		"empty start date — validation error": {
			options: &ListEventsOptions{StartDate: "", EndDate: "2026-05-31"},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(Events{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "StartDate",
		},
		"empty end date — validation error": {
			options: &ListEventsOptions{StartDate: "2026-05-01", EndDate: ""},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(Events{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "EndDate",
		},
		"api error — returns error": {
			options: &ListEventsOptions{StartDate: "2026-05-01", EndDate: "2026-05-31"},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(apiErrorResponse{Error: "forbidden"})
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "forbidden",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			c := newMethodClient(t, tt.mockFn)

			got, err := c.ListEvents(context.Background(), tt.options)
			if tt.err != "" {
				if err == nil {
					t.Fatalf("ListEvents() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("ListEvents() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ListEvents() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ListEvents() returned %d events, want %d", len(got), len(tt.want))
			}
			for i, ev := range got {
				if ev.ID != tt.want[i].ID {
					t.Errorf("ListEvents()[%d].ID = %q, want %q", i, ev.ID, tt.want[i].ID)
				}
			}
		})
	}
}

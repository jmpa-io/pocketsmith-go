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

func Test_ListTransactionAccounts(t *testing.T) {
	tests := map[string]struct {
		mockFn func(*http.Request) *http.Response
		want   TransactionAccounts
		err    string
	}{
		"success — returns two transaction accounts": {
			mockFn: func(req *http.Request) *http.Response {
				want := TransactionAccounts{
					{ID: 10, Name: "Everyday", Type: "checking"},
					{ID: 11, Name: "Savings", Type: "savings"},
				}
				b, _ := json.Marshal(want)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: TransactionAccounts{
				{ID: 10, Name: "Everyday", Type: "checking"},
				{ID: 11, Name: "Savings", Type: "savings"},
			},
		},
		"success — empty list": {
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(TransactionAccounts{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: TransactionAccounts{},
		},
		"api error — returns error": {
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

			got, err := c.ListTransactionAccounts(context.Background())
			if tt.err != "" {
				if err == nil {
					t.Fatalf("ListTransactionAccounts() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("ListTransactionAccounts() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ListTransactionAccounts() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ListTransactionAccounts() returned %d accounts, want %d", len(got), len(tt.want))
			}
			for i, a := range got {
				if a.ID != tt.want[i].ID {
					t.Errorf("ListTransactionAccounts()[%d].ID = %d, want %d", i, a.ID, tt.want[i].ID)
				}
				if a.Name != tt.want[i].Name {
					t.Errorf("ListTransactionAccounts()[%d].Name = %q, want %q", i, a.Name, tt.want[i].Name)
				}
			}
		})
	}
}

func Test_ListTransactionAccountTransactions(t *testing.T) {
	tests := map[string]struct {
		mockFn  func(*http.Request) *http.Response
		options *ListTransactionAccountTransactionsOptions
		want    []Transaction
		err     string
	}{
		"success — returns two transactions": {
			options: &ListTransactionAccountTransactionsOptions{
				TransactionAccountID: "10",
			},
			mockFn: func(req *http.Request) *http.Response {
				txns := []Transaction{
					{ID: 42, Payee: "Coles Supermarket", Amount: -55.30},
					{ID: 43, Payee: "BP Fuel Station", Amount: -80.00},
				}
				b, _ := json.Marshal(txns)
				h := make(http.Header)
				h.Set("Total", strconv.Itoa(2))
				h.Set("Per-Page", strconv.Itoa(1000))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     h,
				}
			},
			want: []Transaction{
				{ID: 42, Payee: "Coles Supermarket", Amount: -55.30},
				{ID: 43, Payee: "BP Fuel Station", Amount: -80.00},
			},
		},
		"success — empty result": {
			options: &ListTransactionAccountTransactionsOptions{
				TransactionAccountID: "10",
			},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal([]Transaction{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: []Transaction{},
		},
		"empty transaction account id — passes through (no validation)": {
			options: &ListTransactionAccountTransactionsOptions{
				TransactionAccountID: "",
			},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal([]Transaction{})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: []Transaction{},
		},
		"api error — returns error": {
			options: &ListTransactionAccountTransactionsOptions{
				TransactionAccountID: "10",
			},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(apiErrorResponse{Error: "internal server error"})
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "internal server error",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			c := newMethodClient(t, tt.mockFn)

			got, err := c.ListTransactionAccountTransactions(context.Background(), tt.options)
			if tt.err != "" {
				if err == nil {
					t.Fatalf("ListTransactionAccountTransactions() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("ListTransactionAccountTransactions() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ListTransactionAccountTransactions() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ListTransactionAccountTransactions() returned %d transactions, want %d", len(got), len(tt.want))
			}
			for i, tx := range got {
				if tx.ID != tt.want[i].ID {
					t.Errorf("ListTransactionAccountTransactions()[%d].ID = %d, want %d", i, tx.ID, tt.want[i].ID)
				}
				if tx.Payee != tt.want[i].Payee {
					t.Errorf("ListTransactionAccountTransactions()[%d].Payee = %q, want %q", i, tx.Payee, tt.want[i].Payee)
				}
			}
		})
	}
}

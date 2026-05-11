package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func Test_GetTransaction(t *testing.T) {
	tests := map[string]struct {
		mockFn  func(*http.Request) *http.Response
		options *GetTransactionOptions
		want    *Transaction
		err     string
	}{
		"success — returns transaction": {
			options: &GetTransactionOptions{TransactionID: 42},
			mockFn: func(req *http.Request) *http.Response {
				tx := Transaction{
					ID:     42,
					Payee:  "Coles Supermarket",
					Amount: -55.30,
					Date:   "2026-05-10",
				}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{
				ID:     42,
				Payee:  "Coles Supermarket",
				Amount: -55.30,
				Date:   "2026-05-10",
			},
		},
		"success — returns empty transaction for id zero (no validation)": {
			options: &GetTransactionOptions{TransactionID: 0},
			mockFn: func(req *http.Request) *http.Response {
				tx := Transaction{ID: 0}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{ID: 0},
		},
		"api error — not found": {
			options: &GetTransactionOptions{TransactionID: 9999},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(apiErrorResponse{Error: "not found"})
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "not found",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			c := newMethodClient(t, tt.mockFn)

			got, err := c.GetTransaction(context.Background(), tt.options)
			if tt.err != "" {
				if err == nil {
					t.Fatalf("GetTransaction() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("GetTransaction() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetTransaction() unexpected error: %v", err)
			}
			if got.ID != tt.want.ID {
				t.Errorf("GetTransaction().ID = %d, want %d", got.ID, tt.want.ID)
			}
			if got.Payee != tt.want.Payee {
				t.Errorf("GetTransaction().Payee = %q, want %q", got.Payee, tt.want.Payee)
			}
			if got.Amount != tt.want.Amount {
				t.Errorf("GetTransaction().Amount = %.2f, want %.2f", got.Amount, tt.want.Amount)
			}
		})
	}
}

func Test_UpdateTransaction(t *testing.T) {
	tests := map[string]struct {
		mockFn  func(*http.Request) *http.Response
		options *UpdateTransactionOptions
		want    *Transaction
		err     string
	}{
		"success — updates payee": {
			options: &UpdateTransactionOptions{
				TransactionID: 42,
				Payee:         "Coles Express",
			},
			mockFn: func(req *http.Request) *http.Response {
				tx := Transaction{
					ID:    42,
					Payee: "Coles Express",
					Date:  "2026-05-10",
				}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{
				ID:    42,
				Payee: "Coles Express",
			},
		},
		"success — updates labels": {
			options: &UpdateTransactionOptions{
				TransactionID: 42,
				Labels:        "groceries,weekly",
			},
			mockFn: func(req *http.Request) *http.Response {
				tx := Transaction{
					ID:     42,
					Payee:  "Coles Supermarket",
					Labels: []string{"groceries", "weekly"},
				}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{
				ID:     42,
				Payee:  "Coles Supermarket",
				Labels: []string{"groceries", "weekly"},
			},
		},
		"zero transaction id — passes through (no validation)": {
			options: &UpdateTransactionOptions{TransactionID: 0},
			mockFn: func(req *http.Request) *http.Response {
				tx := Transaction{ID: 0}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{ID: 0},
		},
		"api error — returns error": {
			options: &UpdateTransactionOptions{TransactionID: 42, Payee: "bad"},
			mockFn: func(req *http.Request) *http.Response {
				b, _ := json.Marshal(apiErrorResponse{Error: "unprocessable entity"})
				return &http.Response{
					StatusCode: http.StatusUnprocessableEntity,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			err: "unprocessable entity",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			c := newMethodClient(t, tt.mockFn)

			got, err := c.UpdateTransaction(context.Background(), tt.options)
			if tt.err != "" {
				if err == nil {
					t.Fatalf("UpdateTransaction() expected error containing %q, got nil", tt.err)
				}
				if !containsStr(err.Error(), tt.err) {
					t.Fatalf("UpdateTransaction() error = %q, want substring %q", err.Error(), tt.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateTransaction() unexpected error: %v", err)
			}
			if got.ID != tt.want.ID {
				t.Errorf("UpdateTransaction().ID = %d, want %d", got.ID, tt.want.ID)
			}
			if got.Payee != tt.want.Payee {
				t.Errorf("UpdateTransaction().Payee = %q, want %q", got.Payee, tt.want.Payee)
			}
		})
	}
}

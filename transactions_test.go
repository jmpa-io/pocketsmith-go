package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

var tdTransaction = newTestdata("transaction")

func Test_GetTransaction(t *testing.T) {
	tests := map[string]struct {
		mockFn  func(*http.Request) *http.Response
		options *GetTransactionOptions
		want    *Transaction
		err     string
	}{
		"success — returns full transaction": {
			options: &GetTransactionOptions{TransactionID: 42},
			mockFn: func(req *http.Request) *http.Response {
				if !strings.Contains(req.URL.Path, "/transactions/42") {
					t.Errorf("GetTransaction() path = %q, want to contain /transactions/42", req.URL.Path)
				}
				if req.Method != http.MethodGet {
					t.Errorf("GetTransaction() method = %q, want GET", req.Method)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(tdTransaction.content)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{
				ID:     42,
				Payee:  "Coles Supermarket",
				Amount: -55.30,
				Date:   "2026-05-10",
				Note:   "⏰ 2026-05-10T20:29+10:00",
				Labels: []string{"Groceries", "Snack"},
				Category: Category{
					ID:    8973822,
					Title: "04 | 🥗 | Food",
				},
			},
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
					t.Fatalf("GetTransaction() expected error %q, got nil", tt.err)
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
				t.Errorf("ID = %d, want %d", got.ID, tt.want.ID)
			}
			if got.Payee != tt.want.Payee {
				t.Errorf("Payee = %q, want %q", got.Payee, tt.want.Payee)
			}
			if got.Amount != tt.want.Amount {
				t.Errorf("Amount = %.2f, want %.2f", got.Amount, tt.want.Amount)
			}
			if got.Date != tt.want.Date {
				t.Errorf("Date = %q, want %q", got.Date, tt.want.Date)
			}
			if got.Note != tt.want.Note {
				t.Errorf("Note = %q, want %q", got.Note, tt.want.Note)
			}
			if len(got.Labels) != len(tt.want.Labels) {
				t.Errorf("Labels len = %d, want %d", len(got.Labels), len(tt.want.Labels))
			} else {
				for i, l := range got.Labels {
					if l != tt.want.Labels[i] {
						t.Errorf("Labels[%d] = %q, want %q", i, l, tt.want.Labels[i])
					}
				}
			}
			if got.Category.ID != tt.want.Category.ID {
				t.Errorf("Category.ID = %d, want %d", got.Category.ID, tt.want.Category.ID)
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
				// verify method and path
				if req.Method != http.MethodPut {
					t.Errorf("UpdateTransaction() method = %q, want PUT", req.Method)
				}
				if !strings.Contains(req.URL.Path, "/transactions/42") {
					t.Errorf("UpdateTransaction() path = %q, want to contain /transactions/42", req.URL.Path)
				}
				// verify request body contains the new payee
				var body map[string]any
				json.NewDecoder(req.Body).Decode(&body)
				if body["payee"] != "Coles Express" {
					t.Errorf("request body payee = %v, want Coles Express", body["payee"])
				}
				tx := Transaction{ID: 42, Payee: "Coles Express", Amount: -55.30, Date: "2026-05-10"}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{ID: 42, Payee: "Coles Express", Amount: -55.30, Date: "2026-05-10"},
		},
		"success — updates note and labels": {
			options: &UpdateTransactionOptions{
				TransactionID: 42,
				Note:          "⏰ 2026-05-10T20:29+10:00",
				Labels:        "Groceries,Snack",
			},
			mockFn: func(req *http.Request) *http.Response {
				// verify request body contains note and labels
				var body map[string]any
				json.NewDecoder(req.Body).Decode(&body)
				if body["note"] != "⏰ 2026-05-10T20:29+10:00" {
					t.Errorf("request body note = %v, want timestamp", body["note"])
				}
				if body["labels"] != "Groceries,Snack" {
					t.Errorf("request body labels = %v, want Groceries,Snack", body["labels"])
				}
				tx := Transaction{
					ID:     42,
					Payee:  "Coles Supermarket",
					Note:   "⏰ 2026-05-10T20:29+10:00",
					Labels: []string{"Groceries", "Snack"},
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
				Note:   "⏰ 2026-05-10T20:29+10:00",
				Labels: []string{"Groceries", "Snack"},
			},
		},
		"success — updates category": {
			options: &UpdateTransactionOptions{
				TransactionID: 42,
				CategoryID:    func() *int32 { v := int32(8973822); return &v }(),
			},
			mockFn: func(req *http.Request) *http.Response {
				var body map[string]any
				json.NewDecoder(req.Body).Decode(&body)
				if body["category_id"] == nil {
					t.Error("request body missing category_id")
				}
				tx := Transaction{
					ID:       42,
					Payee:    "Coles Supermarket",
					Category: Category{ID: 8973822, Title: "04 | 🥗 | Food"},
				}
				b, _ := json.Marshal(tx)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			},
			want: &Transaction{
				ID:       42,
				Payee:    "Coles Supermarket",
				Category: Category{ID: 8973822, Title: "04 | 🥗 | Food"},
			},
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
					t.Fatalf("UpdateTransaction() expected error %q, got nil", tt.err)
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
				t.Errorf("ID = %d, want %d", got.ID, tt.want.ID)
			}
			if got.Payee != tt.want.Payee {
				t.Errorf("Payee = %q, want %q", got.Payee, tt.want.Payee)
			}
			if got.Note != tt.want.Note {
				t.Errorf("Note = %q, want %q", got.Note, tt.want.Note)
			}
			if len(got.Labels) != len(tt.want.Labels) {
				t.Errorf("Labels len = %d, want %d", len(got.Labels), len(tt.want.Labels))
			} else {
				for i, l := range got.Labels {
					if l != tt.want.Labels[i] {
						t.Errorf("Labels[%d] = %q, want %q", i, l, tt.want.Labels[i])
					}
				}
			}
			if got.Category.ID != tt.want.Category.ID {
				t.Errorf("Category.ID = %d, want %d", got.Category.ID, tt.want.Category.ID)
			}
		})
	}
}

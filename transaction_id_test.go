package pocketsmith

import (
	"encoding/json"
	"testing"
)

// Pocketsmith transaction IDs passed math.MaxInt32 in Oct 2026; they must decode.
func TestTransactionID_DecodesBeyondInt32(t *testing.T) {
	var tx Transaction
	if err := json.Unmarshal([]byte(`{"id": 2149437008}`), &tx); err != nil {
		t.Fatal(err)
	}
	if tx.ID != 2149437008 {
		t.Fatalf("got %d", tx.ID)
	}
}

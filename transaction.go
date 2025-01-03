package pocketsmith

import "time"

// Transaction defines a PocketSmith transaction.
type Transaction struct {
	ID                   int32              `json:"id"`
	Date                 string             `json:"date"`
	Payee                string             `json:"payee"`
	OriginalPayee        string             `json:"original_payee"`
	Amount               float64            `json:"amount"`
	UploadSource         string             `json:"upload_source"`
	ClosingBalance       float64            `json:"closing_balance"`
	Memo                 string             `json:"memo"`
	Note                 string             `json:"note"`
	Labels               []string           `json:"labels"`
	Type                 string             `json:"type"`
	Status               string             `json:"status"`
	IsTransfer           bool               `json:"is_transfer"`
	NeedsReview          bool               `json:"needs_review"`
	ChequeNumber         string             `json:"cheque_number"`
	AmountInBaseCurrency float64            `json:"amount_in_base_currency"`
	Category             Category           `json:"category"`
	TransactionAccount   TransactionAccount `json:"transaction_account"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

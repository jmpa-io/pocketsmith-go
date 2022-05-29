package pocketsmith

import (
	"fmt"
	"net/http"
	"time"
)

// Transaction defines a PocketSmith transaction.
// https://developers.pocketsmith.com/reference#get_transactions-id
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

// UpdateTransactionOptions defines the options for updating a transaction.
type UpdateTransactionOptions struct {
	ID           int32   `json:"-"`
	Labels       string  `json:"labels,omitempty"` // must be comma seperated list.
	Payee        string  `json:"payee,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	Date         string  `json:"date,omitempty"`
	IsTransfer   bool    `json:"is_transfer,omitempty"`
	CategoryID   int32   `json:"category_id,omitempty"`
	Note         string  `json:"note,omitempty"`
	Memo         string  `json:"memo,omitempty"`
	ChequeNumber string  `json:"cheque_number,omitempty"`
}

// UpdateTransaction updates a PocketSmith transaction.
// https://developers.pocketsmith.com/reference#put_transactions-id
func (c *Client) UpdateTransaction(options *UpdateTransactionOptions) (Transaction, error) {
	cr := clientRequest{
		method: http.MethodPut,
		path:   fmt.Sprintf("/transactions/%v", options.ID),
		data:   options,
	}
	var transaction Transaction
	_, err := c.sender(cr, &transaction)
	return transaction, err
}

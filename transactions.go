package pocketsmith

import (
	"fmt"
	"net/http"
)

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
// https://developers.pocketsmith.com/reference/put_transactions-id-1
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

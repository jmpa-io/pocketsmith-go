package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Transactions represents a slice of Transaction.
type Transactions []Transaction

// UpdateTransactionOptions defines the options for updating a transaction in
// Pocketsmith, by the given transaction id.
type UpdateTransactionOptions struct {
	TransactionID int32   `json:"-"                       validator:"required"`
	Labels        string  `json:"labels,omitempty"` // must be comma seperated list.
	Payee         string  `json:"payee,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	Date          string  `json:"date,omitempty"`
	IsTransfer    bool    `json:"is_transfer,omitempty"`
	CategoryID    int32   `json:"category_id,omitempty"`
	Note          string  `json:"note,omitempty"`
	Memo          string  `json:"memo,omitempty"`
	ChequeNumber  string  `json:"cheque_number,omitempty"`
}

// UpdateTransaction updates a transaction in Pocketsmith, by the given
// transaction id.
// https://developers.pocketsmith.com/reference/put_transactions-id-1.
func (c *Client) UpdateTransaction(
	ctx context.Context,
	options *UpdateTransactionOptions,
) (transaction *Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "UpdateTransaction")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	//  update transaction.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPut,
		path:   fmt.Sprintf("/transactions/%v", options.TransactionID),
		body:   options,
	}, &transaction)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to update transaction: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return transaction, nil
}

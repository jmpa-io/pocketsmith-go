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

// GetTransactionOptions defines the options for retrieving a single
// transaction from Pocketsmith by its ID.
type GetTransactionOptions struct {
	TransactionID int32 `json:"-" validator:"required"`
}

// GetTransaction retrieves a single transaction from Pocketsmith by its ID.
// https://developers.pocketsmith.com/reference/get_transactions-id-1.
func (c *Client) GetTransaction(
	ctx context.Context,
	options *GetTransactionOptions,
) (transaction *Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetTransaction")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// get transaction.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/transactions/%v", options.TransactionID),
	}, &transaction)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get transaction: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return transaction, nil
}
// Pocketsmith, by the given transaction id.
// Pointer fields are used for values where zero is a valid update
// (e.g. Amount=0.00, CategoryID=0 to un-assign, IsTransfer=false to un-mark).
// Nil pointer means "do not update this field".
type UpdateTransactionOptions struct {
	TransactionID int32    `json:"-"                       validator:"required"`
	Labels        string   `json:"labels,omitempty"`        // comma-separated list; use ClearLabels=true to first clear all existing labels
	ClearLabels   bool     `json:"-"`                       // if true, send labels="" before setting Labels (avoids Pocketsmith dormant-record bug)
	Payee         string   `json:"payee,omitempty"`
	Amount        *float64 `json:"amount,omitempty"`        // pointer: nil = don't update, 0.00 = valid
	Date          string   `json:"date,omitempty"`
	IsTransfer    *bool    `json:"is_transfer,omitempty"`   // pointer: nil = don't update, false = un-mark
	CategoryID    *int32   `json:"category_id,omitempty"`   // pointer: nil = don't update, 0 = un-assign
	Note          string   `json:"note,omitempty"`
	Memo          string   `json:"memo,omitempty"`
	ChequeNumber  string   `json:"cheque_number,omitempty"`
}

// UpdateTransaction updates a transaction in Pocketsmith, by the given
// transaction id.
// https://developers.pocketsmith.com/reference/put_transactions-id-1.
// If options.ClearLabels is true, a preliminary PUT with labels="" is sent
// before the main update to avoid Pocketsmith's dormant-label-record bug
// (where setting a previously-used label re-activates all stored duplicate
// records of that label name).
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

	// If ClearLabels is set, first send an empty-label PUT to wipe all stored
	// label records, then proceed with the intended labels.
	if options.ClearLabels {
		type clearBody struct {
			Labels string `json:"labels"`
		}
		if _, err := c.sender(newCtx, senderRequest{
			method: http.MethodPut,
			path:   fmt.Sprintf("/transactions/%v", options.TransactionID),
			body:   &clearBody{Labels: ""},
		}, nil); err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to clear labels: %v", err))
			span.RecordError(err)
			return nil, fmt.Errorf("failed to clear labels before update: %w", err)
		}
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

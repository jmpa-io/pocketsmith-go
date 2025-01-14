package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// TransactionAccounts represents a slice of TransactionAccount.
type TransactionAccounts []TransactionAccount

type ListTransactionAccountsOptions struct {
	userID int `validator:"required"`
}

// ListTransactionAccounts lists the transaction accounts for the given user id.
// https://developers.pocketsmith.com/reference/get_users-id-transaction-accounts-1.
func (c *Client) ListTransactionAccounts(
	ctx context.Context,
	options *ListTransactionAccountsOptions,
) (accounts TransactionAccounts, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactionAccounts")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list transaction accounts.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/transaction_accounts", options.userID),
	}, &accounts)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get user: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return accounts, nil
}

// ListTransactionAccountsForAuthedUser lists the transaction accounts for the authed user.
func (c *Client) ListTransactionAccountsForAuthedUser(
	ctx context.Context,
) (TransactionAccounts, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactionAccountsForAuthedUser")
	defer span.End()

	// list transaction accounts for authed user.
	return c.ListTransactionAccounts(
		newCtx,
		&ListTransactionAccountsOptions{userID: c.authedUser.ID},
	)
}

// CreateTransactionAccountTransactionOptions defines the options for creating
// a transaction in a transaction account.
type CreateTransactionAccountTransactionOptions struct {
	accountID    int     `json:"-"                       validator:"required"`
	Date         string  `json:"date"`
	Payee        string  `json:"payee"`
	Amount       float64 `json:"amount"`
	Labels       string  `json:"labels,omitempty"` // must be comma seperated.
	CategoryID   int32   `json:"category_id,omitempty"`
	Note         string  `json:"note,omitempty"`
	Memo         string  `json:"memo,omitempty"`
	IsTransfer   bool    `json:"is_transfer,omitempty"`
	ChequeNumber string  `json:"cheque_number,omitempty"`
	NeedsReview  bool    `json:"needs_review,omitempty"`
}

// CreateTransactionAccountTransaction, using the given account id, creates a
// a transaction in an transaction account.
// https://developers.pocketsmith.com/reference/post_transaction-accounts-id-transactions-1
func (c *Client) CreateTransactionAccountTransaction(
	ctx context.Context,
	options *CreateTransactionAccountTransactionOptions,
) (transaction *Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).
		Start(ctx, "CreateTransactionAccountTransaction")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// create transaction account transaction.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions", options.accountID),
		body:   options,
	}, &transaction)
	if err != nil {
		span.SetStatus(
			codes.Error,
			fmt.Sprintf("failed to create transaction account transaction: %v", err),
		)
		span.RecordError(err)
		return nil, err
	}
	return transaction, nil
}

// ListTransactionAccountTransactionsOptions defines the options for listing
// transactions in a transaction account.
type ListTransactionAccountTransactionsOptions struct {
	AccountID         string `json:"-"                            validator:"required"`
	StartDate         string `json:"start_date,omitempty"`
	EndDate           string `json:"end_date,omitempty"`
	OnlyUncategorised int32  `json:"only_uncategorized,omitempty"`
	Type              string `json:"type,omitempty"`
}

// ListTransactionAccountTransactions, using the given account id, lists the
// transactions for a transaction account.
// https://developers.pocketsmith.com/reference/get_transaction-accounts-id-transactions-1
func (c *Client) ListTransactionAccountTransactions(
	ctx context.Context,
	options *ListTransactionAccountTransactionsOptions,
) (transactions []Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).
		Start(ctx, "ListTransactionAccountTransactions")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// setup request.
	sr := senderRequest{
		method: http.MethodGet,
		path: fmt.Sprintf(
			"/transaction_accounts/%v/transactions?per_page=100",
			options.AccountID,
		),
	}

	// list transaction account transactions.
	for {

		// get batch.
		var batch []Transaction
		resp, err := c.sender(newCtx, sr, &batch)
		if err != nil {
			span.SetStatus(
				codes.Error,
				fmt.Sprintf("failed to list transaction account transactions: %v", err),
			)
			span.RecordError(err)
			return nil, err
		}

		// extract batch data.
		transactions = append(transactions, batch...)

		// paginate?
		next := getHeader(resp.Header, "next")
		if next == "" {
			break
		}
		sr.path = strings.Replace(next, c.endpoint, "", -1)
	}
	return transactions, nil
}

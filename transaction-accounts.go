package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// TransactionAccounts represents a slice of TransactionAccount.
type TransactionAccounts []TransactionAccount

// LsitTransactionAccountsForUserOptions defines options for listing
// transaction accounts from Pocketsmith for the given user, by the user id.
type ListTransactionAccountsForUserOptions struct {
	UserID int `validator:"required"`
}

// ListTransactionAccounts lists the transaction accounts from Pocketsmith for
// the given user, by the user id.
// https://developers.pocketsmith.com/reference/get_users-id-transaction-accounts-1.
func (c *Client) ListTransactionAccountsForUser(
	ctx context.Context,
	options *ListTransactionAccountsForUserOptions,
) (accounts TransactionAccounts, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactionAccountsForUser")
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
		path:   fmt.Sprintf("/users/%v/transaction_accounts", options.UserID),
	}, &accounts)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get user: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return accounts, nil
}

// ListTransactionAccounts lists the transaction accounts from Pocketsmith
// under the authed user.
// https://developers.pocketsmith.com/reference/get_users-id-transaction-accounts-1.
func (c *Client) ListTransactionAccounts(ctx context.Context) (TransactionAccounts, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactionAccounts")
	defer span.End()

	// list transaction accounts for authed user.
	return c.ListTransactionAccountsForUser(
		newCtx,
		&ListTransactionAccountsForUserOptions{UserID: c.authedUser.ID},
	)
}

// CreateTransactionAccountTransactionOptions defines the options for creating
// a transaction in the given transaction account in Pocketsmith, by the
// transaction account id.
type CreateTransactionAccountTransactionOptions struct {
	TransactionAccountID int     `json:"-"                       validator:"required"`
	Payee                string  `json:"payee"                   validator:"required"`
	Amount               float64 `json:"amount"                  validator:"required"`
	Date                 string  `json:"date"                    validator:"required"` //TODO: should this be customTime?
	IsTransfer           bool    `json:"is_transfer,omitempty"`
	Labels               string  `json:"labels,omitempty"` // must be comma seperated. // TODO: should this be a []string or a custom type?
	CategoryID           int32   `json:"category_id,omitempty"`
	Note                 string  `json:"note,omitempty"`
	Memo                 string  `json:"memo,omitempty"`
	ChequeNumber         string  `json:"cheque_number,omitempty"`
	NeedsReview          bool    `json:"needs_review,omitempty"`
}

// CreateTransactionAccountTransaction creates a transaction in the given
// transaction account in Pocketsmith, by the transaction account id.
// https://developers.pocketsmith.com/reference/post_transaction-accounts-id-transactions-1.
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
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions", options.TransactionAccountID),
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

type ListTransactionAccountTransactionsOptionType string

const (
	ListTransactionAccountTransactionsOptionTypeDebit  ListTransactionAccountTransactionsOptionType = "debit"
	ListTransactionAccountTransactionsOptionTypeCredit ListTransactionAccountTransactionsOptionType = "credit"
)

// ListTra609534nsactionAccountTransactionsOptions defines the options for l
// isting transactions in a transaction account from Pocketsmith, by the
// transaction account id.
type ListTransactionAccountTransactionsOptions struct {
	TransactionAccountID string                                       `json:"-"                       validator:"required"`
	StartDate            string                                       `json:"start_date,omitempty"` // TODO: should this be customTime?
	EndDate              string                                       `json:"end_date,omitempty"`   // TODO: should this be customTime?
	UpdatedSince         time.Time                                    `json:"updated_since,omitempty"`
	Uncategorised        int32                                        `json:"uncategorized,omitempty"` // TODO: should this be a bool?
	Type                 ListTransactionAccountTransactionsOptionType `json:"type,omitempty"`
	NeedsReview          int32                                        `json:"needs_review,omitempty"` // TODO: should this be a bool?
	Search               string                                       `json:"search,omitempty"`
}

// ListTransactionAccountTransactions lists transactions in a transaction
// account from Pocketsmith, by the transaction account id.
// https://developers.pocketsmith.com/reference/get_transaction-accounts-id-transactions-1.
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
			"/transaction_accounts/%v/transactions",
			options.TransactionAccountID,
		),
		queries: setupQueries(nil),
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

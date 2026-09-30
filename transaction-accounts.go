package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// defaultPageWorkers is the default number of concurrent page-fetch goroutines
// used by fetchAllPages. Pocketsmith's rate limit is 5000 requests/hour;
// 10 concurrent workers is conservative and safe.
const defaultPageWorkers = 10

// TransactionAccounts represents a slice of TransactionAccount.
type TransactionAccounts []TransactionAccount

// LsitTransactionAccountsForUserOptions defines options for listing
// transaction accounts from Pocketsmith for the given user, by the user id.
type ListTransactionAccountsForUserOptions struct {
	UserID int `validate:"required"`
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
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list transaction accounts: %v", err))
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

	userID, err := c.authedUserID(newCtx)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get authed user: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list transaction accounts for authed user.
	accounts, err := c.ListTransactionAccountsForUser(
		newCtx,
		&ListTransactionAccountsForUserOptions{UserID: userID},
	)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list transaction accounts: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return accounts, nil
}

// CreateTransactionAccountTransactionOptions defines the options for creating
// a transaction in the given transaction account in Pocketsmith, by the
// transaction account id.
type CreateTransactionAccountTransactionOptions struct {
	TransactionAccountID int     `json:"-"                       validate:"required"`
	Payee                string  `json:"payee"                   validate:"required"`
	Amount               float64 `json:"amount"` // no validate:"required" — 0.00 is a valid amount
	Date                 string  `json:"date"                    validate:"required"` //TODO: should this be customTime?
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

// ListTransactionAccountTransactionsOptions defines the options for listing
// transactions in a transaction account from Pocketsmith, by the
// transaction account id.
type ListTransactionAccountTransactionsOptions struct {
	TransactionAccountID int                                          `json:"-"                       validate:"required"`
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

	// setup query params — only non-zero/non-empty values are included.
	queryMap := map[string]string{
		"start_date": options.StartDate,
		"end_date":   options.EndDate,
		"search":     options.Search,
		"type":       string(options.Type),
	}
	if options.Uncategorised != 0 {
		queryMap["uncategorised"] = fmt.Sprintf("%v", options.Uncategorised)
	}
	if options.NeedsReview != 0 {
		queryMap["needs_review"] = fmt.Sprintf("%v", options.NeedsReview)
	}
	if !options.UpdatedSince.IsZero() {
		queryMap["updated_since"] = options.UpdatedSince.Format(time.RFC3339)
	}

	// setup request.
	sr := senderRequest{
		method: http.MethodGet,
		path: fmt.Sprintf(
			"/transaction_accounts/%v/transactions",
			options.TransactionAccountID,
		),
		queries: setupQueries(&queryMap),
	}

	// list transaction account transactions concurrently across all pages.
	transactions, err = fetchAllPages[Transaction](newCtx, c, sr, defaultPageWorkers)
	if err != nil {
		span.SetStatus(
			codes.Error,
			fmt.Sprintf("failed to list transaction account transactions: %v", err),
		)
		span.RecordError(err)
		return nil, err
	}
	return transactions, nil
}

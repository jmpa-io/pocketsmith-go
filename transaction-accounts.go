package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
)

// ListTransactionAccounts, using the given user id, lists transaction accounts
// for a user.
// https://developers.pocketsmith.com/reference/get_users-id-transaction-accounts-1
func (c *Client) ListTransactionAccounts(
	ctx context.Context,
	userId int,
) (accounts []TransactionAccount, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListInstitutions")
	defer span.End()

	// list transaction accounts.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/transaction_accounts", userId),
	}, &accounts)
	return accounts, err
}

// ListTransactionAccountsForAuthedUser, using the token attached to the client,
// lists transaction accounts for the authed user.
func (c *Client) ListTransactionAccountsForAuthedUser(
	ctx context.Context,
) ([]TransactionAccount, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListInstitutions")
	defer span.End()

	// list transaction accounts for authed user.
	return c.ListTransactionAccounts(newCtx, c.authedUser.ID)
}

// CreateTransactionAccountTransactionOptions defines the options for creating
// a transaction in a transaction account.
type CreateTransactionAccountTransactionOptions struct {
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
	accountId int,
	options *CreateTransactionAccountTransactionOptions,
) (transaction *Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).
		Start(ctx, "CreateTransactionAccountTransaction")
	defer span.End()

	// create transaction account transaction.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions", accountId),
		body:   options,
	}, &transaction)
	return transaction, err
}

// ListTransactionAccountTransactionsOptions defines the options for listing
// transactions in a transaction account.
type ListTransactionAccountTransactionsOptions struct {
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
	accountId int,
	options *ListTransactionAccountTransactionsOptions,
) (transactions []Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).
		Start(ctx, "ListTransactionAccountTransactions")
	defer span.End()

	// setup request.
	sr := senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions?per_page=100", accountId),
	}

	// list transaction account transactions.
	for {

		// get batch.
		var batch []Transaction
		resp, err := c.sender(newCtx, sr, &batch)
		if err != nil {
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

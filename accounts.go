package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
)

// CreateAccountOptions defines the options for creating an account for a user.
type CreateAccountOptions struct {
	InstitutionID int    `json:"institution_id"`
	Title         string `json:"title"`
	CurrencyCode  string `json:"currency_code"`
	Type          string `json:"type"` // TODO enum?
}

// CreateAccount, using the given user id, creates an account for a user.
// https://developers.pocketsmith.com/reference#post_users-id-accounts.
func (c *Client) CreateAccount(
	ctx context.Context,
	userId int,
	options *CreateAccountOptions,
) (*Account, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAccount")
	defer span.End()

	// create account.
	var account *Account
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/accounts", userId),
		data:   options,
	}, &account)
	return account, err
}

// CreateAccountForAuthedUser, using the token attached to the client, creates
// an account for the authed user.
func (c *Client) CreateAccountForAuthedUser(
	ctx context.Context,
	options *CreateAccountOptions,
) (*Account, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAccountForAuthedUser")
	defer span.End()

	// create account for authed user.
	return c.CreateAccount(newCtx, c.user.ID, options)
}

// DeleteAccount, using the given account id, deletes an account.
// https://developers.pocketsmith.com/reference#delete_accounts-id.
func (c *Client) DeleteAccount(ctx context.Context, accountId int) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteAccount")
	defer span.End()

	// delete account.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/accounts/%v", accountId),
	}, nil)
	return err
}

// ListAccounts, using the given user id, returns a list of account for a user.
// https://developers.pocketsmith.com/reference#get_users-id-accounts.
func (c *Client) ListAccounts(ctx context.Context, userId int) ([]Account, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccounts")
	defer span.End()

	// list accounts.
	var accounts []Account
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/accounts", userId),
	}, &accounts)
	return accounts, err
}

// ListAccountsForAuthedUser, using the token attached to the client, returns a
// list of accounts for the authed user.
func (c *Client) ListAccountsForAuthedUser(ctx context.Context) ([]Account, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccountsForAuthedUser")
	defer span.End()

	// list accounts for authed user.
	return c.ListAccounts(newCtx, c.user.ID)
}

// ListAccountTransactionsOptions defines the options for listing
// transactions in an account.
type ListAccountTransactionsOptions struct {
	StartDate     string `json:"start_date,omitempty"`
	EndDate       string `json:"end_date,omitempty"`
	UpdatedSince  string `json:"updated_since,omitempty"`
	Uncategorised int8   `json:"uncategorised,omitempty"`
	Type          string `json:"type,omitempty"`
	NeedsReview   int8   `json:"needs_review,omitempty"`
	Search        string `json:"search,omitempty"`
	Page          int    `json:"page,omitempty"`
}

// ListTransactionAccountTransactions, using the given account id, lists the
// transactions for an account.
// https://developers.pocketsmith.com/reference/get_accounts-id-transactions-1
func (c *Client) ListAccountTransactions(
	ctx context.Context,
	accountId int,
	options *ListAccountTransactionsOptions,
) (transactions []Transaction, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccountTransactions")
	defer span.End()

	// setup request.
	sr := senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/accounts/%v/transactions?per_page=100", accountId),
		data:   options,
	}

	// retrieve transactions for account.
	for {

		// get respons
		var batch []Transaction
		resp, err := c.sender(newCtx, sr, &batch)
		if err != nil {
			return nil, err
		}
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

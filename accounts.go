package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Accounts represents a slice of Account.
type Accounts []Account

// CreateAccountOptions defines the options for creating an account for a user.
type CreateAccountOptions struct {
	InstitutionID int    `json:"institution_id"`
	Title         string `json:"title"`
	CurrencyCode  string `json:"currency_code"`
	Type          string `json:"type"` // TODO enum?
}

// CreateAccountForUserOptions ...
type CreateAccountForUserOptions struct {
	UserID int `json:"-" validator:"required"`

	CreateAccountOptions
}

// CreateAccountForUser, using the given user id, creates an account for a user.
// https://developers.pocketsmith.com/reference#post_users-id-accounts.
func (c *Client) CreateAccountForUser(
	ctx context.Context,
	options *CreateAccountForUserOptions,
) (*Account, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAccountForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// create account.
	var account *Account
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/accounts", options.UserID),
		body:   options,
	}, &account)

	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to create account: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return account, err
}

// CreateAccount, using the token attached to the client, creates
// an account for the authed user.
func (c *Client) CreateAccount(
	ctx context.Context,
	options *CreateAccountOptions,
) (*Account, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateAccount")
	defer span.End()

	// create account for authed user.
	return c.CreateAccountForUser(
		newCtx,
		&CreateAccountForUserOptions{UserID: c.authedUser.ID, CreateAccountOptions: *options},
	)
}

// DeleteAccountOptions ...
type DeleteAccountOptions struct {
	AccountID int `validator:"required"`
}

// DeleteAccount, using the given account id, deletes an account.
// https://developers.pocketsmith.com/reference#delete_accounts-id.
func (c *Client) DeleteAccount(ctx context.Context, options *DeleteAccountOptions) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteAccount")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return err
	}

	// delete account.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/accounts/%v", options.AccountID),
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to delete account: %v", err))
		span.RecordError(err)
		return err
	}
	return nil
}

// ListAccountsOptions ...
type ListAccountsForUserOptions struct {
	UserID int `validator:"required"`
}

// ListAccountsForUser, using the given user id, returns a list of account for a user.
// https://developers.pocketsmith.com/reference#get_users-id-accounts.
func (c *Client) ListAccountsForUser(
	ctx context.Context,
	options *ListAccountsForUserOptions,
) (accounts Accounts, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccountsForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list accounts.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/accounts", options.UserID),
	}, &accounts)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list accounts: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return accounts, nil
}

// ListAccounts, using the token attached to the client, returns a
// list of accounts for the authed user.
func (c *Client) ListAccounts(ctx context.Context) (Accounts, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccounts")
	defer span.End()

	// list accounts for authed user.
	return c.ListAccountsForUser(newCtx, &ListAccountsForUserOptions{UserID: c.authedUser.ID})
}

// ListAccountTransactionsOptions defines the options for listing
// transactions in an account.
type ListAccountTransactionsOptions struct {
	AccountID     int    `json:"-"                       validator:"required"`
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
	options *ListAccountTransactionsOptions,
) (transactions Transactions, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccountTransactions")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// setup request.
	sr := senderRequest{
		method:  http.MethodGet,
		path:    fmt.Sprintf("/accounts/%v/transactions", options.AccountID),
		body:    options,
		queries: setupQueries(nil),
	}

	// retrieve transactions for account.
	for {

		// get respons
		var batch Transactions
		resp, err := c.sender(newCtx, sr, &batch)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to list account transactions: %v", err))
			span.RecordError(err)
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

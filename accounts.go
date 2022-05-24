package pocketsmith

import (
	"fmt"
	"net/http"
	"time"
)

// Account defines a PocketSmith account.
// https://developers.pocketsmith.com/reference#get_accounts-id
type Account struct {
	ID                           int                  `json:"id"`
	Title                        string               `json:"title"`
	Type                         string               `json:"type"`
	IsNetWorth                   bool                 `json:"is_net_worth"`
	CurrencyCode                 string               `json:"currency_code"`
	CurrentBalance               float64              `json:"current_balance"`
	CurrentBalanceInBaseCurrency float64              `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64              `json:"current_balance_exchange_rate"`
	CurrentBalanceDate           string               `json:"current_balance_date"`
	PrimaryTransactionAccount    TransactionAccount   `json:"primary_transaction_account"`
	TransactionAccounts          []TransactionAccount `json:"transaction_accounts"`
	CreatedAt                    time.Time            `json:"created_at"`
	UpdatedAt                    time.Time            `json:"updated_at"`
}

// CreateAccountOptions defines the options for creating an account for a user.
type CreateAccountOptions struct {
	InstitutionID int    `json:"institution_id"`
	Title         string `json:"title"`
	CurrencyCode  string `json:"currency_code"`
	Type          string `json:"type"` // TODO enum?
}

// CreateAccount, using the given user id, creates an account for a user.
// https://developers.pocketsmith.com/reference#post_users-id-accounts
func (c *Client) CreateAccount(userId int, options CreateAccountOptions) (*Account, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/accounts", userId),
		data:   options,
	}
	var account *Account
	_, err := c.sender(cr, &account)
	return account, err
}

// CreateAccountForAuthedUser, using the token attached to the client, creates
// an account for the authed user.
func (c *Client) CreateAccountForAuthedUser(options CreateAccountOptions) (*Account, error) {
	return c.CreateAccount(c.user.ID, options)
}

// DeleteAccount, using the given account id, deletes an account.
// https://developers.pocketsmith.com/reference#delete_accounts-id
func (c *Client) DeleteAccount(accountId int) error {
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/accounts/%v", accountId),
	}
	_, err := c.sender(cr, nil)
	return err
}

// ListAccounts, using the given user id, returns a list of account for a user.
// https://developers.pocketsmith.com/reference#get_users-id-accounts
func (c *Client) ListAccounts(userId int) ([]Account, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/accounts", userId),
	}
	var accounts []Account
	_, err := c.sender(cr, &accounts)
	return accounts, err
}

// ListAccountsForAuthedUser, using the token attached to the client, returns a
// list of accounts for the authed user.
func (c *Client) ListAccountsForAuthedUser() ([]Account, error) {
	return c.ListAccounts(c.user.ID)
}

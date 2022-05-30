package pocketsmith

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// TransactionAccount defines a PocketSmith transaction account.
// https://developers.pocketsmith.com/reference#get_transaction-accounts-id
type TransactionAccount struct {
	ID                           int         `json:"id"`
	Name                         string      `json:"name"`
	Number                       string      `json:"number"`
	Type                         string      `json:"type"`
	CurrencyCode                 string      `json:"currency_code"`
	CurrentBalance               float64     `json:"current_balance"`
	CurrentBalanceInBaseCurrency float64     `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64     `json:"current_balance_exchange_rate"`
	CurrentBalanceDate           string      `json:"current_balance_date"`
	StartingBalance              float64     `json:"starting_balance"`
	StartingBalanceDate          string      `json:"starting_balance_date"`
	Institution                  Institution `json:"institution"`
	CreatedAt                    time.Time   `json:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at"`
}

// ListTransactionAccounts, using the given user id, lists transaction accounts
// for a user.
// https://developers.pocketsmith.com/reference#get_users-id-transaction-accounts
func (c *Client) ListTransactionAccounts(userId int) ([]TransactionAccount, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/transaction_accounts", userId),
	}
	var accounts []TransactionAccount
	_, err := c.sender(cr, &accounts)
	return accounts, err
}

// ListTransactionAccountsForAuthedUser, using the token attached to the client,
// lists transaction accounts for the authed user.
func (c *Client) ListTransactionAccountsForAuthedUser() ([]TransactionAccount, error) {
	return c.ListTransactionAccounts(c.user.ID)
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
// https://developers.pocketsmith.com/reference#post_transaction-accounts-id-transactions
func (c *Client) CreateTransactionAccountTransaction(accountId int, options *CreateTransactionAccountTransactionOptions) (*Transaction, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions", accountId),
		data:   options,
	}
	var transaction *Transaction
	_, err := c.sender(cr, &transaction)
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
// https://developers.pocketsmith.com/reference#get_transaction-accounts-id-transactions
func (c *Client) ListTransactionAccountTransactions(accountId int, options *ListTransactionAccountTransactionsOptions) ([]Transaction, error) {
	var transactions []Transaction
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/transaction_accounts/%v/transactions?per_page=100", accountId),
	}
	for {
		var batch []Transaction
		resp, err := c.sender(cr, &batch)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, batch...)
		next := getHeader(resp.Header, "next")
		if next == "" {
			break
		}
		cr.path = strings.Replace(next, c.endpoint, "", -1)
	}
	return transactions, nil
}

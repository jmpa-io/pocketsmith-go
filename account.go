package pocketsmith

import "time"

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

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

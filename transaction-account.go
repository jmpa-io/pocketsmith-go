package pocketsmith

import "time"

// TransactionAccount defines a PocketSmith transaction account.
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

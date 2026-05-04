package pocketsmith

import "time"

// Scenario defines a Pocketsmith scenario (account forecast container).
type Scenario struct {
	ID             int       `json:"id"`
	AccountID      int       `json:"account_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	InterestRate   float64   `json:"interest_rate"`
	StartingBalance float64  `json:"starting_balance"`
	CurrentBalance float64   `json:"current_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

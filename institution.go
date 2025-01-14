package pocketsmith

import "time"

// Institution defines a PocketSmith institution.
type Institution struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	CurrencyCode string    `json:"currency_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

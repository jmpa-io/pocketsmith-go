package pocketsmith

import "time"

// Category defines a PocketSmith category.
type Category struct {
	ID         int32       `json:"id"`
	Title      string      `json:"title"`
	Colour     string      `json:"colour"`
	Children   []*Category `json:"children"`
	ParentID   *int        `json:"parent_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	IsTransfer bool        `json:"is_transfer"`
}

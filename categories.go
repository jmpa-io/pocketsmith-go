package pocketsmith

import (
	"fmt"
	"net/http"
	"time"
)

// Category defines a PocketSmith category.
// https://developers.pocketsmith.com/reference#get_categories-id
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

// ListCategories, using the given user id, lists the categories for a user.
// https://developers.pocketsmith.com/reference#get_users-id-categories
func (c *Client) ListCategories(userId int) ([]Category, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/categories", userId),
	}
	var categories []Category
	_, err := c.sender(cr, &categories)
	return categories, err
}

// ListCategoriesForAuthedUser, using the token attached to a client, lists the
// categories for the authed user.
func (c *Client) ListCategoriesForAuthedUser() ([]Category, error) {
	return c.ListCategories(c.user.ID)
}

// GetCategoryByTitle, using the given user id and category, returns the found
// category for a user.
func (c *Client) GetCategoryByTitle(userId int, category string) (*Category, error) {
	categories, err := c.ListCategories(userId)
	if err != nil {
		return nil, err
	}
	for _, c := range categories {
		if c.Title == category {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("failed to find category with title %q for user %v", category, userId)
}

// GetCategoryByTitleForAuthedUser, using the token attached to the client,
// returns the found category for the authed user.
func (c *Client) GetCategoryByTitleForAuthedUser(category string) (*Category, error) {
	return c.GetCategoryByTitle(c.user.ID, category)
}

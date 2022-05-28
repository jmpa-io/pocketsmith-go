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

// CreateCategoryOptions defines the options for creating a catagory for a user.
type CreateCategoryOptions struct {
	Title           string `json:"title"`
	Colour          string `json:"colour,omitempty"`
	ParentID        string `json:"parent_id,omitempty"`
	IsTransfer      bool   `json:"is_transfer,omitempty"`
	IsBill          bool   `json:"is_bill,omitempty"`
	RollUp          bool   `json:"roll_up,omitempty"`
	RefundBehaviour string `json:"refund_behaviour,omitempty"`
}

// CreateCategory, using the given user id, creates a category for a user.
// https://developers.pocketsmith.com/reference/post_users-id-categories-1
func (c *Client) CreateCategory(userId int, options CreateCategoryOptions) error {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/categories", userId),
		data:   options,
	}
	_, err := c.sender(cr, nil)
	return err
}

// CreateCategoryForAuthedUser, using the token attached to a client, creates a
// category for the authed user.
func (c *Client) CreateCategoryForAuthedUser(options CreateCategoryOptions) error {
	return c.CreateCategory(c.user.ID, options)
}

// DeleteCategory, using the given category id, deletes a category.
// https://developers.pocketsmith.com/reference/delete_categories-id-1
func (c *Client) DeleteCategory(categoryId int32) error {
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/categories/%v", categoryId),
	}
	_, err := c.sender(cr, nil)
	return err
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

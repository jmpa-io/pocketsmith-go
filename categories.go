package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

// Categories represents a slice of Category.
type Categories []Category

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
func (c *Client) CreateCategory(
	ctx context.Context,
	userId int,
	options *CreateCategoryOptions,
) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateCategory")
	defer span.End()

	// create category.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/categories", userId),
		body:   options,
	}, nil)
	return err
}

// CreateCategoryForAuthedUser, using the token attached to a client, creates a
// category for the authed user.
func (c *Client) CreateCategoryForAuthedUser(
	ctx context.Context,
	options *CreateCategoryOptions,
) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateCategoryForAuthedUser")
	defer span.End()

	// create category for authed user.
	return c.CreateCategory(newCtx, c.authedUser.ID, options)
}

// DeleteCategory, using the given category id, deletes a category.
// https://developers.pocketsmith.com/reference/delete_categories-id-1
func (c *Client) DeleteCategory(ctx context.Context, categoryId int32) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateCategory")
	defer span.End()

	// delete category.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/categories/%v", categoryId),
	}, nil)
	return err
}

// ListCategories, using the given user id, lists the categories for a user.
// https://developers.pocketsmith.com/reference#get_users-id-categories
func (c *Client) ListCategories(
	ctx context.Context,
	userId int,
) (categories []Category, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListCategories")
	defer span.End()

	// list categories.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/categories", userId),
	}, &categories)
	return categories, err
}

// ListCategoriesForAuthedUser, using the token attached to a client, lists the
// categories for the authed user.
func (c *Client) ListCategoriesForAuthedUser(ctx context.Context) ([]Category, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListCategoriesForAuthedUser")
	defer span.End()

	// list categories for authed user.
	return c.ListCategories(newCtx, c.authedUser.ID)
}

// GetCategoryByTitle, using the given user id and category, returns the found
// category for a user.
func (c *Client) GetCategoryByTitle(
	ctx context.Context,
	userId int,
	category string,
) (*Category, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetCategoryByTitle")
	defer span.End()

	// get category by title.
	categories, err := c.ListCategories(newCtx, userId)
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
func (c *Client) GetCategoryByTitleForAuthedUser(
	ctx context.Context,
	category string,
) (*Category, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetCategoryByTitleForAuthedUser")
	defer span.End()

	// get category by title.
	return c.GetCategoryByTitle(newCtx, c.authedUser.ID, category)
}

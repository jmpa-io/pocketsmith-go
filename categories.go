package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Categories represents a slice of Category.
type Categories []Category

// CreateCategoryOptions defines the options for creating a catagory for a user.
type CreateCategoryOptions struct {
	Title           string `json:"title"                      validator:"required"`
	Colour          string `json:"colour,omitempty"`
	ParentID        string `json:"parent_id,omitempty"`
	IsTransfer      bool   `json:"is_transfer,omitempty"`
	IsBill          bool   `json:"is_bill,omitempty"`
	RollUp          bool   `json:"roll_up,omitempty"`
	RefundBehaviour string `json:"refund_behaviour,omitempty"`
}

// CreateCategoryForUser ...
type CreateCategoryForUserOptions struct {
	UserID int `json:"-" validator:"required"`

	CreateCategoryOptions
}

// CreateCategoryForUser, using the given user id, creates a category for a user.
// https://developers.pocketsmith.com/reference/post_users-id-categories-1
func (c *Client) CreateCategoryForUser(
	ctx context.Context,
	options *CreateCategoryForUserOptions,
) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateCategoryForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return err
	}

	// create category.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/categories", options.UserID),
		body:   options,
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to create category: %v", err))
		span.RecordError(err)
		return err
	}
	return nil
}

// CreateCategoryForAuthedUser, using the token attached to a client, creates a
// category for the authed user.
func (c *Client) CreateCategory(
	ctx context.Context,
	options *CreateCategoryOptions,
) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateCategory")
	defer span.End()

	// create category for authed user.
	return c.CreateCategoryForUser(newCtx, &CreateCategoryForUserOptions{
		UserID:                c.authedUser.ID,
		CreateCategoryOptions: *options,
	})
}

// DeleteCategoryOptions ...
type DeleteCategoryOptions struct {
	CategoryID int32 `json:"-" validator:"required"`
}

// DeleteCategory, using the given category id, deletes a category.
// https://developers.pocketsmith.com/reference/delete_categories-id-1
func (c *Client) DeleteCategory(ctx context.Context, options *DeleteCategoryOptions) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteCategory")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return err
	}

	// delete category.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/categories/%v", options.CategoryID),
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to delete category: %v", err))
		span.RecordError(err)
		return err
	}
	return nil
}

// ListCategoriesOptions ...
type ListCategoriesForUserOptions struct {
	UserID int `json:"-" validator:"required"`
}

// ListCategoriesForUser, using the given user id, lists the categories for a user.
// https://developers.pocketsmith.com/reference#get_users-id-categories
func (c *Client) ListCategoriesForUser(
	ctx context.Context,
	options *ListCategoriesForUserOptions,
) (categories Categories, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListCategoriesForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list categories.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/categories", options.UserID),
	}, &categories)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list categories: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return categories, nil
}

// ListCategoriesForAuthedUser, using the token attached to a client, lists the
// categories for the authed user.
func (c *Client) ListCategories(ctx context.Context) (Categories, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListCategoriesForAuthedUser")
	defer span.End()

	// list categories for authed user.
	return c.ListCategoriesForUser(newCtx, &ListCategoriesForUserOptions{UserID: c.authedUser.ID})
}

// GetCategoryByTitle ...
type GetCategoryByTitleOptions struct {
	Category string `validator:"required"`
}

// GetCategoryByTitleOptions ...
type GetCategoryByTitleForUserOptions struct {
	UserID int `validator:"required"`

	GetCategoryByTitleOptions
}

// GetCategoryByTitleForUser, using the given user id and category, returns the found
// category for a user.
func (c *Client) GetCategoryByTitleForUser(
	ctx context.Context,
	options *GetCategoryByTitleForUserOptions,
) (*Category, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetCategoryByTitleForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// get category by title.
	categories, err := c.ListCategoriesForUser(
		newCtx,
		&ListCategoriesForUserOptions{UserID: options.UserID},
	)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get categories by title: %v", err))
		span.RecordError(err)
		return nil, err
	}
	for _, c := range categories {
		if c.Title == options.Category {
			return &c, nil
		}
	}
	err = fmt.Errorf(
		"category with title %q for user %v doesn't exist",
		options.Category,
		options.UserID,
	)
	span.SetStatus(codes.Error, fmt.Sprintf("failed to find category by title: %v", err))
	span.RecordError(err)
	return nil, err
}

// GetCategoryByTitle, using the token attached to the client,
// returns the found category for the authed user.
func (c *Client) GetCategoryByTitle(
	ctx context.Context,
	options *GetCategoryByTitleOptions,
) (*Category, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetCategoryByTitle")
	defer span.End()

	// get category by title.
	return c.GetCategoryByTitleForUser(
		newCtx,
		&GetCategoryByTitleForUserOptions{
			UserID:                    c.authedUser.ID,
			GetCategoryByTitleOptions: *options,
		},
	)
}

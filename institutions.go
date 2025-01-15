package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Institutions represents a slice of Institution.
type Institutions []Institution

// CreateInstitutionOptions defines the options for creating an institutions in
// Pocketsmith, under the authed user.
type CreateInstitutionOptions struct {
	Title        string `json:"title"         validator:"required"`
	CurrencyCode string `json:"currency_code" validator:"required"`
}

// CreateInstitutionOptionsForUser defines the options for creating an
// institution for the given user in Pocketsmith, by the user id.
type CreateInstitutionOptionsForUser struct {
	UserID int `json:"-" validator:"required"`

	CreateInstitutionOptions
}

// CreateInstitutionForUser creates an institution for the given user in
// Pocketsmith, by the user id.
// https://developers.pocketsmith.com/reference/post_users-id-institutions-1.
func (c *Client) CreateInstitutionForUser(
	ctx context.Context,
	options *CreateInstitutionOptionsForUser,
) (institution *Institution, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateInstitutionForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// create institution.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/institutions", options.UserID),
		body:   options,
	}, &institution)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to create institution: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return institution, nil
}

// CreateInstitution creates an institution in Pocketsmith, under the authed user.
func (c *Client) CreateInstitution(
	ctx context.Context,
	options *CreateInstitutionOptions,
) (*Institution, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateInstitution")
	defer span.End()

	// create Institution for authed user.
	return c.CreateInstitutionForUser(
		newCtx,
		&CreateInstitutionOptionsForUser{
			UserID:                   c.authedUser.ID,
			CreateInstitutionOptions: *options,
		},
	)
}

// DeleteInstitutionOptions defines the options for deleteing an institution.
type DeleteInstitutionOptions struct {
	InstitutionID int `json:"-" validator:"required"`

	MergeIntoInstitutionID int `json:"merge_into_institution_id"`
}

// DeleteInstitution, using the given institution id, creates an institution.
// https://developers.pocketsmith.com/reference/delete_institutions-id-1.
func (c *Client) DeleteInstitution(
	ctx context.Context,
	options *DeleteInstitutionOptions,
) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteInstitution")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return err
	}

	// delete institution.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/institutions/%v", options.InstitutionID),
		body:   options,
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to delete institution: %v", err))
		span.RecordError(err)
		return err
	}
	return nil
}

// ListInstitutionsForUser ...
type ListInstitutionsForUser struct {
	UserID int `json:"-" validator:"required"`
}

// ListInstitutionsForUser, using the given user id, list the institutions for a user.
// https://developers.pocketsmith.com/reference#get_users-id-institutions
func (c *Client) ListInstitutionsForUser(
	ctx context.Context,
	options *ListInstitutionsForUser,
) (institutions Institutions, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListInstitutionsForUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// list institutions.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/institutions", options.UserID),
	}, &institutions)
	return institutions, err
}

// ListInstitutions, using the token attached to the client, lists
// the institutions for the authed user.
func (c *Client) ListInstitutions(
	ctx context.Context,
) ([]Institution, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListInstitutions")
	defer span.End()

	// list institutions for authed user.
	return c.ListInstitutionsForUser(newCtx, &ListInstitutionsForUser{UserID: c.authedUser.ID})
}

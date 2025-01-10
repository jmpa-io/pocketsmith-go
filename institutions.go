package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

// CreateInstitutionOptions defines the options for creating an institution.
type CreateInstitutionOptions struct {
	Title        string `json:"title"`
	CurrencyCode string `json:"currency_code"`
}

// CreateInstitution, using the given user id, creates an institution for a user.
// https://developers.pocketsmith.com/reference#post_users-id-institutions
func (c *Client) CreateInstitution(
	ctx context.Context,
	userId int,
	options *CreateInstitutionOptions,
) (institution *Institution, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateInstitution")
	defer span.End()

	// create institution.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/institutions", userId),
		body:   options,
	}, &institution)
	return institution, err
}

// CreateInstitutionForAuthedUser, using the token attached to the client,
// creates an institution for the authed user.
func (c *Client) CreateInstitutionForAuthedUser(
	ctx context.Context,
	options *CreateInstitutionOptions,
) (*Institution, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateInstitutionForAuthedUser")
	defer span.End()

	// create Institution for authed user.
	return c.CreateInstitution(newCtx, c.authedUser.ID, options)
}

// DeleteInstitutionOptions defines the options for deleteing an institution.
type DeleteInstitutionOptions struct {
	MergeIntoInstitutionId int
}

// DeleteInstitution, using the given institution id, creates an institution.
// https://developers.pocketsmith.com/reference/delete_institutions-id-1
func (c *Client) DeleteInstitution(
	ctx context.Context,
	institutionId int,
	options *DeleteInstitutionOptions,
) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteInstitution")
	defer span.End()

	// delete institution.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/institutions/%v", institutionId),
		body:   options,
	}, nil)
	return err
}

// ListInstitutions, using the given user id, list the institutions for a user.
// https://developers.pocketsmith.com/reference#get_users-id-institutions
func (c *Client) ListInstitutions(
	ctx context.Context,
	userId int,
) (institutions []Institution, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListInstitutions")
	defer span.End()

	// list institutions.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/institutions", userId),
	}, &institutions)
	return institutions, err
}

// ListInstitutionsForAuthedUser, using the token attached to the client, lists
// the institutions for the authed user.
func (c *Client) ListInstitutionsForAuthedUser(ctx context.Context) ([]Institution, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListInstitutionsForAuthedUser")
	defer span.End()

	// list institutions for authed user.
	return c.ListInstitutions(newCtx, c.authedUser.ID)
}

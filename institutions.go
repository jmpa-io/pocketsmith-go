package pocketsmith

import (
	"fmt"
	"net/http"
	"time"
)

// Institution defines a PocketSmith institution.
// https://developers.pocketsmith.com/reference#get_institutions-id
type Institution struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	CurrencyCode string    `json:"currency_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateInstitutionOptions defines the options for creating an institution.
type CreateInstitutionOptions struct {
	Title        string `json:"title"`
	CurrencyCode string `json:"currency_code"`
}

// CreateInstitution, using the given user id, creates an institution for a user.
// https://developers.pocketsmith.com/reference#post_users-id-institutions
func (c *Client) CreateInstitution(
	userId int,
	options *CreateInstitutionOptions,
) (*Institution, error) {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/users/%v/institutions", userId),
		data:   options,
	}
	var institution *Institution
	_, err := c.sender(cr, &institution)
	return institution, err
}

// CreateInstitutionForAuthedUser, using the token attached to the client,
// creates an institution for the authed user.
func (c *Client) CreateInstitutionForAuthedUser(
	options *CreateInstitutionOptions,
) (*Institution, error) {
	return c.CreateInstitution(c.user.ID, options)
}

// DeleteInstitutionOptions defines the options for deleteing an institution.
type DeleteInstitutionOptions struct {
	MergeIntoInstitutionId int
}

// DeleteInstitution, using the given institution id, creates an institution.
// https://developers.pocketsmith.com/reference/delete_institutions-id-1
func (c *Client) DeleteInstitution(institutionId int, options *DeleteInstitutionOptions) error {
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/institutions/%v", institutionId),
		data:   options,
	}
	_, err := c.sender(cr, nil)
	return err
}

// ListInstitutions, using the given user id, list the institutions for a user.
// https://developers.pocketsmith.com/reference#get_users-id-institutions
func (c *Client) ListInstitutions(userId int) ([]Institution, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/institutions", userId),
	}
	var institutions []Institution
	_, err := c.sender(cr, &institutions)
	return institutions, err
}

// ListInstitutionsForAuthedUser, using the token attached to the client, lists
// the institutions for the authed user.
func (c *Client) ListInstitutionsForAuthedUser() ([]Institution, error) {
	return c.ListInstitutions(c.user.ID)
}

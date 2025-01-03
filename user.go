package pocketsmith

import (
	"net/http"
	"time"
)

// User defines a PocketSmith user.
// https://developers.pocketsmith.com/reference#get_users-id
type User struct {
	ID                      int       `json:"id"`
	Login                   string    `json:"login"`
	Name                    string    `json:"name"`
	Email                   string    `json:"email"`
	AvatarURL               string    `json:"avatar_url"`
	TimeZone                string    `json:"time_zone"`
	WeekStartDay            int       `json:"week_start_day"`
	BaseCurrencyCode        string    `json:"base_currency_code"`
	AlwaysShowBaseCurrency  bool      `json:"always_show_base_currency"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	UsingMultipleCurrencies bool      `json:"using_multiple_currencies"`
	LastLoggedInAt          time.Time `json:"last_logged_in_at"`
	LastActivityAt          time.Time `json:"last_activity_at"`
}

// GetAuthedUser, using the token attached to the client, returns information
// about the authed user.
// https://developers.pocketsmith.com/reference#get_me
func (c *Client) GetAuthedUser() (*User, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   "/me",
	}
	var user *User
	_, err := c.sender(cr, &user)
	return user, err
}

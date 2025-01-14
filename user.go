package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// User defines a PocketSmith user.
type User struct {
	ID                      int    `json:"id"`
	Login                   string `json:"login"`
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	AvatarURL               string `json:"avatar_url"`
	BetaUser                bool   `json:"beta_user"`
	TimeZone                string `json:"time_zone"`
	WeekStartDay            int    `json:"week_start_day"`
	IsReviewingTransactions bool   `json:"is_reviewing_transactions"`
	BaseCurrencyCode        string `json:"base_currency_code"`
	AlwaysShowBaseCurrency  bool   `json:"always_show_base_currency"`
	UsingMultipleCurrencies bool   `json:"using_multiple_currencies"`

	AvailableAccounts int `json:"available_accounts"`
	AvailableBudgets  int `json:"available_budgets"`

	ForecastLastUpdatedAt    time.Time  `json:"forecast_last_updated_at"`
	ForecastLastAccessedAt   time.Time  `json:"forecast_last_accessed_at"`
	ForecastStartDate        customTime `json:"forecast_start_date"`
	ForecastEndDate          customTime `json:"forecast_end_date"`
	ForecastDeferRecalculate bool       `json:"forecast_defer_recalculate"`
	ForecastNeedsRecalculate bool       `json:"forecast_needs_recalculate"`

	LastLoggedInAt time.Time `json:"last_logged_in_at"`
	LastActivityAt time.Time `json:"last_activity_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAuthedUser returns data about the user who owns the token used by this client.
// https://developers.pocketsmith.com/reference/get_me-1.
func (c *Client) GetAuthedUser(ctx context.Context) (user *User, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetAuthedUser")
	defer span.End()

	// return the current authed user.
	if c.authedUser != nil {
		span.SetAttributes(attribute.Bool("cached", true))
		return c.authedUser, nil
	}

	// get authed user.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   "/me",
	}, &user)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get authed user: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return user, nil
}

// GetUserOptions ...
type GetUserOptions struct {
	ID int `json:"id" validator:"required"`
}

// GetUser returns data about a user, by their given user id.
// https://developers.pocketsmith.com/reference/get_users-id-1.
func (c *Client) GetUser(ctx context.Context, options *GetUserOptions) (user *User, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetUser")
	defer span.End()

	// validate options.
	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// get user.
	_, err = c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   "/users/id",
		body:   options,
	}, &user)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get user: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return user, nil
}

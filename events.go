package pocketsmith

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Event defines a Pocketsmith forecast event (budget entry).
type Event struct {
	ID               string   `json:"id"`
	SeriesID         int64    `json:"series_id"`
	SeriesStartID    string   `json:"series_start_id"`
	Category         Category `json:"category"`
	Scenario         Scenario `json:"scenario"`
	Amount           float64  `json:"amount"`
	AmountInBase     float64  `json:"amount_in_base_currency"`
	CurrencyCode     string   `json:"currency_code"`
	Date             string   `json:"date"`
	Colour           string   `json:"colour"`
	Note             string   `json:"note"`
	RepeatType       string   `json:"repeat_type"`
	RepeatInterval   int      `json:"repeat_interval"`
	InfiniteSeries   bool     `json:"infinite_series"`
	IsTransfer       bool     `json:"is_transfer"`
}

// Events represents a slice of Event.
type Events []Event

// ---
// List events
// ---

// ListEventsForUserOptions defines options for listing forecast events.
type ListEventsForUserOptions struct {
	UserID    int    `json:"-"          validator:"required"`
	StartDate string `json:"start_date" validator:"required"`
	EndDate   string `json:"end_date"   validator:"required"`
}

// ListEventsOptions defines options for listing forecast events for the authed user.
type ListEventsOptions struct {
	StartDate string `json:"start_date" validator:"required"`
	EndDate   string `json:"end_date"   validator:"required"`
}

// ListEventsForUser lists all forecast events for a user within a date range.
// https://developers.pocketsmith.com/reference/get_users-id-events-1
func (c *Client) ListEventsForUser(
	ctx context.Context,
	options *ListEventsForUserOptions,
) (Events, error) {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListEventsForUser")
	defer span.End()

	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// NOTE: the events API does not support per_page — pass dates directly in the path.
	var events Events
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path: fmt.Sprintf(
			"/users/%v/events?start_date=%s&end_date=%s",
			options.UserID, options.StartDate, options.EndDate,
		),
	}, &events)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list events: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return events, nil
}

// ListEvents lists all forecast events for the authed user within a date range.
func (c *Client) ListEvents(ctx context.Context, options *ListEventsOptions) (Events, error) {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListEvents")
	defer span.End()
	return c.ListEventsForUser(newCtx, &ListEventsForUserOptions{
		UserID:    c.authedUser.ID,
		StartDate: options.StartDate,
		EndDate:   options.EndDate,
	})
}

// ---
// Create event (budget entry)
// ---

// CreateEventOptions defines options for creating a forecast event.
type CreateEventOptions struct {
	ScenarioID     int     `json:"-"                validator:"required"`
	CategoryID     int32   `json:"category_id"      validator:"required"`
	Amount         float64 `json:"amount"           validator:"required"` // negative = expense, positive = income
	Date           string  `json:"date"             validator:"required"` // YYYY-MM-DD, start date of the series
	RepeatType     string  `json:"repeat_type"`     // once, daily, weekly, fortnightly, monthly, yearly
	RepeatInterval int     `json:"repeat_interval"` // e.g. 1 for every month
	InfiniteSeries bool    `json:"infinite_series"`
	Note           string  `json:"note,omitempty"`
	Colour         string  `json:"colour,omitempty"`
}

// CreateEvent creates a new recurring forecast event (budget entry) in a scenario.
// The API only returns the event ID on creation, so we follow up with a ListEvents
// call to return the full event.
// https://developers.pocketsmith.com/reference/post_scenarios-id-events-1
func (c *Client) CreateEvent(
	ctx context.Context,
	options *CreateEventOptions,
) (*Event, error) {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "CreateEvent")
	defer span.End()

	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	// default repeat to monthly if not specified.
	if options.RepeatType == "" {
		options.RepeatType = "monthly"
	}
	if options.RepeatInterval == 0 {
		options.RepeatInterval = 1
	}

	// the API only returns {"id": "..."} on creation.
	// Note: the event may not immediately appear in ListEvents due to Pocketsmith
	// processing delay — callers should use ListEvents after a moment to verify.
	var created struct {
		ID string `json:"id"`
	}
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/scenarios/%v/events", options.ScenarioID),
		body:   options,
	}, &created)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to create event: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return &Event{ID: created.ID}, nil
}

// ---
// Update event
// ---

// UpdateEventBehaviour controls which events in a series are updated.
type UpdateEventBehaviour string

const (
	UpdateEventBehaviourOne     UpdateEventBehaviour = "one"     // only this occurrence
	UpdateEventBehaviourForward UpdateEventBehaviour = "forward" // this and all future occurrences
	UpdateEventBehaviourAll     UpdateEventBehaviour = "all"     // entire series
)

// UpdateEventOptions defines options for updating a forecast event.
type UpdateEventOptions struct {
	EventID    string               `json:"-"          validator:"required"`
	Behaviour  UpdateEventBehaviour `json:"behaviour"  validator:"required"`
	Amount     *float64             `json:"amount,omitempty"`
	Date       string               `json:"date,omitempty"`
	Note       string               `json:"note,omitempty"`
	CategoryID int32                `json:"category_id,omitempty"`
	RepeatType string               `json:"repeat_type,omitempty"`
}

// UpdateEvent updates a forecast event (or series of events).
// The API only returns {"id": "..."} on update.
// https://developers.pocketsmith.com/reference/put_events-id-1
func (c *Client) UpdateEvent(
	ctx context.Context,
	options *UpdateEventOptions,
) (*Event, error) {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "UpdateEvent")
	defer span.End()

	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	var result struct {
		ID string `json:"id"`
	}
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPut,
		path:   fmt.Sprintf("/events/%v", options.EventID),
		body:   options,
	}, &result)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to update event: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return &Event{ID: result.ID}, nil
}

// ---
// Delete event
// ---

// DeleteEventBehaviour controls which events in a series are deleted.
type DeleteEventBehaviour string

const (
	DeleteEventBehaviourOne     DeleteEventBehaviour = "one"
	DeleteEventBehaviourForward DeleteEventBehaviour = "forward"
	DeleteEventBehaviourAll     DeleteEventBehaviour = "all"
)

// DeleteEventOptions defines options for deleting a forecast event.
type DeleteEventOptions struct {
	EventID   string               `json:"-"         validator:"required"`
	Behaviour DeleteEventBehaviour `json:"behaviour" validator:"required"`
}

// DeleteEvent deletes a forecast event (or series).
// https://developers.pocketsmith.com/reference/delete_events-id-1
func (c *Client) DeleteEvent(
	ctx context.Context,
	options *DeleteEventOptions,
) error {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "DeleteEvent")
	defer span.End()

	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return err
	}

	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/events/%v", options.EventID),
		body:   options,
	}, nil)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to delete event: %v", err))
		span.RecordError(err)
		return err
	}
	return nil
}

// ---
// Get budget summary
// ---

// BudgetPeriod defines a single period in a budget summary.
type BudgetPeriod struct {
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	CurrencyCode   string  `json:"currency_code"`
	ForecastAmount float64 `json:"forecast_amount"`
	ActualAmount   float64 `json:"actual_amount"`
	RefundAmount   float64 `json:"refund_amount"`
	Current        bool    `json:"current"`
	OverBudget     bool    `json:"over_budget"`
	UnderBudget    bool    `json:"under_budget"`
	OverBy         float64 `json:"over_by"`
	UnderBy        float64 `json:"under_by"`
	PercentageUsed float64 `json:"percentage_used"`
}

// BudgetEntry defines the budget data for a single category direction (expense or income).
type BudgetEntry struct {
	StartDate            string         `json:"start_date"`
	EndDate              string         `json:"end_date"`
	CurrencyCode         string         `json:"currency_code"`
	TotalActualAmount    float64        `json:"total_actual_amount"`
	AverageActualAmount  float64        `json:"average_actual_amount"`
	TotalForecastAmount  float64        `json:"total_forecast_amount"`
	AverageForecastAmount float64       `json:"average_forecast_amount"`
	TotalUnderBy         float64        `json:"total_under_by"`
	TotalOverBy          float64        `json:"total_over_by"`
	Periods              []BudgetPeriod `json:"periods"`
}

// BudgetItem defines a single category's budget (expense and/or income).
type BudgetItem struct {
	Category   Category     `json:"category"`
	IsTransfer bool         `json:"is_transfer"`
	Expense    *BudgetEntry `json:"expense"`
	Income     *BudgetEntry `json:"income"`
}

// GetBudgetOptions defines options for getting the current budget.
type GetBudgetOptions struct{}

// GetBudget returns the current budget for all categories for the authed user.
// https://developers.pocketsmith.com/reference/get_users-id-budget-1
func (c *Client) GetBudget(ctx context.Context) ([]BudgetItem, error) {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetBudget")
	defer span.End()

	var items []BudgetItem
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/budget", c.authedUser.ID),
	}, &items)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get budget: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return items, nil
}

// GetBudgetSummaryOptions defines options for getting a budget summary over a date range.
type GetBudgetSummaryOptions struct {
	StartDate string `json:"start_date" validator:"required"`
	EndDate   string `json:"end_date"   validator:"required"`
	Period    string `json:"period"`   // months, weeks, years, days
	Interval  int    `json:"interval"` // e.g. 1
	RollUp    bool   `json:"roll_up"`
}

// BudgetSummary holds the income and expense budget summary over a date range.
type BudgetSummary struct {
	Income  *BudgetEntry `json:"income"`
	Expense *BudgetEntry `json:"expense"`
}

// GetBudgetSummary returns a rolled-up budget summary over a date range.
// https://developers.pocketsmith.com/reference/get_users-id-budget-summary-1
func (c *Client) GetBudgetSummary(
	ctx context.Context,
	options *GetBudgetSummaryOptions,
) (*BudgetSummary, error) {
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetBudgetSummary")
	defer span.End()

	if err := c.validator.StructCtx(newCtx, options); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to validate options: %v", err))
		span.RecordError(err)
		return nil, err
	}

	period := options.Period
	if period == "" {
		period = "months"
	}
	interval := options.Interval
	if interval == 0 {
		interval = 1
	}

	var summary BudgetSummary
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/users/%v/budget_summary", c.authedUser.ID),
		queries: setupQueries(&map[string]string{
			"start_date": options.StartDate,
			"end_date":   options.EndDate,
			"period":     period,
			"interval":   fmt.Sprintf("%v", interval),
			"roll_up":    fmt.Sprintf("%v", options.RollUp),
		}),
	}, &summary)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get budget summary: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return &summary, nil
}

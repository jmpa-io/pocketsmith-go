package pocketsmith

import (
	"context"
	"encoding/json"
	"fmt"
)

func toJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		safe, _ := json.Marshal(err.Error())
		return `{"error":` + string(safe) + `}`
	}
	return string(b)
}

func (c *Client) GetAuthedUserString(ctx context.Context) (string, error) {
	user, err := c.GetAuthedUser(ctx)
	if err != nil {
		return "", err
	}
	return toJSON(user), nil
}

func (c *Client) ListAccountsString(ctx context.Context) (string, error) {
	accounts, err := c.ListAccounts(ctx)
	if err != nil {
		return "", err
	}
	return toJSON(accounts), nil
}

func (c *Client) ListTransactionAccountsString(ctx context.Context) (string, error) {
	accounts, err := c.ListTransactionAccounts(ctx)
	if err != nil {
		return "", err
	}
	return toJSON(accounts), nil
}

func (c *Client) ListTransactionAccountTransactionsString(ctx context.Context, opts *ListTransactionAccountTransactionsOptions) (string, error) {
	txns, err := c.ListTransactionAccountTransactions(ctx, opts)
	if err != nil {
		return "", err
	}
	return toJSON(txns), nil
}

func (c *Client) GetTransactionString(ctx context.Context, opts *GetTransactionOptions) (string, error) {
	t, err := c.GetTransaction(ctx, opts)
	if err != nil {
		return "", err
	}
	return toJSON(t), nil
}

func (c *Client) ListCategoriesString(ctx context.Context) (string, error) {
	cats, err := c.ListCategories(ctx)
	if err != nil {
		return "", err
	}
	return toJSON(cats), nil
}

func (c *Client) GetBudgetString(ctx context.Context) (string, error) {
	items, err := c.GetBudget(ctx)
	if err != nil {
		return "", err
	}
	return toJSON(items), nil
}

func (c *Client) ListEventsString(ctx context.Context, opts *ListEventsOptions) (string, error) {
	events, err := c.ListEvents(ctx, opts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Events (%d):\n%s", len(events), toJSON(events)), nil
}

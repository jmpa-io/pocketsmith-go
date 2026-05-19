package pocketsmith

import "log/slog"

// Option configures a departure client.
type Option func(*Client) error

// WithLogLevel sets the log level for the default logger.
func WithLogLevel(level slog.Level) Option {
	return func(c *Client) error {
		c.logLevel = level
		return nil
	}
}

// WithLogger overwrites the default logger with the given custom logger.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// WithHttpClient overwrites the default httpClient used for API communication.
func WithHttpClient(httpClient iHttpClient) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithSkipAuthCheck disables the GetAuthedUser call that normally happens at
// client startup. Use this when the PocketSmith API may be unreachable at
// init time (e.g. restrictive corporate networks). Tools will still fail
// gracefully if the API is unreachable when they are called.
func WithSkipAuthCheck() Option {
	return func(c *Client) error {
		c.skipAuthCheck = true
		return nil
	}
}

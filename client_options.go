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

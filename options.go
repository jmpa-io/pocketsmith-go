package pocketsmith

// Option configures a departure client.
type Option func(*Client) error

// SetLogLevel sets the log level for the client.
func SetLogLevel(level LogLevel) Option {
	return func(c *Client) error {
		c.logLevel = level
		return nil
	}
}

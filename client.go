package pocketsmith

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	// TODO add this log level: https://github.com/rs/zerolog
)

// An iHttpClient is an interface over http.Client.
type iHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines a PocketSmith client; the interface for this package.
type Client struct {
	logLevel   LogLevel
	httpClient iHttpClient
	headers    http.Header
	endpoint   string // the endpoint to query against.

	// misc.
	logger zerolog.Logger

	// metadata.
	user *User // the authed user attached to the token.
}

// New returns a client for this package, which can be used to make
// requests to the PocketSmith API.
func New(token string, options ...Option) (*Client, error) {

	// check args.
	if token == "" {
		return nil, ErrMissingToken{}
	}

	// default client.
	c := &Client{
		logLevel:   LogLevelDebug,
		httpClient: http.DefaultClient,
		endpoint:   "https://api.pocketsmith.com/v2",
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, ErrFailedOptionSet{err}
		}
	}

	// setup logger.
	zerolog.MessageFieldName = "msg"
	var level zerolog.Level
	switch c.logLevel {
	case LogLevelDebug:
		level = zerolog.DebugLevel
	case LogLevelInfo:
		level = zerolog.InfoLevel
	case LogLevelWarn:
		level = zerolog.WarnLevel
	case LogLevelError:
		level = zerolog.ErrorLevel
	default:
		level = zerolog.ErrorLevel
	}
	c.logger = zerolog.New(os.Stderr).
		With().Caller().Logger().
		With().Timestamp().Logger().Level(level)
	c.logger.Debug().Msg("setting up client")

	// setup headers.
	headers := make(http.Header)
	headers.Set("X-Developer-Key", token)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json") // TODO is this needed?
	headers.Set("Accept-Charset", "utf-8")    // TODO is this needed?
	c.headers = headers

	// retrieve authed user, to determine if the token is valid.
	user, err := c.GetAuthedUser()
	if err != nil {
		return nil, err
	}
	c.user = user

	c.logger.Debug().Msg("client setup successfully")
	return c, nil
}

package pocketsmith

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
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
	logger log.Logger

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
		logLevel:   LogLevelError,
		httpClient: http.DefaultClient,
		endpoint:   "https://api.pocketsmith.com/v2",
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, ErrFailedOptionSet{err}
		}
	}

	// retrieve authed user, to determine if the token is valid.
	user, err := c.GetAuthedUser()
	if err != nil {
		return nil, err
	}
	c.user = user

	// setup headers.
	headers := make(http.Header)
	headers.Set("X-Developer-Key", token)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json") // TODO is this needed?
	headers.Set("Accept-Charset", "utf-8")    // TODO is this needed?
	c.headers = headers

	// setup logger.
	config := promlog.AllowedLevel{}
	if err := config.Set(c.logLevel.String()); err != nil {
		return nil, ErrFailedLoggerSetup{err}
	}
	c.logger = promlog.New(&promlog.Config{Level: &config})

	_ = level.Debug(c.logger).Log("msg", "setup client successfully")
	return c, nil
}
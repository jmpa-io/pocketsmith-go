package pocketsmith

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
)

// An iHttpClient is an interface over http.Client.
type iHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines a client for this package.
type Client struct {

	// tracing.
	tracerName string // The name of the tracer output in the traces.

	// config.
	endpoint      string      // The endpoint to query against.
	httpClient    iHttpClient // The http client used when sending / receiving data from the endpoint.
	headers       http.Header // The headers passed to the http client when sending / receiving data from the endpoint.
	skipAuthCheck bool        // Skip the GetAuthedUser call on startup (useful when API is unreachable).

	// misc.
	logLevel  slog.Level          // The log level of the default logger.
	logger    *slog.Logger        // The logger used in this client (custom or default).
	validator *validator.Validate // A validator for validating structs.

	// metadata.
	authedUser *User // the authed user attached to the token.
}

// New creates and returns a new Client, initialized with the provided token.
// The client itself is set up with tracing, logging, and HTTP configuration.
// Additional options can be provided to modify its behavior, via the options
// slice. The client is used for making requests and interacting with the
// Pockestmith API.
func New(ctx context.Context, token string, options ...Option) (*Client, error) {

	// setup tracing.
	tracerName := "pocketsmith-go"
	newCtx, span := otel.Tracer(tracerName).Start(ctx, "New")
	defer span.End()

	// check args.
	if token == "" {
		return nil, ErrClientEmptyToken{}
	}

	// default client — clone http.DefaultTransport so we inherit all of its
	// system-level settings (DNS resolver, TLS config, proxy env vars, etc.)
	// and only override the connection-pool limits we care about:
	//   • a 30-second request timeout (DefaultClient has none — hangs forever)
	//   • MaxIdleConnsPerHost=20 so the concurrent page workers in fetchAllPages
	//     can all reuse keep-alive connections to the same API host instead of
	//     being throttled to DefaultTransport's limit of 2
	dt := http.DefaultTransport.(*http.Transport).Clone()
	dt.MaxIdleConns = 100
	dt.MaxIdleConnsPerHost = 20
	dt.IdleConnTimeout = 90 * time.Second
	c := &Client{
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: dt,
		},
		endpoint: "https://api.pocketsmith.com/v2",
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, ErrClientFailedToSetOption{err}
		}
	}

	// determine if the default logger should be used.
	if c.logger == nil {

		// use default logger.
		c.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: c.logLevel, // default log level is 'INFO'.
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
				}
				return a
			},
		}))

	}

	// setup validator.
	c.validator = validator.New(validator.WithRequiredStructEnabled())

	// setup headers.
	headers := make(http.Header)
	headers.Set("X-Developer-Key", token)
	headers.Set("Content-Type", "application/json")
	c.headers = headers

	// retrieve authed user, to determine if the token is valid.
	// Skipped if WithSkipAuthCheck was used — useful when the API is
	// unreachable at startup (e.g. restrictive corporate networks).
	if !c.skipAuthCheck {
		var err error
		if c.authedUser, err = c.GetAuthedUser(newCtx); err != nil {
			return nil, ErrClientFailedToGetAuthedUser{err}
		}
	}

	c.logger.Debug("client setup successfully")
	return c, nil
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"

	"github.com/jmpa-io/pocketsmith-go"
)

type handler struct {

	// config.
	name        string
	version     string
	environment string

	// clients.
	pocketsmithsvc *pocketsmith.Client

	// misc.
	logger *slog.Logger
}

// run is like main but after the handler is configured.
func (h *handler) run(ctx context.Context) {

	// setup span.
	newCtx, span := otel.Tracer(h.name).Start(ctx, "run")
	defer span.End()

	// get authed user.
	u, err := h.pocketsmithsvc.GetAuthedUser(newCtx)
	if err != nil {
		fmt.Printf("failed to get authed user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", u)

}

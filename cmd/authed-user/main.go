package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmpa-io/pocketsmith-go"
)

func main() {

	// setup tracing.
	ctx := context.TODO()

	// retrieve token.
	token := os.Getenv("POCKETSMITH_TOKEN")

	// setup client.
	c, err := pocketsmith.New(ctx, token, pocketsmith.WithLogLevel(slog.LevelWarn))
	if err != nil {
		fmt.Printf("failed to setup client: %v\n", err)
		os.Exit(1)
	}

	// get authed user.
	u, err := c.GetAuthedUser(ctx)
	if err != nil {
		fmt.Printf("failed to get authed user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", u)
}

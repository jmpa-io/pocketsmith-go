package main

import (
	"fmt"
	"os"

	"github.com/jmpa-io/pocketsmith-go"
)

func main() {

	// setup client.
	c, err := pocketsmith.New("xxxx")
	if err != nil {
		fmt.Printf("failed to setup client: %v\n", err)
		os.Exit(1)
	}

	// do something with the client..
	// like retrieve the authed user attached to the token.
	u, err := c.GetAuthedUser()
	if err != nil {
		fmt.Printf("failed to get authed user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", u)
}

//go:build ignore

// gen downloads an updated version of the PSL list and compiles it into go code.
//
// It is meant to be used by maintainers in conjunction with the go generate tool
// to update the list.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/nerdlem/publicsuffix-go/publicsuffix/generator"
)

const (
	// where the rules will be written
	filename = "rules.go"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	g := generator.NewGenerator()
	g.Verbose = true
	err := g.Write(ctx, filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

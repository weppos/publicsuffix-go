// +build ignore

// gen downloads an updated version of the PSL list and compiles it into go code.
//
// It is meant to be used by maintainers in conjunction with the go generate tool
// to update the list.
package main

import (
	"github.com/weppos/publicsuffix-go/publicsuffix/generator"
)

const (
	// where the rules will be written
	filename = "rules.go"
)

func main() {
	g := generator.NewGenerator()
	g.Verbose = true
	g.Write(filename)
}

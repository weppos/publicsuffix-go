// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func main() {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		elapsed -= elapsed % 1000000
		fmt.Printf("Time elapsed: %s\n", elapsed)
	}()

	fmt.Printf("%d rules loaded\n", publicsuffix.DefaultList.Size())
}

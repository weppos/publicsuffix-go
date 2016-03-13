package main

// +build ignore

import (
	"net/http"
	"go/format"
	"fmt"
	"os"
	"bytes"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

var (

)

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	list := publicsuffix.NewList()
	list.Load(resp.Body, nil)

	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "// This file is generated automatically by `go run gen.go`\n")
	fmt.Fprintf(buf, "// DO NOT EDIT MANUALLY\n\n")
	fmt.Fprintf(buf, "package publicsuffix\n\n")

	fmt.Fprintf(buf, `
import (
	"strconv"
	"strings"
	"fmt"
)

const defaultListVersion = "HEAD"

func initDefaultList() {
        rules := `)

	fmt.Fprintf(buf, "`")
	for _, rule := range list.Rules() {
		private := 0
		if rule.Private {
			private = 1
		}
		fmt.Fprintf(buf, fmt.Sprintf("\n%v,%v,%v,%v", rule.Type, rule.Value, rule.Length, private))
	}

	//for _, rule := range list.Rules() {
	//	fmt.Fprintf(buf,
	//		fmt.Sprintf(`
	//l.AddRule(&Rule{%d, "%s", %d, %v})`,
	//			rule.Type, rule.Value, rule.Length, rule.Private))
	//}
	fmt.Fprintf(buf, "`")
	fmt.Fprintf(buf, `

	for _, rule := range strings.Split(rules, "\n") {
		if len(rule) > 0 {
                        tokens := strings.Split(rule, ",")
                        t, _ := strconv.Atoi(tokens[0])
                        l, _ := strconv.Atoi(tokens[2])
                        v := tokens[1]
                        p := (tokens[3] == "1")
                        fmt.Println(&Rule{t, v, l, p})
			defaultList.AddRule(&Rule{t, v, l, p})
		}
	}
}
`)

	b, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = os.Stdout.Write(b)
}

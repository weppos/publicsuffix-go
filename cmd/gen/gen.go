package main

// +build ignore

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func main() {
	sha, datetime := extractHeadInfo()

	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/publicsuffix/list/%s/public_suffix_list.dat", sha))
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()

	list := publicsuffix.NewList()
	rules, _ := list.Load(resp.Body, nil)
	if err != nil {
		fatal(err)
	}

	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "// This file is generated automatically by `go run gen.go`\n")
	fmt.Fprintf(buf, "// DO NOT EDIT MANUALLY\n\n")
	fmt.Fprintf(buf, "package publicsuffix\n\n")

	fmt.Fprintf(buf, fmt.Sprintf(`
import (
	"strconv"
	"strings"
	"fmt"
)

const defaultListVersion = "PSL version %s (%v)"

func initDefaultList() {
        rules := `, sha, datetime.Format(time.ANSIC)))

	fmt.Fprintf(buf, "`")
	for _, rule := range rules {
		private := 0
		if rule.Private {
			private = 1
		}
		fmt.Fprintf(buf, fmt.Sprintf("\n%v,%v,%v,%v", rule.Type, rule.Value, rule.Length, private))
	}

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

// Is there a better way?
func extractHeadInfo() (sha string, datetime time.Time) {
	var re *regexp.Regexp
	resp, err := http.Get("https://github.com/publicsuffix/list")
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatal(err)
	}

	re = regexp.MustCompile(`<a class="commit-tease-sha" (?:.+)>([\w\s\n]+)<\/a>`)
	sha = strings.TrimSpace(re.FindStringSubmatch(string(data[:]))[1])
	if sha == "" {
		fatal(fmt.Errorf("sha is blank"))
	}

	re = regexp.MustCompile(`<span itemprop="dateModified">(?:.+)datetime="([\w\d\:\-]+)"(?:.+)</span>`)
	stringtime := re.FindStringSubmatch(string(data[:]))[1]
	if stringtime == "" {
		fatal(fmt.Errorf("date is blank"))
	}
	datetime, err = time.Parse(time.RFC3339, stringtime)
	if err != nil {
		fatal(err)
	}

	return sha, datetime
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

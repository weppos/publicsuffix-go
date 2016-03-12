// Package publicsuffix provides a domain name parser
// based on data from the public suffix list http://publicsuffix.org/.
// A public suffix is one under which Internet users can directly register names.
package publicsuffix

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"net/http/cookiejar"
)

const (
	NormalType    = 1
	WildcardType  = 2
	ExceptionType = 3

	listTokenPrivateDomains = "===BEGIN PRIVATE DOMAINS==="
	listTokenComment        = "//"

	defaultListFile = "list.dat"
)

// defaultList is the default List and it is used by Parse and Domain.
var defaultList = NewList()

// DefaultRule is the default Rule that represents "*".
var DefaultRule = NewRule("*")

// DefaultParserOptions are the default options used to parse a Public Suffix list.
var DefaultParserOptions = &ParserOption{PrivateDomains: true}

// DefaultFindOptions are the default options used to perform the lookup of rules in the list.
var DefaultFindOptions = &FindOptions{IgnorePrivate: false}

// Rule represents a single rule in a Public Suffix List.
type Rule struct {
	Type    int
	Value   string
	Length  int
	Private bool
}

// ParserOption are the options you can use to customize the way a List
// is parsed from a file or a string.
type ParserOption struct {
	PrivateDomains bool
}

// FindOptions are the options you can use to customize the way a Rule
// is searched within the list.
type FindOptions struct {
	IgnorePrivate bool
}

// List represents a Public Suffix List.
type List struct {
	// rules is kept private because you should not access rules directly
	// for lookup optimization the list will not be guaranteed to be a simple slice forever
	rules []Rule
}

// NewList creates a new empty list.
func NewList() *List {
	return &List{}
}

// NewListFromString parses a string that represents a Public Suffix source
// and returns a List initialized with the rules in the source.
func NewListFromString(src string, options *ParserOption) (*List, error) {
	l := NewList()
	return l, l.Load(src, options)
}

// NewListFromString parses a string that represents a Public Suffix source
// and returns a List initialized with the rules in the source.
func NewListFromFile(path string, options *ParserOption) (*List, error) {
	l := NewList()
	return l, l.LoadFile(path, options)
}

// experimental
func (l *List) Load(src string, options *ParserOption) error {
	r := strings.NewReader(src)
	return l.parse(r, options)
}

// experimental
func (l *List) LoadFile(path string, options *ParserOption) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return l.parse(f, options)
}

// AddRule adds a new rule to the list.
//
// The exact position of the rule into the list is unpredictable.
// The list may be optimized internally for lookups, therefore the algorithm
// will decide the best position for the new rule.
func (l *List) AddRule(r *Rule) error {
	l.rules = append(l.rules, *r)
	return nil
}

// experimental
func (l *List) Size() int {
	return len(l.rules)
}

func (l *List) parse(r io.Reader, options *ParserOption) error {
	if options == nil {
		options = DefaultParserOptions
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var section int // 1 == ICANN, 2 == PRIVATE

Scanning:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {

		// skip blank lines
		case line == "":
			break

		// include private domains or stop scanner
		case strings.Contains(line, listTokenPrivateDomains):
			if !options.PrivateDomains {
				break Scanning
			}
			section = 2

		// skip comments
		case strings.HasPrefix(line, listTokenComment):
			break

		default:
			rule := NewRule(line)
			rule.Private = (section == 2)
			l.AddRule(rule)
		}

	}

	return nil
}

// Finds and returns the most appropriate rule for the domain name.
func (l *List) Find(name string, options *FindOptions) Rule {
	var rule *Rule

	for _, r := range l.Select(name, options) {
		if r.Type == ExceptionType {
			return r
		}
		if rule == nil || rule.Length < r.Length {
			rule = &r
		}
	}

	if rule != nil {
		return *rule
	}

	return *DefaultRule
}

// experimental
func (l *List) Select(name string, options *FindOptions) []Rule {
	var found []Rule

	if options == nil {
		options = DefaultFindOptions
	}

	// In this phase the search is a simple sequential scan
	for _, rule := range l.rules {
		if !rule.Match(name) {
			continue
		}
		if options.IgnorePrivate == true && rule.Private {
			continue
		}
		found = append(found, rule)
	}

	return found
}

// NewRule parses the rule content, creates and returns a Rule.
func NewRule(content string) *Rule {
	var rule *Rule
	var value string

	switch content[0:1] {
	case "*": // wildcard
		if content == "*" {
			value = ""
		} else {
			value = content[2:len(content)]
		}
		rule = &Rule{Type: WildcardType, Value: value, Length: len(Labels(value)) + 1}
	case "!": // exception
		value = content[1:len(content)]
		rule = &Rule{Type: ExceptionType, Value: value, Length: len(Labels(value))}
	default: // normal
		value = content
		rule = &Rule{Type: NormalType, Value: value, Length: len(Labels(value))}
	}
	return rule
}

// Match checks if the rule matches the name.
//
// A domain name is said to match a rule if and only if all of the following conditions are met:
// - When the domain and rule are split into corresponding labels,
//   that the domain contains as many or more labels than the rule.
// - Beginning with the right-most labels of both the domain and the rule,
//   and continuing for all labels in the rule, one finds that for every pair,
//   either they are identical, or that the label from the rule is "*".
//
// See https://publicsuffix.org/list/
func (r *Rule) Match(name string) bool {
	left := strings.TrimSuffix(name, r.Value)

	// the name contains as many labels than the rule
	// this is a match, unless it's a wildcard
	// because the wildcard requires one more label
	if left == "" {
		return r.Type != WildcardType
	}

	// if there is one more label, the rule match
	// because either the rule is shorter than the domain
	// or the rule is a wildcard and there is one more label
	return left[len(left)-1:] == "."
}

// Decompose takes a name as input and decomposes it into a tuple of <TRD+SLD, TLD>,
// according to the rule definition and type.
func (r *Rule) Decompose(name string) [2]string {
	var parts []string

	switch r.Type {
	case WildcardType:
		parts = []string{}
		parts = append(parts, `.*?`)
		parts = append(parts, r.parts()...)
	default:
		parts = r.parts()
	}

	suffix := strings.Join(parts, `\.`)
	re := regexp.MustCompile(fmt.Sprintf(`^(.+)\.(%s)$`, suffix))

	matches := re.FindStringSubmatch(name)
	if len(matches) < 3 {
		return [2]string{"", ""}
	}

	return [2]string{matches[1], matches[2]}
}

func (r *Rule) parts() []string {
	labels := Labels(r.Value)
	if r.Type == ExceptionType {
		return labels[1:]
	}
	if r.Type == WildcardType && r.Value == "" {
		return []string{}
	}
	return labels
}

// Labels decomposes given domain name into labels,
// corresponding to the dot-separated tokens.
func Labels(name string) []string {
	return strings.Split(name, ".")
}

// DomainName represents a domain name.
type DomainName struct {
	Tld  string
	Sld  string
	Trd  string
	Rule *Rule
}

// String joins the components of the domain name into a single string.
// Empty labels are skipped.
//
// Examples:
//
// 	DomainName{"com", "example"}.String()
//	// example.com
// 	DomainName{"com", "example", "www"}.String()
//	// www.example.com
//
func (d *DomainName) String() string {
	switch {
	case d.Tld == "":
		return ""
	case d.Sld == "":
		return d.Tld
	case d.Trd == "":
		return d.Sld + "." + d.Tld
	default:
		return d.Trd + "." + d.Sld + "." + d.Tld
	}
}

// Domain extract and return the domain name from the input
// using the default (Public Suffix) List.
//
// Examples:
//
// 	publicsuffix.Domain("example.com")
//	// example.com
// 	publicsuffix.Domain("www.example.com)
//	// example.com
// 	publicsuffix.Domain("www.example.co.uk")
//	// example.co.uk
//
func Domain(name string) (string, error) {
	return DomainFromListWithOptions(DefaultList(), name, DefaultFindOptions)
}

// Parse decomposes the name into TLD, SLD, TRD
// using the default (Public Suffix) List,
// and returns the result as a DomainName
//
// Examples:
//
//	list := NewList()
//
// 	publicsuffix.Parse("example.com")
//	// &DomainName{"com", "example"}
// 	publicsuffix.Parse("www.example.com)
//	// &DomainName{"com", "example", "www"}
// 	publicsuffix.Parse("www.example.co.uk")
//	// &DomainName{"co.uk", "example"}
//
func Parse(name string) (*DomainName, error) {
	return ParseFromListWithOptions(DefaultList(), name, DefaultFindOptions)
}

// Ldomain extract and return the domain name from the input
// using the (Public Suffix) list passed as argument.
//
// Examples:
//
//	list := NewList()
//
// 	publicsuffix.DomainFromListWithOptions(list, "example.com")
//	// example.com
// 	publicsuffix.DomainFromListWithOptions(list, "www.example.com)
//	// example.com
// 	publicsuffix.DomainFromListWithOptions(list, "www.example.co.uk")
//	// example.co.uk
//
func DomainFromListWithOptions(l *List, name string, options *FindOptions) (string, error) {
	dn, err := ParseFromListWithOptions(l, name, options)
	if err != nil {
		return "", err
	}

	return dn.Sld + "." + dn.Tld, nil
}

// Lparse decomposes the name into TLD, SLD, TRD
// using the (Public Suffix) list passed as argument,
// and returns the result as a DomainName
//
// Examples:
//
//	list := NewList()
//
// 	publicsuffix.ParseFromListWithOptions(list, "example.com")
//	// &DomainName{"com", "example"}
// 	publicsuffix.ParseFromListWithOptions(list, "www.example.com)
//	// &DomainName{"com", "example", "www"}
// 	publicsuffix.ParseFromListWithOptions(list, "www.example.co.uk")
//	// &DomainName{"co.uk", "example"}
//
func ParseFromListWithOptions(l *List, name string, options *FindOptions) (*DomainName, error) {
	n, err := normalize(name)
	if err != nil {
		return nil, err
	}

	r := l.Find(n, options)
	if tld := r.Decompose(n)[1]; tld == "" {
		return nil, fmt.Errorf("%s is a suffix", n)
	}

	dn := &DomainName{Rule: &r}
	dn.Tld, dn.Sld, dn.Trd = decompose(&r, n)
	return dn, nil
}

// DefaultList returns the default list, initialized with the rules stored in the list.
// The list is lazy-initialized the first time it is requested.
func DefaultList() *List {
	if defaultList.Size() == 0 {
		// build the path to the file using the current directory
		// otherwise when this lib is loaded in a different package
		// the initializer crashes because it can't find the file
		_, filename, _, _ := runtime.Caller(0)
		filepath := path.Join(path.Dir(filename), defaultListFile)

		err := defaultList.LoadFile(filepath, DefaultParserOptions)
		if err != nil {
			panic(err)
		}
	}

	return defaultList
}

func normalize(name string) (string, error) {
	ret := strings.ToLower(name)

	if ret == "" {
		return "", fmt.Errorf("Name is blank")
	}
	if string(ret[0]) == "." {
		return "", fmt.Errorf("Name %s starts with a dot", ret)
	}

	return ret, nil
}

func decompose(r *Rule, name string) (tld, sld, trd string) {
	parts := r.Decompose(name)
	left, tld := parts[0], parts[1]

	dot := strings.LastIndex(left, ".")
	if dot == -1 {
		sld = left
		trd = ""
	} else {
		sld = left[dot+1:]
		trd = left[0:dot]
	}

	return
}

// CookieList implements the cookiejar.PublicSuffixList interface.
var CookieJarList cookiejar.PublicSuffixList = cookiejarList{}

type cookiejarList struct {
	List *List
}

// PublicSuffix implements cookiejar.PublicSuffixList.
func (l cookiejarList) PublicSuffix(domain string) string {
	if l.List == nil {
		l.List = DefaultList()
	}

	rule := l.List.Find(domain, nil)
	return rule.Decompose(domain)[1]
}

// PublicSuffix implements cookiejar.String.
func (cookiejarList) String() string {
	return "github.com/weppos/publicsuffix-go/publicsuffix"
}

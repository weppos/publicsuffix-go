package publicsuffix

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	NormalType    = 1
	WildcardType  = 2
	ExceptionType = 3

	listTokenPrivateDomains = "===BEGIN PRIVATE DOMAINS==="
	listTokenComment        = "//"

	defaultListFile = "list.dat"
)

// DefaultList is the default List and is used by Parse and Domain.
var DefaultList = NewList()

// DefaultRule is the default Rule that represents "*".
var DefaultRule = NewRule("*")

// DefaultParserOptions are the default options used to parse a Public Suffix list.
var DefaultParserOptions = &ParserOption{PrivateDomains: true}

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
func (l *List) Find(name string) Rule {
	var rule *Rule

	for _, r := range l.Select(name) {
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
func (l *List) Select(name string) []Rule {
	var found []Rule

	// In this phase the search is a simple sequential scan
	for _, rule := range l.rules {
		if rule.Match(name) {
			found = append(found, rule)
		}
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
	Tld string
	Sld string
	Trd string
}

// String joins the components of the domain name into a single string.
// Empty labels are skipped.
//
// Example:
// 	DomainName{"com", "example"}.String()
//	// "example.com"
// 	DomainName{"com", "example", "www"}.String()
//	// "www.example.com"
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

func Domain(name string) (string, error) {
	if DefaultList.Size() == 0 {
		initDefaultList()
	}

	return Ldomain(DefaultList, name)
}

func Parse(name string) (*DomainName, error) {
	if DefaultList.Size() == 0 {
		initDefaultList()
	}

	return Lparse(DefaultList, name)
}

func Ldomain(l *List, name string) (string, error) {
	dn, err := Lparse(l, name)
	if err != nil {
		return "", err
	}

	return dn.Sld + "." + dn.Tld, nil
}

func Lparse(l *List, name string) (*DomainName, error) {
	n, err := normalize(name)
	if err != nil {
		return nil, nil
	}

	r := l.Find(n)
	if tld := r.Decompose(n)[1]; tld == "" {
		return nil, fmt.Errorf("%s is a suffix", n)
	}

	dn := &DomainName{}
	dn.Tld, dn.Sld, dn.Trd = decompose(&r, n)
	return dn, nil
}

func initDefaultList() {
	err := DefaultList.LoadFile(defaultListFile, DefaultParserOptions)
	if err != nil {
		panic(err)
	}
}

func normalize(name string) (string, error) {
	return name, nil
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

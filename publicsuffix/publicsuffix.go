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
)

// DefaultList is the default List and is used by Parse.
var DefaultList = NewList()

// DefaultParserOptions are the default options used to parse a Public Suffix list.
var DefaultParserOptions = &ParserOption{PrivateDomains: true}

type Rule struct {
	Type    int
	Value   string
	Length  int
	Private bool
}

type ParserOption struct {
	PrivateDomains bool
}

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
	r := strings.NewReader(src)

	l := NewList()
	err := l.parse(r, options)
	return l, err
}

// NewListFromString parses a string that represents a Public Suffix source
// and returns a List initialized with the rules in the source.
func NewListFromFile(path string, options *ParserOption) (*List, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	l := NewList()
	err = l.parse(f, options)
	return l, err
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

// AddRule adds a new rule to the list.
//
// The exact position of the rule into the list is unpredictable.
// The list may be optimized internally for lookups, therefore the algorithm
// will decide the best position for the new rule.
func (l *List) AddRule(r *Rule) error {
	l.rules = append(l.rules, *r)
	return nil
}

func (l *List) rulesCount() int {
	return len(l.rules)
}

// NewRule parses the rule content, creates and returns a Rule.
func NewRule(content string) *Rule {
	var rule *Rule
	var value string

	switch content[0:1] {
	case "*": // wildcard
		value = content[2:len(content)]
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
//
// For exa
//
func (r *Rule) Decompose(name string) [2]string {
	var re *regexp.Regexp
	suffix := strings.Join(r.parts(), "./")

	switch r.Type {
	case WildcardType:
		re = regexp.MustCompile(fmt.Sprintf(`^(.+)\.(.*?\.%s)$`, suffix))
	default:
		re = regexp.MustCompile(fmt.Sprintf(`^(.+)\.(%s)$`, suffix))
	}

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
	return labels
}

func Labels(name string) []string {
	return strings.Split(name, ".")
}

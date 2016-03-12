package publicsuffix

import (
	"strings"
	"regexp"
	"fmt"
)

const (
	NormalType    = 1
	WildcardType  = 2
	ExceptionType = 3
)

type Rule struct {
	Type    int
	Value   string
	Length  int
	Private bool
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
	default:  // normal
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

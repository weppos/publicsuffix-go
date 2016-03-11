package publicsuffix

import (
	"strings"
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

func NewRule(content string) *Rule {
	var rule *Rule
	var value string

	switch content[0:1] {
	case "*":
		value = content[2:len(content)]
		rule = &Rule{Type: WildcardType, Value: value, Length: len(Labels(value)) + 1}
	case "!":
		value = content[1:len(content)]
		rule = &Rule{Type: ExceptionType, Value: value, Length: len(Labels(value))}
	default:
		value = content
		rule = &Rule{Type: NormalType, Value: value, Length: len(Labels(value))}
	}
	return rule
}

func Labels(name string) []string {
	return strings.Split(name, ".")
}

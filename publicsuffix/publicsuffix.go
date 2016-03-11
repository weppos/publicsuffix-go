package publicsuffix

import ()

const (
	NormalType    = 1
	WildcardType  = 2
	ExceptionType = 3
)

type Rule struct {
	Type    int
	Value   string
	Private bool
}

func NewRule(content string) *Rule {
	var rtype int
	var value string

	switch content[0:1] {
	case "*":
		rtype = WildcardType
		value = content[2:len(content)]
	case "!":
		rtype = ExceptionType
		value = content[1:len(content)]
	default:
		rtype = NormalType
		value = content
	}
	return &Rule{Type: rtype, Value: value}
}

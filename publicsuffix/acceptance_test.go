package publicsuffix

import (
	"reflect"
	"testing"
)

type validTestCase struct {
	input  string
	domain string
	parsed *DomainName
}

func TestValid(t *testing.T) {
	testCases := []validTestCase{
		validTestCase{"example.com", "example.com", &DomainName{"com", "example", ""}},
		validTestCase{"foo.example.com", "example.com", &DomainName{"com", "example", "foo"}},
	}

	for _, testCase := range testCases {
		got, err := Parse(testCase.input)
		if err != nil {
			t.Errorf("Parse(%v) returned error: %v", testCase.input, err)
		}
		if want := testCase.parsed; !reflect.DeepEqual(got, want) {
			t.Errorf("Parse(%v) = %v, want %v", testCase.input, got, want)
		}

		str, err := Domain(testCase.input)
		if err != nil {
			t.Errorf("Domain(%v) returned error: %v", testCase.input, err)
		}
		if want := testCase.domain; testCase.domain != str {
			t.Errorf("Domain(%v) = %v, want %v", testCase.input, str, want)
		}
	}
}

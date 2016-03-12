package publicsuffix

import (
	"reflect"
	"testing"
)

func TestNewRule_Normal(t *testing.T) {
	rule := NewRule("com")
	want := &Rule{Type: NormalType, Value: "com", Length: 1}

	if !reflect.DeepEqual(want, rule) {
		t.Errorf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Wildcard(t *testing.T) {
	rule := NewRule("*.example.com")
	want := &Rule{Type: WildcardType, Value: "example.com", Length: 3}

	if !reflect.DeepEqual(want, rule) {
		t.Errorf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Exception(t *testing.T) {
	rule := NewRule("!example.com")
	want := &Rule{Type: ExceptionType, Value: "example.com", Length: 2}

	if !reflect.DeepEqual(want, rule) {
		t.Errorf("NewRule returned %v, want %v", rule, want)
	}
}

type matchTestCase struct {
	rule     *Rule
	input    string
	expected bool
}

func TestMatch(t *testing.T) {
	testCases := []matchTestCase{
		// standard match
		matchTestCase{NewRule("uk"), "uk", true},
		matchTestCase{NewRule("uk"), "example.uk", true},
		matchTestCase{NewRule("uk"), "example.co.uk", true},
		matchTestCase{NewRule("co.uk"), "example.co.uk", true},

		// special rules match
		matchTestCase{NewRule("*.com"), "com", false},
		matchTestCase{NewRule("*.com"), "example.com", true},
		matchTestCase{NewRule("*.com"), "foo.example.com", true},
		matchTestCase{NewRule("!example.com"), "com", false},
		matchTestCase{NewRule("!example.com"), "example.com", true},
		matchTestCase{NewRule("!example.com"), "foo.example.com", true},

		// TLD mismatch
		matchTestCase{NewRule("gk"), "example.uk", false},
		matchTestCase{NewRule("gk"), "example.co.uk", false},

		// general mismatch
		matchTestCase{NewRule("uk.co"), "example.co.uk", false},
		matchTestCase{NewRule("go.uk"), "example.co.uk", false},
		// rule is longer than input, should not match
		matchTestCase{NewRule("co.uk"), "uk", false},

		// partial matches/mismatches
		matchTestCase{NewRule("co"), "example.co.uk", false},
		matchTestCase{NewRule("example"), "example.uk", false},
		matchTestCase{NewRule("le.it"), "example.it", false},
		matchTestCase{NewRule("le.it"), "le.it", true},
		matchTestCase{NewRule("le.it"), "foo.le.it", true},
	}

	for _, testCase := range testCases {
		if testCase.rule.Match(testCase.input) != testCase.expected {
			t.Errorf("Expected %v to %v match %v", testCase.rule.Value, testCase.expected, testCase.input)
		}
	}
}

type decomposeTestCase struct {
	rule     *Rule
	input    string
	expected [2]string
}

func TestDecompose(t *testing.T) {
	testCases := []decomposeTestCase{
		decomposeTestCase{NewRule("com"), "com", [2]string{"", ""}},
		decomposeTestCase{NewRule("com"), "example.com", [2]string{"example", "com"}},
		decomposeTestCase{NewRule("com"), "foo.example.com", [2]string{"foo.example", "com"}},

		decomposeTestCase{NewRule("!british-library.uk"), "uk", [2]string{"", ""}},
		decomposeTestCase{NewRule("!british-library.uk"), "british-library.uk", [2]string{"british-library", "uk"}},
		decomposeTestCase{NewRule("!british-library.uk"), "foo.british-library.uk", [2]string{"foo.british-library", "uk"}},

		decomposeTestCase{NewRule("*.com"), "com", [2]string{"", ""}},
		decomposeTestCase{NewRule("*.com"), "example.com", [2]string{"", ""}},
		decomposeTestCase{NewRule("*.com"), "foo.example.com", [2]string{"foo", "example.com"}},
		decomposeTestCase{NewRule("*.com"), "bar.foo.example.com", [2]string{"bar.foo", "example.com"}},
	}

	for _, testCase := range testCases {
		if got := testCase.rule.Decompose(testCase.input); !reflect.DeepEqual(got, testCase.expected) {
			t.Errorf("Expected %v to decompose %v into %v, got %v", testCase.rule.Value, testCase.input, testCase.expected, got)
		}
	}
}

func TestLabels(t *testing.T) {
	testCases := map[string][]string{
		"com":             []string{"com"},
		"example.com":     []string{"example", "com"},
		"www.example.com": []string{"www", "example", "com"},
	}

	for input, expected := range testCases {
		if output := Labels(input); !reflect.DeepEqual(output, expected) {
			t.Errorf("Labels(%v) = %v, want %v", input, output, expected)
		}
	}
}

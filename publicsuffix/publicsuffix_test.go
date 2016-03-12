package publicsuffix

import (
	"reflect"
	"testing"
)

func TestNewListFromString(t *testing.T) {
	src := `
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// ===BEGIN ICANN DOMAINS===

// ac : http://en.wikipedia.org/wiki/.ac
ac
com.ac

// ===END ICANN DOMAINS===
// ===BEGIN PRIVATE DOMAINS===

// Google, Inc.
blogspot.com

// ===END PRIVATE DOMAINS===
	`

	list, err := NewListFromString(src, nil)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	if want, got := 3, list.rulesCount(); want != got {
		t.Errorf("Parse returned a list with %v rules, want %v", got, want)
		t.Fatalf("%v", list.rules)
	}

	rules := list.rules
	var testRules []Rule

	testRules = []Rule{}
	for _, rule := range rules {
		if rule.Private == false {
			testRules = append(testRules, rule)
		}
	}
	if want, got := 2, len(testRules); want != got {
		t.Errorf("Parse returned a list with %v IANA rules, want %v", got, want)
		t.Fatalf("%v", testRules)
	}

	testRules = []Rule{}
	for _, rule := range rules {
		if rule.Private == true {
			testRules = append(testRules, rule)
		}
	}
	if want, got := 1, len(testRules); want != got {
		t.Errorf("Parse returned a list with %v PRIVATE rules, want %v", got, want)
		t.Fatalf("%v", testRules)
	}
}

func TestNewListFromFile(t *testing.T) {
	list, err := NewListFromFile("../fixtures/test.txt", nil)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	if want, got := 3, list.rulesCount(); want != got {
		t.Errorf("Parse returned a list with %v rules, want %v", got, want)
		t.Fatalf("%v", list.rules)
	}

	rules := list.rules
	var testRules []Rule

	testRules = []Rule{}
	for _, rule := range rules {
		if rule.Private == false {
			testRules = append(testRules, rule)
		}
	}
	if want, got := 2, len(testRules); want != got {
		t.Errorf("Parse returned a list with %v IANA rules, want %v", got, want)
		t.Fatalf("%v", testRules)
	}

	testRules = []Rule{}
	for _, rule := range rules {
		if rule.Private == true {
			testRules = append(testRules, rule)
		}
	}
	if want, got := 1, len(testRules); want != got {
		t.Errorf("Parse returned a list with %v PRIVATE rules, want %v", got, want)
		t.Fatalf("%v", testRules)
	}
}

func TestListAddRule(t *testing.T) {
	list := &List{}

	if list.rulesCount() != 0 {
		t.Fatalf("Empty list should have 0 rules, got %v", list.rulesCount())
	}

	rule := NewRule("com")
	list.AddRule(rule)
	if list.rulesCount() != 1 {
		t.Fatalf("List should have 1 rule, got %v", list.rulesCount())
	}
	if got := &list.rules[0]; !reflect.DeepEqual(rule, got) {
		t.Fatalf("List[0] expected to be %v, got %v", rule, got)
	}
}

func TestNewRule_Normal(t *testing.T) {
	rule := NewRule("com")
	want := &Rule{Type: NormalType, Value: "com", Length: 1}

	if !reflect.DeepEqual(want, rule) {
		t.Fatalf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Wildcard(t *testing.T) {
	rule := NewRule("*.example.com")
	want := &Rule{Type: WildcardType, Value: "example.com", Length: 3}

	if !reflect.DeepEqual(want, rule) {
		t.Fatalf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Exception(t *testing.T) {
	rule := NewRule("!example.com")
	want := &Rule{Type: ExceptionType, Value: "example.com", Length: 2}

	if !reflect.DeepEqual(want, rule) {
		t.Fatalf("NewRule returned %v, want %v", rule, want)
	}
}

type matchTestCase struct {
	rule     *Rule
	input    string
	expected bool
}

func TestRuleMatch(t *testing.T) {
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

func TestRuleDecompose(t *testing.T) {
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

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

	if want, got := 3, list.Size(); want != got {
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

func TestNewListFromString_IDNAInputIsUnicode(t *testing.T) {
	src := `
// xn--d1alf ("mkd", Macedonian) : MK
// MARnet
мкд

// xn--l1acc ("mon", Mongolian) : MN
xn--l1acc
	`

	list, err := NewListFromString(src, nil)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	if want, got := 2, list.Size(); want != got {
		t.Errorf("Parse returned a list with %v rules, want %v", got, want)
		t.Fatalf("%v", list.rules)
	}

	if rule := list.Find("hello.xn--d1alf", &FindOptions{DefaultRule: nil}); rule == nil {
		t.Fatalf("Find(%v) returned nil", "hello.xn--d1alf")
	}
	if rule := list.Find("hello.мкд", &FindOptions{DefaultRule: nil}); rule != nil {
		t.Fatalf("Find(%v) expected to return nil, got %v", "hello.xn--d1alf", rule)
	}
	if rule := list.Find("hello.xn--l1acc", &FindOptions{DefaultRule: nil}); rule == nil {
		t.Fatalf("Find(%v) returned nil", "hello.xn--l1acc")
	}
}

func TestNewListFromString_IDNAInputIsAscii(t *testing.T) {
	src := `
// xn--d1alf ("mkd", Macedonian) : MK
// MARnet
xn--d1alf

// xn--l1acc ("mon", Mongolian) : MN
xn--l1acc
	`

	list, err := NewListFromString(src, &ParserOption{ASCIIEncoded: true})
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	if want, got := 2, list.Size(); want != got {
		t.Errorf("Parse returned a list with %v rules, want %v", got, want)
		t.Fatalf("%v", list.rules)
	}

	if rule := list.Find("hello.xn--d1alf", &FindOptions{DefaultRule: nil}); rule == nil {
		t.Fatalf("Find(%v) returned nil", "hello.xn--d1alf")
	}
	if rule := list.Find("hello.мкд", &FindOptions{DefaultRule: nil}); rule != nil {
		t.Fatalf("Find(%v) expected to return nil, got %v", "hello.xn--d1alf", rule)
	}
	if rule := list.Find("hello.xn--l1acc", &FindOptions{DefaultRule: nil}); rule == nil {
		t.Fatalf("Find(%v) returned nil", "hello.xn--l1acc")
	}
}

func TestNewListFromFile(t *testing.T) {
	list, err := NewListFromFile("../fixtures/list-simple.txt", nil)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	if want, got := 3, list.Size(); want != got {
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
	list := NewList()

	if list.Size() != 0 {
		t.Fatalf("Empty list should have 0 rules, got %v", list.Size())
	}

	rule := MustNewRule("com")
	list.AddRule(rule)
	if list.Size() != 1 {
		t.Fatalf("List should have 1 rule, got %v", list.Size())
	}
	if got := &list.rules[0]; !reflect.DeepEqual(rule, got) {
		t.Fatalf("List[0] expected to be %v, got %v", rule, got)
	}
}

type listFindTestCase struct {
	input    string
	expected *Rule
}

func TestListFind(t *testing.T) {
	src := `
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// ===BEGIN ICANN DOMAINS===

// com
com

// uk
*.uk
*.sch.uk
!bl.uk
!british-library.uk

// io
io

// ===END ICANN DOMAINS===
// ===BEGIN PRIVATE DOMAINS===

// Google, Inc.
blogspot.com

// ===END PRIVATE DOMAINS===
	`

	// TODO(weppos): ability to set type to a rule.
	p1 := MustNewRule("blogspot.com")
	p1.Private = true

	testCases := []listFindTestCase{
		// match IANA
		{"example.com", MustNewRule("com")},
		{"foo.example.com", MustNewRule("com")},

		// match wildcard
		{"example.uk", MustNewRule("*.uk")},
		{"example.co.uk", MustNewRule("*.uk")},
		{"foo.example.co.uk", MustNewRule("*.uk")},

		// match exception
		{"british-library.uk", MustNewRule("!british-library.uk")},
		{"foo.british-library.uk", MustNewRule("!british-library.uk")},

		// match default rule
		{"test", DefaultRule},
		{"example.test", DefaultRule},
		{"foo.example.test", DefaultRule},

		// match private
		{"blogspot.com", p1},
		{"foo.blogspot.com", p1},
	}

	list, err := NewListFromString(src, nil)
	if err != nil {
		t.Fatalf("Unable to parse list: %v", err)
	}

	for _, testCase := range testCases {
		if want, got := testCase.expected, list.Find(testCase.input, nil); !reflect.DeepEqual(want, got) {
			t.Errorf("Find(%v) = %v, want %v", testCase.input, got, want)
		}
	}
}

func TestNewRule_Normal(t *testing.T) {
	rule := MustNewRule("com")
	want := &Rule{Type: NormalType, Value: "com", Length: 1}

	if !reflect.DeepEqual(want, rule) {
		t.Fatalf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Wildcard(t *testing.T) {
	rule := MustNewRule("*.example.com")
	want := &Rule{Type: WildcardType, Value: "example.com", Length: 3}

	if !reflect.DeepEqual(want, rule) {
		t.Fatalf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Exception(t *testing.T) {
	rule := MustNewRule("!example.com")
	want := &Rule{Type: ExceptionType, Value: "example.com", Length: 2}

	if !reflect.DeepEqual(want, rule) {
		t.Fatalf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_FromASCII(t *testing.T) {
	rule, _ := NewRule("xn--l1acc")

	if want := "xn--l1acc"; rule.Value != want {
		t.Fatalf("NewRule == %v, want %v", rule.Value, want)
	}
}
func TestNewRule_FromUnicode(t *testing.T) {
	rule, _ := NewRule("мон")

	// No transformation is performed
	if want := "мон"; rule.Value != want {
		t.Fatalf("NewRule == %v, want %v", rule.Value, want)
	}
}

func TestNewRuleUnicode_FromASCII(t *testing.T) {
	rule, _ := NewRuleUnicode("xn--l1acc")

	if want := "xn--l1acc"; rule.Value != want {
		t.Fatalf("NewRule == %v, want %v", rule.Value, want)
	}
}

func TestNewRuleUnicode_FromUnicode(t *testing.T) {
	rule, _ := NewRuleUnicode("мон")

	if want := "xn--l1acc"; rule.Value != want {
		t.Fatalf("NewRule == %v, want %v", rule.Value, want)
	}
}

type ruleMatchTestCase struct {
	rule     *Rule
	input    string
	expected bool
}

func TestRuleMatch(t *testing.T) {
	testCases := []ruleMatchTestCase{
		// standard match
		{MustNewRule("uk"), "uk", true},
		{MustNewRule("uk"), "example.uk", true},
		{MustNewRule("uk"), "example.co.uk", true},
		{MustNewRule("co.uk"), "example.co.uk", true},

		// special rules match
		{MustNewRule("*.com"), "com", false},
		{MustNewRule("*.com"), "example.com", true},
		{MustNewRule("*.com"), "foo.example.com", true},
		{MustNewRule("!example.com"), "com", false},
		{MustNewRule("!example.com"), "example.com", true},
		{MustNewRule("!example.com"), "foo.example.com", true},

		// TLD mismatch
		{MustNewRule("gk"), "example.uk", false},
		{MustNewRule("gk"), "example.co.uk", false},

		// general mismatch
		{MustNewRule("uk.co"), "example.co.uk", false},
		{MustNewRule("go.uk"), "example.co.uk", false},
		// rule is longer than input, should not match
		{MustNewRule("co.uk"), "uk", false},

		// partial matches/mismatches
		{MustNewRule("co"), "example.co.uk", false},
		{MustNewRule("example"), "example.uk", false},
		{MustNewRule("le.it"), "example.it", false},
		{MustNewRule("le.it"), "le.it", true},
		{MustNewRule("le.it"), "foo.le.it", true},
	}

	for _, testCase := range testCases {
		if testCase.rule.Match(testCase.input) != testCase.expected {
			t.Errorf("Expected %v to %v match %v", testCase.rule.Value, testCase.expected, testCase.input)
		}
	}
}

type ruleDecomposeTestCase struct {
	rule     *Rule
	input    string
	expected [2]string
}

func TestRuleDecompose(t *testing.T) {
	testCases := []ruleDecomposeTestCase{
		{MustNewRule("com"), "com", [2]string{"", ""}},
		{MustNewRule("com"), "example.com", [2]string{"example", "com"}},
		{MustNewRule("com"), "foo.example.com", [2]string{"foo.example", "com"}},

		{MustNewRule("!british-library.uk"), "uk", [2]string{"", ""}},
		{MustNewRule("!british-library.uk"), "british-library.uk", [2]string{"british-library", "uk"}},
		{MustNewRule("!british-library.uk"), "foo.british-library.uk", [2]string{"foo.british-library", "uk"}},

		{MustNewRule("*.com"), "com", [2]string{"", ""}},
		{MustNewRule("*.com"), "example.com", [2]string{"", ""}},
		{MustNewRule("*.com"), "foo.example.com", [2]string{"foo", "example.com"}},
		{MustNewRule("*.com"), "bar.foo.example.com", [2]string{"bar.foo", "example.com"}},
	}

	for _, testCase := range testCases {
		if got := testCase.rule.Decompose(testCase.input); !reflect.DeepEqual(got, testCase.expected) {
			t.Errorf("Expected %v to decompose %v into %v, got %v", testCase.rule.Value, testCase.input, testCase.expected, got)
		}
	}
}

func TestLabels(t *testing.T) {
	testCases := map[string][]string{
		"com":             {"com"},
		"example.com":     {"example", "com"},
		"www.example.com": {"www", "example", "com"},
	}

	for input, expected := range testCases {
		if output := Labels(input); !reflect.DeepEqual(output, expected) {
			t.Errorf("Labels(%v) = %v, want %v", input, output, expected)
		}
	}
}

func TestCookieJarList(t *testing.T) {
	testCases := map[string]string{
		"example.com":              "com",
		"www.example.com":          "com",
		"example.co.uk":            "co.uk",
		"www.example.co.uk":        "co.uk",
		"example.blogspot.com":     "blogspot.com",
		"www.example.blogspot.com": "blogspot.com",
		"parliament.uk":            "uk",
		"www.parliament.uk":        "uk",
		// not listed
		"www.example.test": "test",
	}

	for input, suffix := range testCases {
		if output := CookieJarList.PublicSuffix(input); output != suffix {
			t.Errorf("CookieJarList.PublicSuffix(%v) = %v, want %v", input, output, suffix)
		}
	}
}

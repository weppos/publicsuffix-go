package publicsuffix

import (
	"testing"
)

type validTestCase struct {
	input  string
	domain string
	parsed *DomainName
}

func TestValid(t *testing.T) {
	testCases := []validTestCase{
		{"example.com", "example.com", &DomainName{"com", "example", "", MustNewRule("com")}},
		{"foo.example.com", "example.com", &DomainName{"com", "example", "foo", MustNewRule("com")}},

		{"verybritish.co.uk", "verybritish.co.uk", &DomainName{"co.uk", "verybritish", "", MustNewRule("*.uk")}},
		{"foo.verybritish.co.uk", "verybritish.co.uk", &DomainName{"co.uk", "verybritish", "foo", MustNewRule("*.uk")}},

		{"parliament.uk", "parliament.uk", &DomainName{"uk", "parliament", "", MustNewRule("!parliament.uk")}},
		{"foo.parliament.uk", "parliament.uk", &DomainName{"uk", "parliament", "foo", MustNewRule("!parliament.uk")}},

		{"foo.blogspot.com", "foo.blogspot.com", &DomainName{"blogspot.com", "foo", "", MustNewRule("blogspot.com")}},
		{"bar.foo.blogspot.com", "foo.blogspot.com", &DomainName{"blogspot.com", "foo", "bar", MustNewRule("blogspot.com")}},
	}

	for _, testCase := range testCases {
		got, err := Parse(testCase.input)
		if err != nil {
			t.Errorf("TestValid(%v) returned error: %v", testCase.input, err)
		}
		if want := testCase.parsed; want.String() != got.String() {
			t.Errorf("TestValid(%v) = %v, want %v", testCase.input, got, want)
		}

		str, err := Domain(testCase.input)
		if err != nil {
			t.Errorf("TestValid(%v) returned error: %v", testCase.input, err)
		}
		if want := testCase.domain; want != str {
			t.Errorf("TestValid(%v) = %v, want %v", testCase.input, str, want)
		}
	}
}

type privateTestCase struct {
	input  string
	domain string
	ignore bool
	error  bool
}

func TestIncludePrivate(t *testing.T) {
	testCases := []privateTestCase{
		{"blogspot.com", "", false, true},
		{"blogspot.com", "blogspot.com", true, false},

		{"foo.blogspot.com", "foo.blogspot.com", false, false},
		{"foo.blogspot.com", "blogspot.com", true, false},
	}

	for _, testCase := range testCases {
		got, err := DomainFromListWithOptions(DefaultList, testCase.input, &FindOptions{IgnorePrivate: testCase.ignore})

		if testCase.error && err == nil {
			t.Errorf("TestIncludePrivate(%v) should have returned error, got: %v", testCase.input, got)
			continue
		}
		if !testCase.error && err != nil {
			t.Errorf("TestIncludePrivate(%v) returned error: %v", testCase.input, err)
			continue
		}

		if want := testCase.domain; want != got {
			t.Errorf("Domain(%v) = %v, want %v", testCase.input, got, want)
		}
	}
}

type idnaTestCase struct {
	input  string
	domain string
	error  bool
}

func TestIDNA(t *testing.T) {
	testACases := []idnaTestCase{
		// A-labels are supported
		// Check single IDN part
		{"xn--p1ai", "", true},
		{"example.xn--p1ai", "example.xn--p1ai", false},
		{"subdomain.example.xn--p1ai", "example.xn--p1ai", false},
		// Check multiple IDN parts
		{"xn--example--3bhk5a.xn--p1ai", "xn--example--3bhk5a.xn--p1ai", false},
		{"subdomain.xn--example--3bhk5a.xn--p1ai", "xn--example--3bhk5a.xn--p1ai", false},
		// Check multiple IDN rules
		{"example.xn--o1ach.xn--90a3ac", "example.xn--o1ach.xn--90a3ac", false},
		{"sudbomain.example.xn--o1ach.xn--90a3ac", "example.xn--o1ach.xn--90a3ac", false},
	}

	for _, testCase := range testACases {
		got, err := DomainFromListWithOptions(DefaultList, testCase.input, nil)

		if testCase.error && err == nil {
			t.Errorf("A-label %v should have returned error, got: %v", testCase.input, got)
			continue
		}
		if !testCase.error && err != nil {
			t.Errorf("A-label %v returned error: %v", testCase.input, err)
			continue
		}

		if want := testCase.domain; want != got {
			t.Errorf("A-label Domain(%v) = %v, want %v", testCase.input, got, want)
		}
	}

	// These tests validates the non-acceptance of U-labels.
	//
	// TODO(weppos): some tests are passing because of the default rule *
	// Consider to add some tests overriding the default rule to nil.
	// Right now, setting the default rule to nil with cause a panic if the lookup results in a nil.
	testUCases := []idnaTestCase{
		// U-labels are NOT supported
		// Check single IDN part
		{"рф", "", true},
		{"example.рф", "example.рф", false},           // passes because of *
		{"subdomain.example.рф", "example.рф", false}, // passes because of *
		// Check multiple IDN parts
		{"example-упр.рф", "example-упр.рф", false},           // passes because of *
		{"subdomain.example-упр.рф", "example-упр.рф", false}, // passes because of *
		// Check multiple IDN rules
		{"example.упр.срб", "упр.срб", false},
		{"sudbomain.example.упр.срб", "упр.срб", false},
	}

	for _, testCase := range testUCases {
		got, err := DomainFromListWithOptions(DefaultList, testCase.input, nil)

		if testCase.error && err == nil {
			t.Errorf("U-label %v should have returned error, got: %v", testCase.input, got)
			continue
		}
		if !testCase.error && err != nil {
			t.Errorf("U-label %v returned error: %v", testCase.input, err)
			continue
		}

		if want := testCase.domain; want != got {
			t.Errorf("U-label Domain(%v) = %v, want %v", testCase.input, got, want)
		}
	}
}

func TestFindRuleIANA(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		// TLD with only 1 rule.
		{"biz", "biz"},
		{"input.biz", "biz"},
		{"b.input.biz", "biz"},

		// The relevant {kobe,kyoto}.jp rules are:
		// jp
		// *.kobe.jp
		// !city.kobe.jp
		// kyoto.jp
		// ide.kyoto.jp
		{"jp", "jp"},
		{"kobe.jp", "jp"},
		{"c.kobe.jp", "c.kobe.jp"},
		{"b.c.kobe.jp", "c.kobe.jp"},
		{"a.b.c.kobe.jp", "c.kobe.jp"},
		{"city.kobe.jp", "kobe.jp"},
		{"www.city.kobe.jp", "kobe.jp"},
		{"kyoto.jp", "kyoto.jp"},
		{"test.kyoto.jp", "kyoto.jp"},
		{"ide.kyoto.jp", "ide.kyoto.jp"},
		{"b.ide.kyoto.jp", "ide.kyoto.jp"},
		{"a.b.ide.kyoto.jp", "ide.kyoto.jp"},

		// Domain with a private public suffix should return the ICANN public suffix.
		{"foo.compute-1.amazonaws.com", "com"},
		// Domain equal to a private public suffix should return the ICANN public suffix.
		{"cloudapp.net", "net"},
	}

	for _, tc := range testCases {
		rule := DefaultList.Find(tc.input, &FindOptions{IgnorePrivate: true, DefaultRule: nil})

		if rule == nil {
			t.Errorf("TestFindRuleIANA(%v) nil rule", tc.input)
			continue
		}

		suffix := rule.Decompose(tc.input)[1]
		// If the TLD is empty, it means name is actually a suffix.
		// In fact, decompose returns an array of empty strings in this case.
		if suffix == "" {
			suffix = tc.input
		}

		if suffix != tc.want {
			t.Errorf("TestFindRuleIANA(%v) = %v, want %v", tc.input, suffix, tc.want)
		}
	}
}

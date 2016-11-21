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
		validTestCase{"example.com", "example.com", &DomainName{"com", "example", "", MustNewRule("com")}},
		validTestCase{"foo.example.com", "example.com", &DomainName{"com", "example", "foo", MustNewRule("com")}},

		validTestCase{"verybritish.co.uk", "verybritish.co.uk", &DomainName{"co.uk", "verybritish", "", MustNewRule("*.uk")}},
		validTestCase{"foo.verybritish.co.uk", "verybritish.co.uk", &DomainName{"co.uk", "verybritish", "foo", MustNewRule("*.uk")}},

		validTestCase{"parliament.uk", "parliament.uk", &DomainName{"uk", "parliament", "", MustNewRule("!parliament.uk")}},
		validTestCase{"foo.parliament.uk", "parliament.uk", &DomainName{"uk", "parliament", "foo", MustNewRule("!parliament.uk")}},

		validTestCase{"foo.blogspot.com", "foo.blogspot.com", &DomainName{"blogspot.com", "foo", "", MustNewRule("blogspot.com")}},
		validTestCase{"bar.foo.blogspot.com", "foo.blogspot.com", &DomainName{"blogspot.com", "foo", "bar", MustNewRule("blogspot.com")}},
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
		privateTestCase{"blogspot.com", "", false, true},
		privateTestCase{"blogspot.com", "blogspot.com", true, false},

		privateTestCase{"foo.blogspot.com", "foo.blogspot.com", false, false},
		privateTestCase{"foo.blogspot.com", "blogspot.com", true, false},
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
		idnaTestCase{"xn--p1ai", "", true},
		idnaTestCase{"example.xn--p1ai", "example.xn--p1ai", false},
		idnaTestCase{"subdomain.example.xn--p1ai", "example.xn--p1ai", false},
		// Check multiple IDN parts
		idnaTestCase{"xn--example--3bhk5a.xn--p1ai", "xn--example--3bhk5a.xn--p1ai", false},
		idnaTestCase{"subdomain.xn--example--3bhk5a.xn--p1ai", "xn--example--3bhk5a.xn--p1ai", false},
		// Check multiple IDN rules
		idnaTestCase{"example.xn--o1ach.xn--90a3ac", "example.xn--o1ach.xn--90a3ac", false},
		idnaTestCase{"sudbomain.example.xn--o1ach.xn--90a3ac", "example.xn--o1ach.xn--90a3ac", false},
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
		idnaTestCase{"рф", "", true},
		idnaTestCase{"example.рф", "example.рф", false},           // passes because of *
		idnaTestCase{"subdomain.example.рф", "example.рф", false}, // passes because of *
		// Check multiple IDN parts
		idnaTestCase{"example-упр.рф", "example-упр.рф", false},           // passes because of *
		idnaTestCase{"subdomain.example-упр.рф", "example-упр.рф", false}, // passes because of *
		// Check multiple IDN rules
		idnaTestCase{"example.упр.срб", "упр.срб", false},
		idnaTestCase{"sudbomain.example.упр.срб", "упр.срб", false},
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

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

		validTestCase{"verybritish.co.uk", "verybritish.co.uk", &DomainName{"co.uk", "verybritish", ""}},
		validTestCase{"foo.verybritish.co.uk", "verybritish.co.uk", &DomainName{"co.uk", "verybritish", "foo"}},

		validTestCase{"parliament.uk", "parliament.uk", &DomainName{"uk", "parliament", ""}},
		validTestCase{"foo.parliament.uk", "parliament.uk", &DomainName{"uk", "parliament", "foo"}},

		validTestCase{"foo.blogspot.com", "foo.blogspot.com", &DomainName{"blogspot.com", "foo", ""}},
		validTestCase{"bar.foo.blogspot.com", "foo.blogspot.com", &DomainName{"blogspot.com", "foo", "bar"}},
	}

	for _, testCase := range testCases {
		got, err := Parse(testCase.input)
		if err != nil {
			t.Errorf("TestValid(%v) returned error: %v", testCase.input, err)
		}
		if want := testCase.parsed; !reflect.DeepEqual(got, want) {
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
	ignore bool
	error  bool
	domain string
}

func TestIncludePrivate(t *testing.T) {
	testCases := []privateTestCase{
		privateTestCase{"blogspot.com", false, true, ""},
		privateTestCase{"blogspot.com", true, false, "blogspot.com"},

		privateTestCase{"foo.blogspot.com", false, false, "foo.blogspot.com"},
		privateTestCase{"foo.blogspot.com", true, false, "blogspot.com"},
	}

	for _, testCase := range testCases {
		got, err := DomainFromListWithOptions(DefaultList(), testCase.input, &FindOptions{IgnorePrivate: testCase.ignore})

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

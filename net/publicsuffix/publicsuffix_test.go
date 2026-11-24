package publicsuffix_test

import (
	"testing"

	wpsl "github.com/weppos/publicsuffix-go/net/publicsuffix"
	xpsl "golang.org/x/net/publicsuffix"
)

func TestPublicSuffix(t *testing.T) {
	testCases := []string{
		"example.com",
		"www.example.com",
		"example.co.uk",
		"www.example.co.uk",
		"example.blogspot.com",
		"www.example.blogspot.com",
		"parliament.uk",
		"www.parliament.uk",
		// not listed
		"www.example.test",
	}

	for _, testCase := range testCases {
		ws, wb := wpsl.PublicSuffix(testCase)
		xs, xb := xpsl.PublicSuffix(testCase)

		if ws != xs || wb != xb {
			t.Errorf("PublicSuffix(%v): x/psl -> (%v, %v) != w/psl -> (%v, %v)", testCase, xs, xb, ws, wb)
		}
	}
}

func TestPublicSuffixEdgeCases(t *testing.T) {
	testCases := []struct {
		domain string
		suffix string
		icann  bool
	}{
		// the domain itself is a public suffix
		{"co.uk", "co.uk", true},
		{"com", "com", true},
		{"foo.ck", "foo.ck", true},
		{"blogspot.com", "blogspot.com", false},
		{"test", "test", false},
		// trailing dot
		{"example.com.", "", false},
		{"test.", "", false},
		{".", "", false},
		// IP address
		{"127.0.0.1", "127.0.0.1", false},
		{"::1", "::1", false},
	}

	for _, tc := range testCases {
		suffix, icann := wpsl.PublicSuffix(tc.domain)
		if suffix != tc.suffix || icann != tc.icann {
			t.Errorf("PublicSuffix(%q) = (%q, %v), want (%q, %v)", tc.domain, suffix, icann, tc.suffix, tc.icann)
		}
	}
}

func TestEffectiveTLDPlusOne(t *testing.T) {
	testCases := []string{
		"example.com",
		"www.example.com",
		"example.co.uk",
		"www.example.co.uk",
		"example.blogspot.com",
		"www.example.blogspot.com",
		"parliament.uk",
		"www.parliament.uk",
		// not listed
		"www.example.test",
	}

	for _, testCase := range testCases {
		ws, we := wpsl.EffectiveTLDPlusOne(testCase)
		xs, xe := xpsl.EffectiveTLDPlusOne(testCase)

		if ws != xs || we != xe {
			t.Errorf("EffectiveTLDPlusOne(%v): x/psl -> (%v, %v) != w/psl -> (%v, %v)", testCase, xs, xe, ws, we)
		}
	}
}

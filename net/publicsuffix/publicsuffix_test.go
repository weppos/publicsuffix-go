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

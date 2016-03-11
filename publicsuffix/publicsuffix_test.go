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

func TestLabels(t *testing.T) {
	testCases := map[string][]string {
		"com": []string{"com"},
		"example.com": []string{"example", "com"},
		"www.example.com": []string{"www", "example", "com"},
	}

	for input, expected := range(testCases) {
		if output := Labels(input); !reflect.DeepEqual(output, expected) {
			t.Errorf("Labels(%v) = %v, want %v", input, output, expected)
		}
	}
}

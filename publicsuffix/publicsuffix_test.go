package publicsuffix

import (
	"reflect"
	"testing"
)

func TestNewRule_Normal(t *testing.T) {
	rule := NewRule("com")
	want := &Rule{Type: NormalType, Value: "com"}

	if !reflect.DeepEqual(want, rule) {
		t.Errorf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Wildcard(t *testing.T) {
	rule := NewRule("*.example.com")
	want := &Rule{Type: WildcardType, Value: "example.com"}

	if !reflect.DeepEqual(want, rule) {
		t.Errorf("NewRule returned %v, want %v", rule, want)
	}
}

func TestNewRule_Exception(t *testing.T) {
	rule := NewRule("!example.com")
	want := &Rule{Type: ExceptionType, Value: "example.com"}

	if !reflect.DeepEqual(want, rule) {
		t.Errorf("NewRule returned %v, want %v", rule, want)
	}
}

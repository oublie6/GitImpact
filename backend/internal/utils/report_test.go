package utils

import "testing"

func TestParseStructuredJSON(t *testing.T) {
	ok := `{"summary":"s","changed_modules":["a"]}`
	if _, err := ParseStructuredJSON(ok); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseStructuredJSON("not-json"); err == nil {
		t.Fatal("expected parse error")
	}
}

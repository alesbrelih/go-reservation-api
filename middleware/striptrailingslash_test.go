package middleware_test

import (
	"testing"

	"github.com/alesbrelih/go-reservation-api/middleware"
)

func TestStripImplementation(t *testing.T) {
	original := "/hi/hello/"
	stripped := middleware.StripImplementation(original)

	expected := "/hi/hello"
	if stripped != "/hi/hello" {
		t.Errorf("Stripped is incorrect. Got: %v, expected: %v ", stripped, expected)
	}
}

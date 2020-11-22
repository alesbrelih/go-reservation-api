package myutil_test

import (
	"testing"

	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
)

func TestGetEnvOrDefault_MissingEnv(t *testing.T) {
	fallbackValue := "hello test"
	myVal := myutil.GetEnvOrDefault("SOME_STUFF", fallbackValue)

	if myVal != fallbackValue {
		t.Errorf("Env val is incorrect, got: %s, expected %s", myVal, fallbackValue)
	}

}

package util_test

import (
	"testing"

	"github.com/alesbrelih/go-reservation-api/internal/util"
)

func TestGetEnvOrDefault_MissingEnv(t *testing.T) {
	fallbackValue := "hello test"
	myVal := util.GetEnvOrDefault("SOME_STUFF", fallbackValue)

	if myVal != fallbackValue {
		t.Errorf("Env val is incorrect, got: %s, expected %s", myVal, fallbackValue)
	}

}

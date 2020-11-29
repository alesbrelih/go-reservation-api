package myutil_test

import (
	"fmt"
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

func TestMissingEnvVariableMsg(t *testing.T) {
	myEnv := "ENV_ENV"
	msg := myutil.MissingEnvVariableMsg(myEnv)
	expected := fmt.Sprintf("Missing enviroment variable: %v", myEnv)
	if msg != expected {
		t.Errorf("Excepted: %v, got: %v", expected, msg)
	}
}

func TestEmpty_NonEmptyString(t *testing.T) {
	val := "sadasd"
	isEmpty := myutil.Empty(val)

	if isEmpty == true {
		t.Error("Is empty is true, expected false")
	}
}

func TestEmpty_EmptyString(t *testing.T) {
	val := ""
	isEmpty := myutil.Empty(val)

	if isEmpty == false {
		t.Error("Is empty is false, expected true")
	}
}
func TestEmpty_EmptySpaces(t *testing.T) {
	val := "  "
	isEmpty := myutil.Empty(val)

	if isEmpty == false {
		t.Error("Is empty is false, expected true")
	}
}

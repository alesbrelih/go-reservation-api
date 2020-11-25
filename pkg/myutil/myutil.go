package myutil

import (
	"fmt"
	"os"
)

// https://github.com/joeshaw/envdecode <- maybe try to use this
func GetEnvOrDefault(envName string, fallback string) string {
	envValue, found := os.LookupEnv(envName)
	if !found {
		return fallback
	}
	return envValue
}

func MissingEnvVariableMsg(envName string) string {
	return fmt.Sprintf("Missing enviroment variable: %v", envName)
}

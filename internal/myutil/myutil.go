package myutil

import (
	"os"
)

func GetEnvOrDefault(envName string, fallback string) string {
	envValue, found := os.LookupEnv(envName)
	if !found {
		return fallback
	}
	return envValue
}

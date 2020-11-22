package main

import (
	"net/http"

	"github.com/alesbrelih/go-reservation-api/internal/myutil"
)

func main() {

	// start server
	port := myutil.GetEnvOrDefault("APPLICATION_PORT", "8080")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err.Error())
	}
}

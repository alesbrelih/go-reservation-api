package main

import (
	"net/http"

	"github.com/alesbrelih/go-reservation-api/internal/util"
)

func main() {

	// start server
	port := util.GetEnvOrDefault("APPLICATION_PORT", "8080")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err.Error())
	}
}

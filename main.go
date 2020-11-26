package main

import (
	"net/http"

	"github.com/alesbrelih/go-reservation-api/middleware"
	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
	"github.com/alesbrelih/go-reservation-api/router"
)

func main() {

	port := myutil.GetEnvOrDefault("APPLICATION_PORT", "8080")

	mux := router.InitializeRouter()

	// start server
	err := http.ListenAndServe(":"+port, middleware.StripTrailingSlash(mux))
	if err != nil {
		panic(err.Error())
	}
}

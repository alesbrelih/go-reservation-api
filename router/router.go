package router

import "net/http"

func InitializeRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}

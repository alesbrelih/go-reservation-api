package middleware

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
)

func NewJsonMiddleware() JsonMiddleware {
	return &jsonMiddleware{}
}

type JsonMiddleware interface {
	SetResponseHeader(http.Handler) http.Handler
}

type jsonMiddleware struct {
	log hclog.Logger
}

func (j *jsonMiddleware) SetResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

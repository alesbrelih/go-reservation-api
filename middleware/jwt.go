package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alesbrelih/go-reservation-api/services"
	"github.com/hashicorp/go-hclog"
)

func NewJwt(auth services.AuthService, log hclog.Logger) Jwt {
	return &jwt{
		log:  log,
		auth: auth,
	}
}

type Jwt interface {
	ValidateUser(http.Handler) http.Handler
}

type jwt struct {
	log  hclog.Logger
	auth services.AuthService
}

type JwtClaimsContextKey struct{}

func (j *jwt) ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("authorization")
		if header == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		split := strings.Split(header, " ")
		if len(split) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		claims, err := j.auth.GetClaims(split[1])

		if err != nil {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		claimsCtx := context.WithValue(r.Context(), &JwtClaimsContextKey{}, claims)

		next.ServeHTTP(w, r.WithContext(claimsCtx))
	})
}

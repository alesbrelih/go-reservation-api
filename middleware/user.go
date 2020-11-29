package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/alesbrelih/go-reservation-api/models"
)

type UserBodyContextKey struct{}

type User struct {
	log log.Logger
}

func (u *User) GetBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := &models.UserReqBody{}
		err := user.FromJSON(r.Body)

		if err != nil {
			u.log.Printf("Cant read create User body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), &UserBodyContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewUserMiddleware() *User {
	return &User{}
}

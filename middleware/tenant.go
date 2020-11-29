package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/alesbrelih/go-reservation-api/models"
)

type TenantBodyContextKey struct{}

type Tenant struct {
	log *log.Logger
}

func (m *Tenant) GetBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tenant := &models.Tenant{}
		err := tenant.FromJSON(r.Body)
		if err != nil {
			m.log.Printf("Cant read tenant json. Error: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		cxt := context.WithValue(r.Context(), &TenantBodyContextKey{}, tenant)
		next.ServeHTTP(w, r.WithContext(cxt))
	})
}

func NewTenantMiddleware(log *log.Logger) *Tenant {
	return &Tenant{
		log: log,
	}
}

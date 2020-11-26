package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/alesbrelih/go-reservation-api/models"
)

type Item struct {
	log *log.Logger
}

type ItemBodyKeyType struct{}

var ItemBodyVar = "item-body-var"

func (im *Item) GetBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		item := &models.Item{}
		err := item.FromJSON(r.Body)
		if err != nil {
			im.log.Printf("Cant read Create item body: %v", err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), &ItemBodyKeyType{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewItemMiddleware(log *log.Logger) *Item {
	return &Item{
		log: log,
	}
}

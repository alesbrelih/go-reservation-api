package router

import (
	"github.com/alesbrelih/go-reservation-api/item"
	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/item").Handler(item.Router(&item.DefaultItemController{}))
	return r
}

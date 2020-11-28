package router

import (
	"log"

	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/item").Handler(controller.NewItemRouter(&controller.DefaultItemController{}))

	// user handler
	userHandler := controller.NewUserHandler(&stores.UserStoreSql{}, &log.Logger{})
	r.PathPrefix("/user").Handler(userHandler.NewRouter())

	return r
}

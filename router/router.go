package router

import (
	"log"
	"os"

	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/item").Handler(controller.NewItemRouter(&controller.DefaultItemController{}))

	// user handler
	userHandler := controller.NewUserHandler(&stores.UserStoreSql{}, log.New(os.Stdout, "user-controller ", log.LstdFlags))
	r.PathPrefix("/user").Handler(userHandler.NewRouter())

	return r
}

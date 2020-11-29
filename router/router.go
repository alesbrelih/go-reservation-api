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

	itemHandler := controller.NewItemHandler(&stores.ItemStoreSql{}, log.New(os.Stdout, "item-controller ", log.LstdFlags))
	r.PathPrefix("/item").Handler(itemHandler.NewItemRouter())

	// user handler
	userHandler := controller.NewUserHandler(&stores.UserStoreSql{}, log.New(os.Stdout, "user-controller ", log.LstdFlags))
	r.PathPrefix("/user").Handler(userHandler.NewRouter())

	// user handler
	tenantHandler := controller.NewTenantHandler(&stores.TenantStoreSql{}, log.New(os.Stdout, "user-controller ", log.LstdFlags))
	r.PathPrefix("/tenant").Handler(tenantHandler.NewRouter())

	return r
}

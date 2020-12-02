package router

import (
	"log"
	"os"

	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/services"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
)

func InitializeRouter(db db.DbFactory) *mux.Router {
	r := mux.NewRouter()

	itemStore := stores.NewItemStoreSql(db)
	itemLogger := log.New(os.Stdout, "item-controller ", log.LstdFlags)
	itemHandler := controller.NewItemHandler(itemStore, itemLogger)
	r.PathPrefix("/item").Handler(itemHandler.NewItemRouter())

	// user handler
	userStore := stores.NewUserStore(db)
	userLogger := log.New(os.Stdout, "user-controller ", log.LstdFlags)
	userHandler := controller.NewUserHandler(userStore, userLogger)
	r.PathPrefix("/user").Handler(userHandler.NewRouter())

	// tenant handler
	tenantStore := stores.NewTenantStore(db)
	tenantLogger := log.New(os.Stdout, "tenant-controller ", log.LstdFlags)
	tenantHandler := controller.NewTenantHandler(tenantStore, tenantLogger)
	r.PathPrefix("/tenant").Handler(tenantHandler.NewRouter())

	// auth handler
	authStore := stores.NewAuthStoreSql(db)
	authService := services.NewAuthHandler()
	authLogger := log.New(os.Stdout, "auth-controller", log.LstdFlags)
	authHandler := controller.NewAuthHandler(authStore, authService, authLogger)
	r.PathPrefix("/auth").Handler(authHandler.NewRouter())

	return r
}

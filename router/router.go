package router

import (
	"log"
	"os"

	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/db"
	"github.com/alesbrelih/go-reservation-api/middleware"
	"github.com/alesbrelih/go-reservation-api/services"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

func InitializeRouter(db db.DbFactory) *mux.Router {
	r := mux.NewRouter()

	authService := services.NewAuthHandler()

	jwt := middleware.NewJwt(authService, hclog.Default())

	itemStore := stores.NewItemStoreSql(db)
	itemLogger := log.New(os.Stdout, "item-controller ", log.LstdFlags)
	itemHandler := controller.NewItemHandler(itemStore, itemLogger)
	itemRouter := itemHandler.NewItemRouter()
	itemRouter.Use(jwt.ValidateUser)
	r.PathPrefix("/item").Handler(itemRouter)

	// user handler
	userStore := stores.NewUserStore(db)
	userLogger := log.New(os.Stdout, "user-controller ", log.LstdFlags)
	userHandler := controller.NewUserHandler(userStore, userLogger)
	userRouter := userHandler.NewRouter()
	userRouter.Use(jwt.ValidateUser)
	r.PathPrefix("/user").Handler(userRouter)

	// tenant handler
	tenantStore := stores.NewTenantStore(db)
	tenantLogger := log.New(os.Stdout, "tenant-controller ", log.LstdFlags)
	tenantHandler := controller.NewTenantHandler(tenantStore, tenantLogger)
	tenantRouter := tenantHandler.NewRouter()
	tenantRouter.Use(jwt.ValidateUser)
	r.PathPrefix("/tenant").Handler(tenantRouter)

	// auth handler
	authStore := stores.NewAuthStoreSql(db)
	authLogger := log.New(os.Stdout, "auth-controller", log.LstdFlags)
	authHandler := controller.NewAuthHandler(authStore, authService, authLogger)
	r.PathPrefix("/auth").Handler(authHandler.NewRouter())

	return r
}

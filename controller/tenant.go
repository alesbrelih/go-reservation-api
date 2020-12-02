package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/middleware"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
)

func NewTenantHandler(store stores.TenantStore, log *log.Logger) TenantHandler {
	return &tenantHandler{
		store: store,
		log:   log,
	}
}

type TenantHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	NewRouter() *mux.Router
}

type tenantHandler struct {
	log   *log.Logger
	store stores.TenantStore
}

func (h *tenantHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.store.GetAll(r.Context())
	if err != nil {
		h.log.Printf("Error retrieving tenants: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tenants.ToJSON(w)
}

func (h *tenantHandler) GetOne(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"]) // validated by regex already

	tenant, err := h.store.GetOne(r.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		h.log.Printf("Error retrieving tenant with id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tenant.ToJSON(w)
}

func (c *tenantHandler) Create(w http.ResponseWriter, r *http.Request) {

	tenant := r.Context().Value(&middleware.TenantBodyContextKey{}).(*models.Tenant)

	// validate
	createValidate.SetTagName("create")
	err := createValidate.Struct(tenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.store.Create(tenant)
	if err != nil {
		c.log.Printf("Error creating tenant: %#v. Error: %v", tenant, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, id)
}

func (c *tenantHandler) Update(w http.ResponseWriter, r *http.Request) {

	tenant := r.Context().Value(&middleware.TenantBodyContextKey{}).(*models.Tenant)

	// validate
	// TODO: shouldn tgo through
	updateValidate.SetTagName("update")
	err := updateValidate.Struct(tenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.store.Update(tenant)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		c.log.Printf("Error updating tenant: %#v. Error: %v", tenant, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (c *tenantHandler) Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"]) // validated by regex already

	err := c.store.Delete(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		c.log.Printf("Error deleting tenant with id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (c *tenantHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewTenantMiddleware(log.New(os.Stdout, "tenant-middleware ", log.LstdFlags))

	get := r.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/tenant", c.GetAll)
	get.HandleFunc("/tenant/{id:[\\d]+}", c.GetOne)

	post := r.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/tenant", c.Create)
	post.Use(middleware.GetBody)

	put := r.Methods(http.MethodPut).Subrouter()
	put.HandleFunc("/tenant", c.Update)
	put.Use(middleware.GetBody)

	delete := r.Methods(http.MethodDelete).Subrouter()
	delete.HandleFunc("/tenant/{id:[\\d]+}", c.Delete)

	return r
}

package controller

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/middleware"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/gorilla/mux"
)

type TenantStore interface {
	GetAll(ctx context.Context) (models.Tenants, error)
	GetOne(ctx context.Context, id int64) (*models.Tenant, error)
	Create(*models.Tenant) (int64, error)
	Update(*models.Tenant) error
	Delete(id int64) error
}

type TenantHandler struct {
	log   *log.Logger
	store TenantStore
}

func (h *TenantHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.store.GetAll(r.Context())
	if err != nil {
		h.log.Printf("Error retrieving tenants: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tenants.ToJSON(w)
}

func (h *TenantHandler) getOne(w http.ResponseWriter, r *http.Request) {

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

func (c *TenantHandler) create(w http.ResponseWriter, r *http.Request) {

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

func (c *TenantHandler) update(w http.ResponseWriter, r *http.Request) {

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

func (c *TenantHandler) delete(w http.ResponseWriter, r *http.Request) {

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

func (h *TenantHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewTenantMiddleware(log.New(os.Stdout, "tenant-middleware ", log.LstdFlags))

	get := r.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/tenant", h.getAll)
	get.HandleFunc("/tenant/{id:[\\d]+}", h.getOne)

	post := r.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/tenant", h.create)
	post.Use(middleware.GetBody)

	put := r.Methods(http.MethodPut).Subrouter()
	put.HandleFunc("/tenant", h.update)
	put.Use(middleware.GetBody)

	delete := r.Methods(http.MethodDelete).Subrouter()
	delete.HandleFunc("/tenant/{id:[\\d]+}", h.delete)

	return r
}

func NewTenantHandler(store TenantStore, log *log.Logger) *TenantHandler {
	return &TenantHandler{
		store: store,
		log:   log,
	}
}

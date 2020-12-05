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
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var baseValidate *validator.Validate
var createValidate *validator.Validate
var updateValidate *validator.Validate

func init() {
	baseValidate = validator.New()

	createValidate = validator.New()
	createValidate.SetTagName("create")

	updateValidate = validator.New()
	updateValidate.SetTagName("update")
}

func NewItemHandler(store stores.ItemStore, log *log.Logger) ItemHandler {
	return &itemHandler{
		store: store,
		log:   log,
	}
}

type ItemHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	NewItemRouter() *mux.Router
}

type itemHandler struct {
	log   *log.Logger
	store stores.ItemStore
}

func (h *itemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.GetAll(r.Context())
	if err != nil {
		h.log.Printf("Error retrieving items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	items.ToJSON(w)
}

func (h *itemHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"]) // validated by regex already

	item, err := h.store.GetOne(r.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		h.log.Printf("Error retrieving item with id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	item.ToJSON(w)
}

func (h *itemHandler) Create(w http.ResponseWriter, req *http.Request) {

	// .( ) <- type assertion
	item := req.Context().Value(&middleware.ItemBodyKeyType{}).(*models.Item)

	createValidate.SetTagName("create")
	err := createValidate.Struct(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.store.Create(req.Context(), item)

	if err != nil {
		h.log.Printf("Error creating item: %#v. Error: %v", item, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, id)
}

func (h *itemHandler) Update(w http.ResponseWriter, r *http.Request) {

	item := r.Context().Value(&middleware.ItemBodyKeyType{}).(*models.Item)

	updateValidate.SetTagName("update")
	err := updateValidate.Struct(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.store.Update(r.Context(), item)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		h.log.Printf("Error updating item: %#v. Error: %v", item, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (h *itemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	err := h.store.Delete(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		h.log.Printf("Error deleting item with id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *itemHandler) NewItemRouter() *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewItemMiddleware(log.New(os.Stdout, "item-middleware ", log.LstdFlags))

	getSubrouter := r.Methods(http.MethodGet).Subrouter()
	getSubrouter.HandleFunc("/item", h.GetAll)
	getSubrouter.HandleFunc("/item/{id:[\\d]+}", h.GetOne)

	postSubrouter := r.Methods(http.MethodPost).Subrouter()
	postSubrouter.HandleFunc("/item", h.Create)
	postSubrouter.Use(middleware.GetBody)

	putSubrouter := r.Methods(http.MethodPut).Subrouter()
	putSubrouter.HandleFunc("/item", h.Update)
	putSubrouter.Use(middleware.GetBody)

	deleteSubgrouter := r.Methods(http.MethodDelete).Subrouter()
	deleteSubgrouter.HandleFunc("/item/{id:[\\d]+}", h.Delete)

	return r
}

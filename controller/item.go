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
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var createValidate *validator.Validate
var updateValidate *validator.Validate

func init() {
	createValidate = validator.New()
	createValidate.SetTagName("create")

	updateValidate = validator.New()
	updateValidate.SetTagName("update")
}

type ItemStore interface {
	GetAll(ctx context.Context) (models.Items, error)
	GetOne(ctx context.Context, id int64) (*models.Item, error)
	Create(ctx context.Context, item *models.Item) (int64, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(id int64) error
}

type ItemHandler struct {
	log   *log.Logger
	store ItemStore
}

func (h *ItemHandler) getAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.GetAll(r.Context())
	if err != nil {
		h.log.Printf("Error retrieving items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	items.ToJSON(w)

}

func (h *ItemHandler) getOne(w http.ResponseWriter, r *http.Request) {
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

func (h *ItemHandler) create(w http.ResponseWriter, req *http.Request) {

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

func (h *ItemHandler) update(w http.ResponseWriter, r *http.Request) {

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

func (h *ItemHandler) delete(w http.ResponseWriter, r *http.Request) {
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

func (h *ItemHandler) NewItemRouter() *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewItemMiddleware(log.New(os.Stdout, "item-middleware ", log.LstdFlags))

	getSubrouter := r.Methods(http.MethodGet).Subrouter()
	getSubrouter.HandleFunc("/item", h.getAll)
	getSubrouter.HandleFunc("/item/{id:[\\d]+}", h.getOne)

	postSubrouter := r.Methods(http.MethodPost).Subrouter()
	postSubrouter.HandleFunc("/item", h.create)
	postSubrouter.Use(middleware.GetBody)

	putSubrouter := r.Methods(http.MethodPut).Subrouter()
	putSubrouter.HandleFunc("/item", h.update)
	putSubrouter.Use(middleware.GetBody)

	deleteSubgrouter := r.Methods(http.MethodDelete).Subrouter()
	deleteSubgrouter.HandleFunc("/item/{id:[\\d]+}", h.delete)

	return r
}

func NewItemHandler(store ItemStore, log *log.Logger) *ItemHandler {
	return &ItemHandler{
		store: store,
		log:   log,
	}
}

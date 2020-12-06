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

func NewUserHandler(store stores.UserStore, log *log.Logger) UserHandler {
	return &userHandler{
		store: store,
		log:   log,
	}
}

type UserHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	NewRouter() *mux.Router
}

type userHandler struct {
	log   *log.Logger
	store stores.UserStore
}

func (c *userHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := c.store.GetAll(r.Context())
	if err != nil {
		c.log.Printf("Error retrieving users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	items.ToJSON(w)
}

func (c *userHandler) GetOne(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"]) // validated by regex already

	item, err := c.store.GetOne(r.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		c.log.Printf("Error retrieving user with id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	item.ToJSON(w)
}

func (c *userHandler) Create(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(&middleware.UserBodyContextKey{}).(*models.UserReqBody)

	// validate
	createValidate.SetTagName("create")
	err := createValidate.Struct(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.store.Create(user)
	if err != nil {
		c.log.Printf("Error creating user. Error: %v", user)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)
}

func (c *userHandler) Update(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(&middleware.UserBodyContextKey{}).(*models.UserReqBody)

	// validate
	// TODO: shouldn tgo through
	updateValidate.SetTagName("update")
	err := updateValidate.Struct(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.store.Update(user)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		c.log.Printf("Error updating user: %#v. Error: %v", user, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (c *userHandler) Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"]) // validated by regex already

	err := c.store.Delete(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		c.log.Printf("Error deleting user with id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (h *userHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewUserMiddleware(log.New(os.Stdout, "user-middleware ", log.LstdFlags))

	get := r.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/user", h.GetAll)
	get.HandleFunc("/user/{id:[\\d]+}", h.GetOne)

	post := r.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/user", h.Create)
	post.Use(middleware.GetBody)

	put := r.Methods(http.MethodPut).Subrouter()
	put.HandleFunc("/user", h.Update)
	put.Use(middleware.GetBody)

	delete := r.Methods(http.MethodDelete).Subrouter()
	delete.HandleFunc("/user/{id:[\\d]+}", h.Delete)

	return r
}

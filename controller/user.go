package controller

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/middleware"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/gorilla/mux"
)

type UserStore interface {
	GetAll(ctx context.Context) (models.Users, error)
	GetOne(ctx context.Context, id int64) (*models.User, error)
	Create(*models.UserReqBody) (int64, error)
	Update(*models.UserReqBody) error
	Delete(id int64) error
}

type DefaultUserController struct {
	log *log.Logger
}

func (c *UserHandler) getAll(w http.ResponseWriter, r *http.Request) {
	items, err := c.store.GetAll(r.Context())
	if err != nil {
		c.log.Printf("Error retrieving users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	items.ToJSON(w)
}

func (c *UserHandler) getOne(w http.ResponseWriter, r *http.Request) {

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

func (c *UserHandler) create(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(&middleware.UserBodyContextKey{}).(*models.UserReqBody)

	// validate
	Validate.SetTagName("create")
	err := Validate.Struct(user)
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
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, id)
}

func (c *UserHandler) update(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(&middleware.UserBodyContextKey{}).(*models.UserReqBody)

	// validate
	// TODO: shouldn tgo through
	Validate.SetTagName("update")
	err := Validate.Struct(user)
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

func (c *UserHandler) delete(w http.ResponseWriter, r *http.Request) {

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

type UserHandler struct {
	log   *log.Logger
	store UserStore
}

func (h *UserHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	middleware := middleware.NewUserMiddleware()

	get := r.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/user", h.getAll)
	get.HandleFunc("/user/{id:[\\d]+}", h.getOne)

	post := r.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/user", h.create)
	post.Use(middleware.GetBody)

	put := r.Methods(http.MethodPut).Subrouter()
	put.HandleFunc("/user", h.update)
	put.Use(middleware.GetBody)

	delete := r.Methods(http.MethodDelete).Subrouter()
	delete.HandleFunc("/user/{id:[\\d]+}", h.delete)

	return r
}

func NewUserHandler(store UserStore, log *log.Logger) *UserHandler {
	return &UserHandler{
		store: store,
		log:   log,
	}
}

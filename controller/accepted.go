package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

func NewAcceptedHandler(store stores.AcceptedStore, log hclog.Logger) AcceptedHandler {
	return &acceptedHandler{
		store: store,
		log:   log,
	}
}

type AcceptedHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	ProcessInquiry(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	NewRouter() *mux.Router
}

type acceptedHandler struct {
	log   hclog.Logger
	store stores.AcceptedStore
}

func (a *acceptedHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := a.store.GetAll(r.Context())
	if err != nil {
		a.log.Error("Error retrieving accepted list (controller)", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	items.ToJSON(w)
}

func (a *acceptedHandler) ProcessInquiry(w http.ResponseWriter, r *http.Request) {

	accepted := &models.Accepted{}
	if err := accepted.FromJSON(r.Body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := baseValidate.Struct(accepted); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := a.store.ProcessInquiry(r.Context(), accepted)
	if err != nil {
		a.log.Error("Error processing inquiry", err)
		return
	}

	models.NewIdResponse(id).ToJSON(w)
}

func (a *acceptedHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64) // no need to check due to mux route

	err := a.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		a.log.Error("Error deleting accepted from db. Id: ", id, " Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *acceptedHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	getSubrouter := r.Methods(http.MethodGet).Subrouter()
	getSubrouter.HandleFunc("/accepted", a.GetAll)

	postSubrouter := r.Methods(http.MethodPost).Subrouter()
	postSubrouter.HandleFunc("/accepted/process", a.ProcessInquiry)

	deleteSubrouter := r.Methods(http.MethodDelete).Subrouter()
	deleteSubrouter.HandleFunc("/accepted/{id:[\\d+]}", a.Delete)

	return r
}

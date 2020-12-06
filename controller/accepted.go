package controller

import (
	"net/http"

	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

func NewAcceptedHandler(store stores.AcceptedStore, log hclog.Logger) AcceptedHandler {
	return &acceptedHandler{
		store: store,
		log:   log,
	}
}

type AcceptedHandler interface {
	ProcessInquiry(w http.ResponseWriter, r *http.Request)
	NewRouter() *mux.Router
}

type acceptedHandler struct {
	log   hclog.Logger
	store stores.AcceptedStore
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

func (a *acceptedHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	postSubrouter := r.Methods(http.MethodPost).Subrouter()
	postSubrouter.HandleFunc("/accepted/process", a.ProcessInquiry)

	return r
}

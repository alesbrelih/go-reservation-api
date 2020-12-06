package controller

import (
	"net/http"

	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

func NewInquiryHandler(store stores.InquiryStore, log hclog.Logger) InquiryHandler {
	return &inquiryHandler{
		store: store,
		log:   log,
	}
}

type InquiryHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	NewRouter() *mux.Router
}

type inquiryHandler struct {
	log   hclog.Logger
	store stores.InquiryStore
}

func (i *inquiryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	inquiries, err := i.store.GetAll(r.Context())
	if err != nil {
		i.log.Error("Error retrieving inquiries", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	inquiries.ToJSON(w)
}

func (i *inquiryHandler) Create(w http.ResponseWriter, r *http.Request) {

	ic := &models.InquiryCreate{}
	err := ic.FromJSON(r.Body)
	defer r.Body.Close()

	if err != nil {
		i.log.Debug("Invalid request body to create inquiry. Error: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = baseValidate.Struct(ic)
	if err != nil {
		i.log.Debug("Validation failed for inquiry create. Request body: %v. Error: %v", ic, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = i.store.Create(r.Context(), ic)
	if err != nil {
		i.log.Debug("Error saving inquiry in database. Request body", ic, " Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (i *inquiryHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	get := r.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/inquiry", i.GetAll)

	post := r.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/inquiry", i.Create)

	return r
}

package item

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ItemController interface {
	GetAll(w http.ResponseWriter, req *http.Request)
	GetOne(w http.ResponseWriter, req *http.Request)
}

type DefaultItemController struct{}

func (h *DefaultItemController) GetAll(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("IS THIS THEORY"))
}

func (h *DefaultItemController) GetOne(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("IS THIS MAGIC"))
}

func Router(controller ItemController) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/item", controller.GetAll).Methods("GET")
	r.HandleFunc("/item/{id:[\\d]+}", controller.GetOne).Methods("GET")
	return r
}

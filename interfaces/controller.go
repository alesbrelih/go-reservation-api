package interfaces

import "net/http"

type Controller interface {
	GetAll(w http.ResponseWriter, req *http.Request)
	GetOne(w http.ResponseWriter, req *http.Request)
	Create(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

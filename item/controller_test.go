package item_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alesbrelih/go-reservation-api/item"
	"github.com/gorilla/mux"
)

type MyFakeController struct{}

func (h *MyFakeController) GetAll(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GetAll"))
}

func (h *MyFakeController) GetOne(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GetOne"))
}

func testRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/item").Handler(item.Router(&MyFakeController{}))
	return r
}

func TestGetOne_Routing(t *testing.T) {

	req, _ := http.NewRequest("GET", "/item", nil)
	res := httptest.NewRecorder()

	testRouter().ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "GetAll" {
		t.Errorf("Response body should be GetAll but got %s", res.Body.String())
	}

}

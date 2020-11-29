package controller_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"code.soquee.net/testlog"
	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type MyFakeItemStore struct {
	mock.Mock
}

func (h *MyFakeItemStore) GetAll(ctx context.Context) (models.Items, error) {
	args := h.Called(ctx)
	return args.Get(0).(models.Items), args.Error(1)
}

func (h *MyFakeItemStore) GetOne(ctx context.Context, id int64) (*models.Item, error) {
	args := h.Called(ctx, id)
	return args.Get(0).(*models.Item), args.Error(1)
}

func (h *MyFakeItemStore) Create(item *models.Item) (int64, error) {
	args := h.Called(item)
	return args.Get(0).(int64), args.Error(1)
}

func (h *MyFakeItemStore) Update(item *models.Item) error {
	args := h.Called(item)
	return args.Error(0)
}

func (h *MyFakeItemStore) Delete(id int64) error {
	args := h.Called(id)
	return args.Error(0)
}

func testRouter(store controller.ItemStore, t *testing.T) *mux.Router {
	r := mux.NewRouter()

	itemHandler := controller.NewItemHandler(store, testlog.New(t))
	r.PathPrefix("/item").Handler(itemHandler.NewItemRouter())
	return r
}

func TestItem_GetAll_JSONFromDb(t *testing.T) {

	// setup mocking
	// returned from "db"
	title := "Hello"
	showTo := time.Now()
	showFrom := time.Now().AddDate(1, 0, 0)

	mockedItems := models.Items{
		{
			Id:       1,
			Title:    &title,
			ShowTo:   &showFrom,
			ShowFrom: &showTo,
		},
	}

	itemStore := &MyFakeItemStore{}
	itemStore.On("GetAll", mock.Anything).Return(mockedItems, nil)
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("GET", "/item", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedItems.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}
}

func TestItem_GetAll_DbError(t *testing.T) {
	var items models.Items
	itemStore := &MyFakeItemStore{}
	itemStore.On("GetAll", mock.Anything).Return(items, errors.New("Some error"))
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("GET", "/item", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestItem_GetOne_GetResult(t *testing.T) {
	// setup mocking
	// returned from "db"
	title := "Hello"
	showTo := time.Now()
	showFrom := time.Now().AddDate(1, 0, 0)

	mockedItem := &models.Item{
		Id:       1,
		Title:    &title,
		ShowTo:   &showFrom,
		ShowFrom: &showTo,
	}
	req, _ := http.NewRequest("GET", "/item/1", nil)

	itemStore := &MyFakeItemStore{}
	itemStore.On("GetOne", mock.Anything, int64(1)).Return(mockedItem, nil)
	router := testRouter(itemStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedItem.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}
}

func TestItem_GetOne_RoutingOnStringParameter(t *testing.T) {
	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("GET", "/item/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestItem_GetOne_RoutingOnMixedParameter(t *testing.T) {
	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("GET", "/item/12g1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestItem_Create_Success(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"title":"MyTitle","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("POST", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 201 {
		t.Errorf("Get all status code should be 201 but got %v", res.Result().StatusCode)
	}

	json := &models.Item{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	itemStore.AssertCalled(t, "Create", json)
}

func TestItem_Create_BadRequest_TitleLength(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"title":"My","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("POST", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Item.Title' Error:Field validation for 'Title' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestItem_Create_BadRequest_TitleMissing(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("POST", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Item.Title' Error:Field validation for 'Title' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestItem_Create_DbError(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Create", mock.Anything).Return(int64(0), errors.New("Some error"))
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"title":"MyTitle","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("POST", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestItem_Update_Success(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Update", mock.Anything).Return(nil)
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"id":1,"title":"MyTitle","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("PUT", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "" {
		t.Errorf("Excepted res body to be \"\" but got: %v", res.Body.String())
	}

	json := &models.Item{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	itemStore.AssertCalled(t, "Update", json)
}

func TestItem_Update_BadRequest_TitleLength(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"id": 1,"title":"My","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("PUT", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Item.Title' Error:Field validation for 'Title' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestItem_Update_BadRequest_TitleMissing(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"id":1,"ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("PUT", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Item.Title' Error:Field validation for 'Title' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestItem_Update_BadRequest_IdMissing(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"title":"hellothere","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("PUT", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Item.Id' Error:Field validation for 'Id' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestItem_Update_DbError(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Update", mock.Anything).Return(errors.New("Some error"))
	router := testRouter(itemStore, t)

	jsonStr := []byte(`{"id":5,"title":"MyTitle","ShowFrom":"2020-11-07T13:37:09.511Z","ShowTo":"2021-11-07T13:37:09.511Z"}`)
	req, _ := http.NewRequest("PUT", "/item", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestItem_Delete_GetResult(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/item/1", nil)

	itemStore := &MyFakeItemStore{}
	itemStore.On("Delete", int64(1)).Return(nil)
	router := testRouter(itemStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	expected := ""
	if res.Body.String() != expected {
		t.Errorf("Response body should be %#v but got %#v", expected, res.Body.String())
	}
}

func TestItem_Delete_RoutingOnStringParameter(t *testing.T) {
	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("DELETE", "/item/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestItem_Delete_RoutingOnMixedParameter(t *testing.T) {
	itemStore := &MyFakeItemStore{}
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("DELETE", "/item/12g1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestItem_Delete_DbError(t *testing.T) {

	itemStore := &MyFakeItemStore{}
	itemStore.On("Delete", int64(1)).Return(errors.New("Some error"))
	router := testRouter(itemStore, t)

	req, _ := http.NewRequest("DELETE", "/item/1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

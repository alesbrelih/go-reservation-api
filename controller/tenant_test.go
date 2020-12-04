package controller_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.soquee.net/testlog"
	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type MyFakeTenantStore struct {
	mock.Mock
}

func (h *MyFakeTenantStore) GetAll(ctx context.Context) (models.Tenants, error) {
	args := h.Called(ctx)
	return args.Get(0).(models.Tenants), args.Error(1)
}

func (h *MyFakeTenantStore) GetOne(ctx context.Context, id int64) (*models.Tenant, error) {
	args := h.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (h *MyFakeTenantStore) Create(item *models.Tenant) (int64, error) {
	args := h.Called(item)
	return args.Get(0).(int64), args.Error(1)
}

func (h *MyFakeTenantStore) Update(item *models.Tenant) error {
	args := h.Called(item)
	return args.Error(0)
}

func (h *MyFakeTenantStore) Delete(id int64) error {
	args := h.Called(id)
	return args.Error(0)
}

func tenantTestRouter(store stores.TenantStore, t *testing.T) *mux.Router {
	r := mux.NewRouter()

	tenantHandler := controller.NewTenantHandler(store, testlog.New(t))
	r.PathPrefix("/tenant").Handler(tenantHandler.NewRouter())
	return r
}

func TestTenant_GetAll_JSONFromDb(t *testing.T) {

	// setup mocking
	// returned from "db"
	mockedTenants := models.Tenants{
		{
			Id:    1,
			Title: "hello",
			Email: "hello@hello.email",
		},
		{
			Id:    2,
			Title: "byebye",
			Email: "bye@hello.email",
		},
	}

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("GetAll", mock.Anything).Return(mockedTenants, nil)
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("GET", "/tenant", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedTenants.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}
}

func TestTenant_GetAll_DbError(t *testing.T) {
	var tenants models.Tenants
	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("GetAll", mock.Anything).Return(tenants, errors.New("Some error"))
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("GET", "/tenant", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestTenant_GetOne_GetResult(t *testing.T) {
	// setup mocking
	// returned from "db"
	mockedTenant := &models.Tenant{
		Id:    1,
		Title: "hello",
		Email: "hello@hello.email",
	}
	req, _ := http.NewRequest("GET", "/tenant/1", nil)

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("GetOne", mock.Anything, int64(1)).Return(mockedTenant, nil)
	router := tenantTestRouter(tenantStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedTenant.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}

	tenantStore.AssertExpectations(t)
}

func TestTenant_GetOne_ErrorNoRows(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tenant/1", nil)
	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("GetOne", mock.Anything, int64(1)).Return(nil, sql.ErrNoRows)
	router := tenantTestRouter(tenantStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Fatalf("Error code should be 400, but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Bad request\n" {
		t.Fatalf("Body should be Bad request, but got %v", res.Body.String())
	}
}

func TestTenant_GetOne_SomeError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tenant/1", nil)
	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("GetOne", mock.Anything, int64(1)).Return(nil, errors.New("Some error"))
	router := tenantTestRouter(tenantStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Fatalf("Error code should be 500, but got %v", res.Result().StatusCode)
	}
	if res.Body.String() != "Internal server error\n" {
		t.Fatalf("Body should be Internal server error, but got %v", res.Body.String())
	}
}

func TestTenant_GetOne_RoutingOnStringParameter(t *testing.T) {
	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("GET", "/tenant/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestTenant_GetOne_RoutingOnMixedParameter(t *testing.T) {
	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("GET", "/tenant/12g1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestTenant_Create_Success(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"title":"my-tenant","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 201 {
		t.Errorf("Get all status code should be 201 but got %v", res.Result().StatusCode)
	}

	json := &models.Tenant{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	tenantStore.AssertCalled(t, "Create", json)
}

func TestTenant_Create_BadRequest_TitleLength(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"title":"my","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestTenant_Create_BadRequest_TitleMissing(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"title":null,"email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestTenant_Create_DbError(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Create", mock.Anything).Return(int64(0), errors.New("Some error"))
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"title":"my-supertitle","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestTenant_Update_Success(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Update", mock.Anything).Return(nil)
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"id":1,"title":"my-supertitle","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "" {
		t.Errorf("Excepted res body to be \"\" but got: %v", res.Body.String())
	}

	json := &models.Tenant{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	tenantStore.AssertCalled(t, "Update", json)
}

func TestTenant_Update_BadRequest_TitleLength(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"id":1,"title":"my","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestTenant_Update_BadRequest_TitleMissing(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"id":1,"title":null,"email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestTenant_Update_BadRequest_IdMissing(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"title":"my-supertitle","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'Tenant.Id' Error:Field validation for 'Id' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestTenant_Update_DbError(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Update", mock.Anything).Return(errors.New("Some error"))
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"id":1,"title":"my-supertitle","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestTenant_Update_NoRowsError(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Update", mock.Anything).Return(sql.ErrNoRows)
	router := tenantTestRouter(tenantStore, t)

	jsonStr := []byte(`{"id":1,"title":"my-supertitle","email":"mytenant@tenant.com"}`)
	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Bad request\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestTenant_Delete_GetResult(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/tenant/1", nil)

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Delete", int64(1)).Return(nil)
	router := tenantTestRouter(tenantStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	expected := ""
	if res.Body.String() != expected {
		t.Errorf("Response body should be %#v but got %#v", expected, res.Body.String())
	}

	tenantStore.AssertExpectations(t)
}

func TestTenant_Delete_RoutingOnStringParameter(t *testing.T) {
	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("DELETE", "/tenant/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestTenant_Delete_RoutingOnMixedParameter(t *testing.T) {
	tenantStore := &MyFakeTenantStore{}
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("DELETE", "/tenant/12g1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestTenant_Delete_DbError(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Delete", int64(1)).Return(errors.New("Some error"))
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("DELETE", "/tenant/1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestTenant_Delete_SqlNoRows(t *testing.T) {

	tenantStore := &MyFakeTenantStore{}
	tenantStore.On("Delete", int64(1)).Return(sql.ErrNoRows)
	router := tenantTestRouter(tenantStore, t)

	req, _ := http.NewRequest("DELETE", "/tenant/1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Bad request\n" {
		t.Errorf("Response body should be %#v but got %#v", "Bad request", res.Body.String())
	}
}

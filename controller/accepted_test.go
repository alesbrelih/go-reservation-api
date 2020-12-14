package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alesbrelih/go-reservation-api/controller"
	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/alesbrelih/go-reservation-api/test_util"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/mock"
)

type MyFakeAcceptedStore struct {
	mock.Mock
}

func (h *MyFakeAcceptedStore) GetAll(ctx context.Context) (models.AcceptedList, error) {
	args := h.Called(ctx)
	return args.Get(0).(models.AcceptedList), args.Error(1)
}

func (h *MyFakeAcceptedStore) ProcessInquiry(ctx context.Context, accepted *models.Accepted) (int64, error) {
	// removed controller parameters because not necessary in this test
	// i would be using mock.Anything anyways
	args := h.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (h *MyFakeAcceptedStore) Delete(ctx context.Context, id int64) error {
	args := h.Called(ctx, id)
	return args.Error(0)
}

func acceptedTestRouter(store stores.AcceptedStore, log hclog.Logger, t *testing.T) *mux.Router {
	r := mux.NewRouter()

	tenantHandler := controller.NewAcceptedHandler(store, log)
	r.PathPrefix("/accepted").Handler(tenantHandler.NewRouter())
	return r
}

func TestAccepted_GetAll_JSONFromDb(t *testing.T) {
	logMock := &test_util.HcLogMock{}

	// setup mocking
	// returned from "db"
	now := time.Now()
	someday := time.Now().AddDate(0, 2, 0)
	mockedAccepted := models.AcceptedList{
		{
			Id:                 1,
			Inquirer:           "john doe",
			InquirerEmail:      "john.doe@doe.com",
			InquirerPhone:      "",
			InquirerComment:    "some comment",
			ItemId:             1,
			ItemTitle:          "Item Title",
			ItemPrice:          200,
			Notes:              "he needs something special",
			DateReservation:    &now,
			DateInquiryCreated: &someday,
			DateAccepted:       &now,
		},
		{
			Id:                 2,
			Inquirer:           "jane buck",
			InquirerEmail:      "jane.buck@buck.com",
			InquirerPhone:      "",
			InquirerComment:    "some buck",
			ItemId:             1,
			ItemTitle:          "buck Title",
			ItemPrice:          200,
			Notes:              "he needs something buck",
			DateReservation:    &someday,
			DateInquiryCreated: &now,
			DateAccepted:       &someday,
		},
	}

	acceptedStore := &MyFakeAcceptedStore{}
	acceptedStore.On("GetAll", mock.Anything).Return(mockedAccepted, nil)
	router := acceptedTestRouter(acceptedStore, logMock, t)

	req, _ := http.NewRequest("GET", "/accepted", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedAccepted.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}
}

func TestAccepted_GetAll_DbError(t *testing.T) {
	err := errors.New("Some error")

	logMock := &test_util.HcLogMock{}
	logMock.On("Error", "Error retrieving accepted list (controller)", mock.Anything)

	var acceptedList models.AcceptedList
	acceptedStore := &MyFakeAcceptedStore{}
	acceptedStore.On("GetAll", mock.Anything).Return(acceptedList, err)
	router := acceptedTestRouter(acceptedStore, logMock, t)

	req, _ := http.NewRequest("GET", "/accepted", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}

	logMock.AssertExpectations(t)
}

func TestAccepted_ProcessInquiry_Success(t *testing.T) {
	logMock := &test_util.HcLogMock{}

	someTime := time.Now()
	mockedBody := &models.Accepted{
		ItemId:             1,
		Inquirer:           "john Doe",
		InquirerEmail:      "john.doe@doe.com",
		InquirerPhone:      "+38641666666",
		InquirerComment:    "My Comment",
		ItemTitle:          "Some item",
		ItemPrice:          4000,
		Notes:              "Some notes",
		DateReservation:    &someTime,
		DateInquiryCreated: &someTime,
		DateAccepted:       &someTime,
	}

	acceptedStore := &MyFakeAcceptedStore{}
	acceptedStore.On("ProcessInquiry").Return(int64(1), nil)
	router := acceptedTestRouter(acceptedStore, logMock, t)

	body, err := json.Marshal(mockedBody)
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("POST", "/accepted/process", bytes.NewReader(body))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	responseExpectation, _ := json.Marshal(models.NewIdResponse(1))
	if strings.TrimSpace(res.Body.String()) != string(responseExpectation) {
		t.Errorf("Response body should be %#v but got %#v", "{\"id\":1}\n", res.Body.String())
	}
}

// func TestTenant_GetOne_GetResult(t *testing.T) {
// 	// setup mocking
// 	// returned from "db"
// 	mockedTenant := &models.Tenant{
// 		Id:    1,
// 		Title: "hello",
// 		Email: "hello@hello.email",
// 	}
// 	req, _ := http.NewRequest("GET", "/tenant/1", nil)

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("GetOne", mock.Anything, int64(1)).Return(mockedTenant, nil)
// 	router := tenantTestRouter(tenantStore, t)

// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 200 {
// 		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
// 	}

// 	buf := new(bytes.Buffer)
// 	mockedTenant.ToJSON(buf)
// 	if res.Body.String() != buf.String() {
// 		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
// 	}

// 	tenantStore.AssertExpectations(t)
// }

// func TestTenant_GetOne_ErrorNoRows(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/tenant/1", nil)
// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("GetOne", mock.Anything, int64(1)).Return(nil, sql.ErrNoRows)
// 	router := tenantTestRouter(tenantStore, t)

// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Fatalf("Error code should be 400, but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "Bad request\n" {
// 		t.Fatalf("Body should be Bad request, but got %v", res.Body.String())
// 	}
// }

// func TestTenant_GetOne_SomeError(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/tenant/1", nil)
// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("GetOne", mock.Anything, int64(1)).Return(nil, errors.New("Some error"))
// 	router := tenantTestRouter(tenantStore, t)

// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 500 {
// 		t.Fatalf("Error code should be 500, but got %v", res.Result().StatusCode)
// 	}
// 	if res.Body.String() != "Internal server error\n" {
// 		t.Fatalf("Body should be Internal server error, but got %v", res.Body.String())
// 	}
// }

// func TestTenant_GetOne_RoutingOnStringParameter(t *testing.T) {
// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	req, _ := http.NewRequest("GET", "/tenant/hello", nil)
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 404 {
// 		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
// 	}
// }

// func TestTenant_GetOne_RoutingOnMixedParameter(t *testing.T) {
// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	req, _ := http.NewRequest("GET", "/tenant/12g1", nil)
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 404 {
// 		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
// 	}
// }

// func TestTenant_Create_Success(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Create", mock.Anything).Return(int64(1), nil)
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"title":"my-tenant","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 201 {
// 		t.Errorf("Get all status code should be 201 but got %v", res.Result().StatusCode)
// 	}

// 	json := &models.Tenant{}
// 	json.FromJSON(bytes.NewBuffer(jsonStr))
// 	tenantStore.AssertCalled(t, "Create", json)
// }

// func TestTenant_Create_BadRequest_TitleLength(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Create", mock.Anything).Return(int64(1), nil)
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"title":"my","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
// 	}

// 	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'gt' tag\n"
// 	if res.Body.String() != expectedErr {
// 		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
// 	}
// }

// func TestTenant_Create_BadRequest_TitleMissing(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Create", mock.Anything).Return(int64(1), nil)
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"title":null,"email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
// 	}

// 	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'required' tag\n"
// 	if res.Body.String() != expectedErr {
// 		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
// 	}
// }

// func TestTenant_Create_DbError(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Create", mock.Anything).Return(int64(0), errors.New("Some error"))
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"title":"my-supertitle","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("POST", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 500 {
// 		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "Internal server error\n" {
// 		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
// 	}
// }

// func TestTenant_Update_Success(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Update", mock.Anything).Return(nil)
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"id":1,"title":"my-supertitle","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 200 {
// 		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "" {
// 		t.Errorf("Excepted res body to be \"\" but got: %v", res.Body.String())
// 	}

// 	json := &models.Tenant{}
// 	json.FromJSON(bytes.NewBuffer(jsonStr))
// 	tenantStore.AssertCalled(t, "Update", json)
// }

// func TestTenant_Update_BadRequest_TitleLength(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"id":1,"title":"my","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
// 	}

// 	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'gt' tag\n"
// 	if res.Body.String() != expectedErr {
// 		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
// 	}
// }

// func TestTenant_Update_BadRequest_TitleMissing(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"id":1,"title":null,"email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
// 	}

// 	expectedErr := "Key: 'Tenant.Title' Error:Field validation for 'Title' failed on the 'required' tag\n"
// 	if res.Body.String() != expectedErr {
// 		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
// 	}
// }

// func TestTenant_Update_BadRequest_IdMissing(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"title":"my-supertitle","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
// 	}

// 	expectedErr := "Key: 'Tenant.Id' Error:Field validation for 'Id' failed on the 'required' tag\n"
// 	if res.Body.String() != expectedErr {
// 		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
// 	}
// }

// func TestTenant_Update_DbError(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Update", mock.Anything).Return(errors.New("Some error"))
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"id":1,"title":"my-supertitle","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 500 {
// 		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "Internal server error\n" {
// 		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
// 	}
// }

// func TestTenant_Update_NoRowsError(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Update", mock.Anything).Return(sql.ErrNoRows)
// 	router := tenantTestRouter(tenantStore, t)

// 	jsonStr := []byte(`{"id":1,"title":"my-supertitle","email":"mytenant@tenant.com"}`)
// 	req, _ := http.NewRequest("PUT", "/tenant", bytes.NewBuffer(jsonStr))
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "Bad request\n" {
// 		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
// 	}
// }

// func TestTenant_Delete_GetResult(t *testing.T) {

// 	req, _ := http.NewRequest("DELETE", "/tenant/1", nil)

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Delete", int64(1)).Return(nil)
// 	router := tenantTestRouter(tenantStore, t)

// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 200 {
// 		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
// 	}

// 	expected := ""
// 	if res.Body.String() != expected {
// 		t.Errorf("Response body should be %#v but got %#v", expected, res.Body.String())
// 	}

// 	tenantStore.AssertExpectations(t)
// }

// func TestTenant_Delete_RoutingOnStringParameter(t *testing.T) {
// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	req, _ := http.NewRequest("DELETE", "/tenant/hello", nil)
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 404 {
// 		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
// 	}
// }

// func TestTenant_Delete_RoutingOnMixedParameter(t *testing.T) {
// 	tenantStore := &MyFakeTenantStore{}
// 	router := tenantTestRouter(tenantStore, t)

// 	req, _ := http.NewRequest("DELETE", "/tenant/12g1", nil)
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 404 {
// 		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
// 	}
// }

// func TestTenant_Delete_DbError(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Delete", int64(1)).Return(errors.New("Some error"))
// 	router := tenantTestRouter(tenantStore, t)

// 	req, _ := http.NewRequest("DELETE", "/tenant/1", nil)
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 500 {
// 		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "Internal server error\n" {
// 		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
// 	}
// }

// func TestTenant_Delete_SqlNoRows(t *testing.T) {

// 	tenantStore := &MyFakeTenantStore{}
// 	tenantStore.On("Delete", int64(1)).Return(sql.ErrNoRows)
// 	router := tenantTestRouter(tenantStore, t)

// 	req, _ := http.NewRequest("DELETE", "/tenant/1", nil)
// 	res := httptest.NewRecorder()

// 	router.ServeHTTP(res, req)

// 	if res.Result().StatusCode != 400 {
// 		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
// 	}

// 	if res.Body.String() != "Bad request\n" {
// 		t.Errorf("Response body should be %#v but got %#v", "Bad request", res.Body.String())
// 	}
// }

package controller_test

import (
	"bytes"
	"context"
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

type MyFakeUserStore struct {
	mock.Mock
}

func (h *MyFakeUserStore) GetAll(ctx context.Context) (models.Users, error) {
	args := h.Called(ctx)
	return args.Get(0).(models.Users), args.Error(1)
}

func (h *MyFakeUserStore) GetOne(ctx context.Context, id int64) (*models.User, error) {
	args := h.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (h *MyFakeUserStore) Create(user *models.UserReqBody) (int64, error) {
	args := h.Called(user)
	return args.Get(0).(int64), args.Error(1)
}

func (h *MyFakeUserStore) Update(user *models.UserReqBody) error {
	args := h.Called(user)
	return args.Error(0)
}

func (h *MyFakeUserStore) Delete(id int64) error {
	args := h.Called(id)
	return args.Error(0)
}

func userRouter(store stores.UserStore, t *testing.T) *mux.Router {
	r := mux.NewRouter()

	userHandler := controller.NewUserHandler(store, testlog.New(t))
	r.PathPrefix("/user").Handler(userHandler.NewRouter())
	return r
}

func TestUser_GetAll_JSONFromDb(t *testing.T) {

	// setup mocking
	// returned from "db"

	mockedUsers := models.Users{
		{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Username:  "john.doe",
			Email:     "john.doe@does.com",
		},
	}

	userStore := &MyFakeUserStore{}
	userStore.On("GetAll", mock.Anything).Return(mockedUsers, nil)
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedUsers.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}
}

func TestUser_GetAll_DbError(t *testing.T) {
	var users models.Users
	userStore := &MyFakeUserStore{}
	userStore.On("GetAll", mock.Anything).Return(users, errors.New("Some error"))
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestUser_GetOne_GetResult(t *testing.T) {
	// setup mocking
	// returned from "db"

	mockedUser := &models.User{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "john.doe",
		Email:     "john.doe@does.com",
	}
	req, _ := http.NewRequest("GET", "/user/1", nil)

	userStore := &MyFakeUserStore{}
	userStore.On("GetOne", mock.Anything, int64(1)).Return(mockedUser, nil)
	router := userRouter(userStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	buf := new(bytes.Buffer)
	mockedUser.ToJSON(buf)
	if res.Body.String() != buf.String() {
		t.Errorf("Response body should be %#v but got %#v", buf.String(), res.Body.String())
	}

	userStore.AssertExpectations(t)
}

func TestUser_GetOne_RoutingOnStringParameter(t *testing.T) {
	userStore := &MyFakeUserStore{}
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("GET", "/user/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestUser_GetOne_RoutingOnMixedParameter(t *testing.T) {
	userStore := &MyFakeUserStore{}
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("GET", "/user/12g1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestUser_Create_Success(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"john","lastName":"doe","username":"john.doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 201 {
		t.Errorf("Get all status code should be 201 but got %v", res.Result().StatusCode)
	}

	json := &models.UserReqBody{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	userStore.AssertCalled(t, "Create", json)
}

func TestUser_Create_BadRequest_MissingName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":null,"lastName":"doe","username":"john.doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_ShortName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"e","lastName":"doe","username":"john.doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.FirstName' Error:Field validation for 'FirstName' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_MissingLastName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName": null,"username":"john.doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.LastName' Error:Field validation for 'LastName' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_ShortLastName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"d","username":"john.doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.LastName' Error:Field validation for 'LastName' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_MissingUsername(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":null,"email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Username' Error:Field validation for 'Username' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_ShortUsername(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName": "Doe","username":"john.a","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Username' Error:Field validation for 'Username' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_MissingEmail(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":null,"password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Email' Error:Field validation for 'Email' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_BadEmail(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Email' Error:Field validation for 'Email' failed on the 'email' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}
func TestUser_Create_BadRequest_MissingPassword(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":null, "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'eqfield' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_ShortPassword(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"shors", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'gt' tag\nKey: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'eqfield' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_MissingConfirm(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"password", "confirm":null}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag\nKey: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_ShortConfirm(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"shorsts", "confirm":"passw"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag\nKey: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Create_BadRequest_PasswordDontMatch(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(1), nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"shorsts", "confirm":"passwaa"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag\nKey: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'eqfield' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}
func TestUser_Create_DbError(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Create", mock.Anything).Return(int64(0), errors.New("Some error"))
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestUser_Update_Success(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"john","lastName":"doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "" {
		t.Errorf("Excepted res body to be \"\" but got: %v", res.Body.String())
	}

	json := &models.UserReqBody{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	userStore.AssertCalled(t, "Update", json)
}

func TestUser_Update_Success_NoPass(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"john","lastName":"doe","email":"john.doe@does.com"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "" {
		t.Errorf("Excepted res body to be \"\" but got: %v", res.Body.String())
	}

	json := &models.UserReqBody{}
	json.FromJSON(bytes.NewBuffer(jsonStr))
	userStore.AssertCalled(t, "Update", json)
}
func TestUser_Update_BadRequest_MissingId(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"firstName":"John","lastName":"doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Id' Error:Field validation for 'Id' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_MissingName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":null,"lastName":"doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_ShortName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"e","lastName":"doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.FirstName' Error:Field validation for 'FirstName' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_MissingLastName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":null,"email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.LastName' Error:Field validation for 'LastName' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_ShortLastName(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"d","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.LastName' Error:Field validation for 'LastName' failed on the 'gt' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_MissingEmail(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":null,"password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Email' Error:Field validation for 'Email' failed on the 'required' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_BadEmail(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Email' Error:Field validation for 'Email' failed on the 'email' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_PasswordDontMatch(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"shorsts", "confirm":"passwaa"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag\nKey: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'eqfield' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}
func TestUser_Update_BadRequest_PasswordDontMatch_2(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":null, "confirm":"passwaaa"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Confirm' Error:Field validation for 'Confirm' failed on the 'eqfield' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_BadRequest_PasswordDontMatch_3(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(nil)
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"Doe","username":"john.doe@does.com","email":"john.doe@does.com","password":"aaaaaaa", "confirm":null}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 400 {
		t.Errorf("Get all status code should be 400 but got %v", res.Result().StatusCode)
	}

	expectedErr := "Key: 'UserReqBody.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag\n"
	if res.Body.String() != expectedErr {
		t.Errorf("Expectes body to be: %v but got: %v", expectedErr, res.Body.String())
	}
}

func TestUser_Update_DbError(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Update", mock.Anything).Return(errors.New("Some error"))
	router := userRouter(userStore, t)

	jsonStr := []byte(`{"id":1,"firstName":"John","lastName":"Doe","email":"john.doe@does.com","password":"password", "confirm":"password"}`)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonStr))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

func TestUser_Delete_GetResult(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/user/1", nil)

	userStore := &MyFakeUserStore{}
	userStore.On("Delete", int64(1)).Return(nil)
	router := userRouter(userStore, t)

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("Get all status code should be 200 but got %v", res.Result().StatusCode)
	}

	expected := ""
	if res.Body.String() != expected {
		t.Errorf("Response body should be %#v but got %#v", expected, res.Body.String())
	}

	userStore.AssertExpectations(t)
}

func TestUser_Delete_RoutingOnStringParameter(t *testing.T) {
	userStore := &MyFakeUserStore{}
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("DELETE", "/user/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestUser_Delete_RoutingOnMixedParameter(t *testing.T) {
	userStore := &MyFakeUserStore{}
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("DELETE", "/user/12g1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 404 {
		t.Errorf("Get all status code should be 404 but got %v", res.Result().StatusCode)
	}
}

func TestUser_Delete_DbError(t *testing.T) {

	userStore := &MyFakeUserStore{}
	userStore.On("Delete", int64(1)).Return(errors.New("Some error"))
	router := userRouter(userStore, t)

	req, _ := http.NewRequest("DELETE", "/user/1", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Result().StatusCode != 500 {
		t.Errorf("Get all status code should be 500 but got %v", res.Result().StatusCode)
	}

	if res.Body.String() != "Internal server error\n" {
		t.Errorf("Response body should be %#v but got %#v", "Internal server error", res.Body.String())
	}
}

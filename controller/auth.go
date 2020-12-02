package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/alesbrelih/go-reservation-api/services"
	"github.com/alesbrelih/go-reservation-api/stores"
	"github.com/gorilla/mux"
)

func NewAuthHandler(authStore stores.AuthStore, authService services.AuthService, log *log.Logger) AuthHandler {
	return &authHandler{
		authStore:   authStore,
		authService: authService,
		log:         log,
	}
}

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	NewRouter() *mux.Router
}

type authHandler struct {
	log         *log.Logger
	authStore   stores.AuthStore
	authService services.AuthService
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {

	login := &models.LoginRequest{}
	err := login.FromJSON(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, err := a.authStore.Authenticate(r.Context(), login.Username, login.Password)
	if err != nil {
		if errors.Unwrap(errors.Unwrap(err)) == sql.ErrNoRows {
			http.Error(w, "Invalid authentication", http.StatusBadRequest)
			return
		}
		a.log.Printf("Error authenticating user login request: %#v. Error: %v", login, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tokenPair, err := a.authService.GenerateJwtPair(strconv.FormatInt(id, 10))
	if err != nil {
		a.log.Printf("Error generating jwt pair for id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tokenPair.ToJSON(w)
}

func (a *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken := &models.RefreshToken{}
	err := refreshToken.FromJSON(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	claims, err := a.authService.GetClaims(refreshToken.Refresh)
	if err != nil {
		if err == services.InvalidTokenError {
			a.log.Printf("Invalid token error: %v", err)
			http.Error(w, "Internal server error", http.StatusUnauthorized)
			return
		}
		a.log.Printf("Error decoding jwt: %v. Error: %v", refreshToken.Refresh, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		// TODO: CHECK TYPES OF ERRORS FROM DECODING -> COULD MEAN 401
		a.log.Printf("Error converting subject to integer. Subject: %v. Error: %v", claims.Subject, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = a.authStore.HasAccess(r.Context(), int64(id))
	if err != nil {
		http.Error(w, "Unathenticated", http.StatusUnauthorized)
		return
	}

	tokenPair, err := a.authService.GenerateJwtPair(claims.Subject)
	if err != nil {
		a.log.Printf("Error generating jwt pair for id: %v. Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tokenPair.ToJSON(w)
}

func (a *authHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()

	post := r.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/auth/login", a.Login)
	post.HandleFunc("/auth/refresh", a.RefreshToken)

	return r
}

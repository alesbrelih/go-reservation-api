package models

import (
	"encoding/json"
	"io"

	"github.com/dgrijalva/jwt-go"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *LoginRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(l)
}

type RefreshToken struct {
	Refresh string `json:"refresh"`
}

func (rt *RefreshToken) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(rt)
}

type TokenPair struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

func (t *TokenPair) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

type Claims struct {
	jwt.StandardClaims
}

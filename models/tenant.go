package models

import (
	"encoding/json"
	"io"
)

type Tenant struct {
	Id    int64  `json:"id" create:"omitempty" update:"required,number"`
	Title string `json:"title" create:"required,gt=3" update:"required,gt=3"`
	Email string `json:"email" create:"required,email" update:"required,email"`
	Users Users  `json:"users,omitempty"`
}

func (t *Tenant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(t)
}

func (t *Tenant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

type Tenants []Tenant

func (t *Tenants) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

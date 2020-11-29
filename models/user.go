package models

import (
	"encoding/json"
	"io"
)

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
}

type UserReqBody struct {
	Id        int64  `json:"id" create:"number,omitempty" update:"required,number"`
	FirstName string `json:"firstName" create:"required,gt=1" update:"required,gt=1"`
	LastName  string `json:"lastName" create:"required,gt=1" update:"required,gt=1"`
	Username  string `json:"username" create:"required,gt=6" update:"omitempty"`
	Email     string `json:"email" create:"required" update:"required"`
	Password  string `json:"password,omitempty" create:"required" update:"omitempty"`
	Confirm   string `json:"confirm,omitempty" create:"required" update:"omitempty"`
}

func (u *UserReqBody) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

type Users []User

func (u Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

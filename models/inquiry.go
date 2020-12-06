package models

import (
	"encoding/json"
	"io"
	"time"
)

type Inquiry struct {
	Id              int64     `json:"id,omitempty"`
	Inquirer        string    `json:"inquirer"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Item            Item      `json:"item"`
	DateReservation time.Time `json:"dateReservation"`
	DateCreated     time.Time `json:"dateCreated"`
	Comment         string    `json:"comment"`
}

type Inquiries []Inquiry

func (i Inquiries) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// TODO: on create check if date in future

type InquiryCreate struct {
	Inquirer string     `json:"inquirer" validate:"required,gt=2"`
	Email    string     `json:"email" validate:"omitempty,required_without=Phone,email"`
	Phone    string     `json:"phone" validate:"omitempty,required_without=Email,e164"`
	ItemId   int64      `json:"itemId" validate:"required"`
	Date     *time.Time `json:"date" validate:"required"`
	Comment  string     `json:"comment"`
}

func (ic *InquiryCreate) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ic)
}

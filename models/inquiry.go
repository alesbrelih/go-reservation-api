package models

import (
	"encoding/json"
	"io"
	"time"
)

type Inquiry struct {
	Id              int64 `json:"id,omitempty" update:"required"`
	Item            Item
	DateReservation time.Time
	DateCreated     time.Time
}

// TODO: on create check if date in future

type InquiryCreate struct {
	Inquirer string     `json:"inquirer" validate:"required,gt=2"`
	Email    string     `json:"email" validate:"omitempty,required_without=Phone,email"`
	Phone    string     `json:"phone" validate:"omitempty,required_without=Email,e164"`
	ItemId   int64      `json:"itemId" validate:"required"`
	Date     *time.Time `json:"date" validate:"required"`
}

func (ic *InquiryCreate) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ic)
}

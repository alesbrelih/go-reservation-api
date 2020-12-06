package models

import (
	"encoding/json"
	"io"
	"time"
)

// todo: price checking needs to check if int only!
type Accepted struct {
	Id                 int64     `json:"id" validate:"omitempty,required"`
	Inquirer           string    `json:"inquirer" validate:"required"`
	InquirerEmail      string    `json:"inquirerEmail" validate:"omitempty,required_without=Phone,email"`
	InquirerPhone      string    `json:"inquirerPhone" validate:"omitempty,required_without=Email,e164"`
	InquirerComment    string    `json:"inquirerComment"`
	ItemId             int64     `json:"itemId" validate:"omitempty,number,required_without=ItemTitle ItemPrice"`
	ItemTitle          string    `json:"itemTitle" validate:"omitempty,required_without=ItemId"`
	ItemPrice          int64     `json:"itemPrice" validate:"omitempty,number,required_without=ItemId"`
	Notes              string    `json:"notes"`
	DateReservation    time.Time `json:"dateReservation" validate:"required"`
	DateInquiryCreated time.Time `json:"dateInquiryCreated"`
	DateAccepted       time.Time `json:"dateAccepted,omitempty"`
}

func (a *Accepted) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Accepted) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

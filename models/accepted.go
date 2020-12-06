package models

import (
	"encoding/json"
	"io"
	"time"
)

// todo: price checking needs to check if int only!
// Implement custom marshaller for dates to avoid pointers
//  had to use pointer else omitempty doesnt work
type Accepted struct {
	Id                 int64      `json:"id,omitempty" validate:"omitempty,required"`
	Inquirer           string     `json:"inquirer,omitempty" validate:"required"`
	InquirerEmail      string     `json:"inquirerEmail,omitempty" validate:"omitempty,required_without=Phone,email"`
	InquirerPhone      string     `json:"inquirerPhone,omitempty" validate:"omitempty,required_without=Email,e164"`
	InquirerComment    string     `json:"inquirerComment,omitempty"`
	ItemId             int64      `json:"itemId,omitempty" validate:"omitempty,number,required_without=ItemTitle ItemPrice"`
	ItemTitle          string     `json:"itemTitle,omitempty" validate:"omitempty,required_without=ItemId"`
	ItemPrice          int64      `json:"itemPrice,omitempty" validate:"omitempty,number,required_without=ItemId"`
	Notes              string     `json:"notes,omitempty"`
	DateReservation    *time.Time `json:"dateReservation,omitempty" validate:"required"`
	DateInquiryCreated *time.Time `json:"dateInquiryCreated,omitempty"`
	DateAccepted       *time.Time `json:"dateAccepted,omitempty"`
}

func (a *Accepted) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Accepted) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

type AcceptedList []*Accepted

func (al *AcceptedList) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(al)
}

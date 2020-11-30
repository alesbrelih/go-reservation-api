package models

import (
	"encoding/json"
	"io"
	"time"
)

type Item struct {
	Id       int64      `json:"id" create:"number,omitempty" update:"required,number"`
	Title    *string    `json:"title" create:"required,gt=3" update:"required,gt=3"`
	ShowFrom *time.Time `json:"showFrom"`
	ShowTo   *time.Time `json:"showTo,omitempty"`
	Price    int64      `json:"price" create:"number,omitempty" update:"number,omitempty"`
}

func (i *Item) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func (i *Item) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

type Items []Item

func (i Items) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

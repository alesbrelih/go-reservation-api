package models

import (
	"database/sql"
	"encoding/json"
	"io"
)

type Item struct {
	Id       int64        `json:"id" validate:"number,omitempty"`
	Title    string       `json:"title" validate:"required,gt=3"`
	ShowFrom sql.NullTime `json:"showFrom"`
	ShowTo   sql.NullTime `json:"showTo"`
}

func (i *Item) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

type Items []Item

func (i Items) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

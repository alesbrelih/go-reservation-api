package models

import (
	"encoding/json"
	"io"
	"time"
)

type Item struct {
	Id       *int64     `json:"id" validate:"number"`
	Title    *string    `json:"title" validate:"required,gt=3"`
	ShowFrom *time.Time `json:"showFrom"`
	ShowTo   *time.Time `json:"showTo"`
}

func (i *Item) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

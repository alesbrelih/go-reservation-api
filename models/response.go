package models

import (
	"encoding/json"
	"io"
)

func NewIdResponse(id int64) *IdResponse {
	return &IdResponse{Id: id}
}

type IdResponse struct {
	Id int64 `json:"id"`
}

func (i *IdResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

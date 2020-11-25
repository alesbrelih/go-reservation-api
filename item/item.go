package item

import "time"

type Item struct {
	Id       int64      `json:"id"`
	Title    *string    `json:"title"`
	ShowFrom *time.Time `json:"showFrom"`
	ShowTo   *time.Time `json:"showTo"`
}

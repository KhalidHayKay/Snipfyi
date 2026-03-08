package model

import "time"

type Url struct {
	Id       int64     `json:"id"`
	Original string    `json:"original"`
	Short    string    `json:"short"`
	Visited  int64     `json:"visited"`
	Created  time.Time `json:"created"`
}

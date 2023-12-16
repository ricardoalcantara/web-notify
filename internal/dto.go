package internal

import "time"

//Simple message model
type Message struct {
	ID         int       `json:"id"`
	Sender     string    `json:"sender"`
	Text       string    `json:"text"`
	CreateTime time.Time `json:"createtime"`
}

package notification

import (
	"github.com/ricardoalcantara/web-notify/internal/models"
)

type Message struct {
	Id      uint
	Date    string
	Payload string
	Read    bool
}

func NewMessage(message *models.Message) Message {
	return Message{
		Id:      message.ID,
		Date:    message.CreatedAt.String(),
		Payload: message.Payload,
		Read:    message.Read,
	}

}

type EventPackt struct {
	Type    string
	UserId  string
	Message Message
}

type UserChannel = chan EventPackt

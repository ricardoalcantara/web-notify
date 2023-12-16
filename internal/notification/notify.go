package notification

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ricardoalcantara/web-notify/internal/models"
	"github.com/rs/zerolog/log"
)

var channels map[string][]UserChannel = make(map[string][]UserChannel)
var messageQueue chan EventPackt = make(chan EventPackt, 100)

func StartNotification() {
	go worker()
}

func worker() {
	log.Debug().Msg("StartNotification")

	var ids []uint
	for packt := range messageQueue {
		for _, channel := range channels[packt.UserId] {
			ids = append(ids, packt.Message.Id)
			log.Debug().Interface("message", packt).Send()
			channel <- packt
		}
		err := models.MessagesMarkAsNotified(ids)
		if err != nil {
			log.Error().Err(err).Stack().Send()
		}
	}
}

func AttachUser(userId string, c *gin.Context) UserChannel {
	userChannel := make(UserChannel)

	channels[userId] = append(channels[userId], userChannel)

	return userChannel
}

func Remove(userId string, u UserChannel) {
	for i, c := range channels[userId] {
		if c == u {
			slice := channels[userId]
			channels[userId] = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	close(u)
}

func NotifyPacket(packt EventPackt) error {
	select {
	case messageQueue <- packt:
		return nil
	default:
		return errors.New("max capacity reached")
	}
}

func Notify(userId string, eventType string, message Message) error {
	return NotifyPacket(EventPackt{
		Type:    eventType,
		UserId:  userId,
		Message: message,
	})
}

package sse

import (
	"github.com/gin-gonic/gin"
	"github.com/ricardoalcantara/web-notify/internal/models"
	"github.com/ricardoalcantara/web-notify/internal/notification"
	"github.com/ricardoalcantara/web-notify/internal/utils"
	"github.com/rs/zerolog/log"
)

func get(c *gin.Context) {

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	userId := c.GetString("x-user-id")

	log.Debug().Str("UserId", userId).Msg("Connected")
	userChannel := notification.AttachUser(userId, c)

	messages, err := models.MessageListNotNotified(userId, models.DefaultPagination())
	if err != nil {
		c.SSEvent("fail", utils.PrintError(err))
		return
	}

	var ids []uint
	for _, m := range messages {
		c.SSEvent(m.EventType, notification.NewMessage(&m))
		ids = append(ids, m.ID)
	}
	c.Writer.Flush()

	err = models.MessagesMarkAsNotified(ids)
	// models.MessagesUpdate(messages)
	if err != nil {
		c.SSEvent("fail", utils.PrintError(err))
		return
	}

channelLoop:
	for {
		select {
		case packt := <-userChannel:
			c.SSEvent(packt.Type, packt.Message)
			c.Writer.Flush()
		case <-c.Writer.CloseNotify():
			notification.Remove(userId, userChannel)
			break channelLoop
		}
	}

	log.Debug().Str("UserId", userId).Msg("Desconnected")
}

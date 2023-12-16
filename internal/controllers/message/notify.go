package message

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/ricardoalcantara/web-notify/internal/controllers"
	"github.com/ricardoalcantara/web-notify/internal/models"
	"github.com/ricardoalcantara/web-notify/internal/notification"
	"github.com/ricardoalcantara/web-notify/internal/utils"
)

func notify(c *gin.Context) {
	var input NofityMessage
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorResponse{Error: utils.PrintError(err)})
		return
	}

	messages := make([]models.Message, len(input.Users))

	for i, userId := range input.Users {
		message := models.Message{
			UserId:    userId,
			EventType: "message",
			Payload:   input.Payload,
			Read:      false,
		}
		messages[i] = message
	}

	err := models.MessagesSave(messages)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorResponse{Error: utils.PrintError(err)})
		return
	}

	for _, message := range messages {
		err := notification.Notify(message.UserId, message.EventType, notification.NewMessage(&message))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorResponse{Error: utils.PrintError(err)})
			return
		}
	}

	c.Status(http.StatusAccepted)
}

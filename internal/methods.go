package internal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MessageChannel = chan Message
type ClientChannels = map[string][]MessageChannel

var messageId = 0

// var messages
var clients ClientChannels = make(ClientChannels)

func newMessage(text string) Message {
	messageId += 1
	return Message{
		ID:         messageId,
		CreateTime: time.Now(),
		Text:       text,
	}
}

func StartPing() {
	go func() {
		for {
			msg := newMessage("Ping")
			for _, client := range clients {
				for _, sseChan := range client {
					sseChan <- msg
				}
			}

			time.Sleep(2 * time.Second)
		}
	}()
}

func SendMessage(c *gin.Context) {
	msg := newMessage(c.Query("text"))
	if to, exist := c.GetQuery("to"); exist {
		for _, sseChan := range clients[to] {
			sseChan <- msg
		}
	} else {
		for _, client := range clients {
			for _, sseChan := range client {
				sseChan <- msg
			}
		}
	}

	c.JSON(http.StatusOK, msg)
}

func ReceiveMessage(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	var userId string
	var exists bool
	if userId, exists = c.GetQuery("userId"); !exists {
		c.SSEvent("fail", "userId required")
		c.Writer.Flush()
		return
	}

	log.Debug().Str("userId", userId).Send()

	// Create a new channel for the new user.
	messageChan := make(chan Message)

	clients[userId] = append(clients[userId], messageChan)
	for {
		select {
		case sseMsg := <-messageChan:
			c.SSEvent("event", sseMsg)
			c.Writer.Flush()
		case <-c.Writer.CloseNotify():
			close(messageChan)
			clients = removeClient(clients, messageChan)
			return
		}
	}
}

func removeClient(clients ClientChannels, messageChannel MessageChannel) ClientChannels {
	for key, client := range clients {
		for i, clientChannel := range client {
			if clientChannel == messageChannel {
				client[i] = client[len(client)-1]
				clients[key] = client[:len(client)-1]
				return clients
			}
		}
	}
	return clients
}

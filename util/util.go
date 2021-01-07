package util

import (
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func WriteMsg(c *websocket.Conn, message []byte, logger *logrus.Logger) {
	err := c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		logger.Error("write message err:", err)
	}
}

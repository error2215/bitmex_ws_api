package ws

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/error2215/bitmex_ws_api/models"
	"github.com/error2215/bitmex_ws_api/storage"
	"github.com/error2215/bitmex_ws_api/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	logger := logrus.StandardLogger()
	logger = logger.WithField("id", util.RandString(8)).Logger

	subscribeCh := make(chan []byte)
	var subscribes []string

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("SubscribeHandler err on upgrader.Upgrade", err)
		return
	}
	defer c.Close()
	logger.Info("New client connected...")
	go sendSubscribedData(subscribeCh, c, logger)
	for {
		// read message from client
		_, message, err := c.ReadMessage()
		if err != nil {
			logger.Error("SubscribeHandler err on ReadMessage", err)
			util.WriteMsg(c, []byte("cannot read message"), logger)
			continue
		}
		// try to unmarshal message
		var msgStr models.WsMsg
		if err := json.Unmarshal(message, &msgStr); err != nil {
			logger.Error("SubscribeHandler err on json.Unmarshal", err)
			util.WriteMsg(c, []byte("bad message"), logger)
			continue
		}
		// depends on message type do some stuff
		switch msgStr.Action {
		case "subscribe":
			for _, symbol := range msgStr.Symbols {
				subscribes = append(subscribes, symbol)
				storage.AddSubscribe(subscribeCh, symbol, logger)

				//send last price
				msg := storage.ReadSymbol(symbol)
				data, err := json.Marshal(msg)
				if err != nil {
					logger.Errorln(err)
				}
				util.WriteMsg(c, data, logger)
			}
		case "unsubscribe":
			for _, symbol := range subscribes {
				storage.RemoveSubscribe(subscribeCh, symbol, logger)
			}
		// unsupported message
		default:
			util.WriteMsg(c, []byte("unsupported method"), logger)
		}
	}
}

func sendSubscribedData(subscribeCh chan []byte, c *websocket.Conn, logger *logrus.Logger) {
	for {
		select {
		case msg := <-subscribeCh:
			util.WriteMsg(c, msg, logger)
		}
	}
}

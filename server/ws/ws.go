package ws

import (
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"

	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}

type WSServer struct {
	port string
	*logrus.Logger
}

func NewWSServer(port string) *WSServer {
	return &WSServer{
		port: port,
	}
}

func (serv *WSServer) Start(wg *sync.WaitGroup) {
	serv.Logger = logrus.StandardLogger()
	r := gin.Default()

	r.GET("/", gin.WrapF(WSHandler))
	if err := r.Run(); err != nil {
		serv.Error(err)
	}
	defer wg.Done()
}

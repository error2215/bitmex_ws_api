package ws

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type WSServer struct {
	port   string
	logger *logrus.Entry
}

func NewWSServer(port string) *WSServer {
	return &WSServer{
		port: port,
	}
}

func (serv *WSServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
}

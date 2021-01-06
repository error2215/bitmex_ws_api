package http

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type HTTPServer struct {
	port   string
	logger *logrus.Entry
}

func NewHTTPServer(port string) *HTTPServer {
	return &HTTPServer{
		port: port,
	}
}

func (serv *HTTPServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
}

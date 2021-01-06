package server

import (
	"sync"
)

type Server interface {
	Start(wg *sync.WaitGroup)
}

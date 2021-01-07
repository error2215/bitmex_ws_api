package models

import (
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

type symbolList struct {
	sync.Mutex
	atomic.Value
}

var globalSymbolList symbolList

func init() {
	globalSymbolList.Store(make(map[string]BitmexSymbol))
}

func UpdateSymbol(key string, value BitmexSymbol, logger *logrus.Logger) {
	globalSymbolList.Lock()
	defer globalSymbolList.Unlock()
	m1 := globalSymbolList.Load().(map[string]BitmexSymbol)
	m2 := make(map[string]BitmexSymbol)
	for k, v := range m1 {
		m2[k] = v
	}
	m2[key] = value
	//logger.Logf(logrus.InfoLevel, "Updated %s: price - %f	. Number of %s subscribers - %d", key, value.Price, key, len(value.Clients))
	globalSymbolList.Store(m2)
}

func ReadSymbol(key string) BitmexSymbol {
	m1 := globalSymbolList.Load().(map[string]BitmexSymbol)
	return m1[key]
}

func AddSubscribe(channel chan []byte, key string, logger *logrus.Logger) {
	globalSymbolList.Lock()
	defer globalSymbolList.Unlock()
	m1 := globalSymbolList.Load().(map[string]BitmexSymbol)
	m2 := make(map[string]BitmexSymbol)
	for k, v := range m1 {
		m2[k] = v
	}
	newVal := m2[key]
	newVal.Clients = append(newVal.Clients, channel)
	m2[key] = newVal
	logger.Logf(logrus.InfoLevel, "New subscriber to %s. Number of subscribers - %d", key, len(m2[key].Clients))
	globalSymbolList.Store(m2)
}

func RemoveSubscribe(channel chan []byte, key string, logger *logrus.Logger) {
	globalSymbolList.Lock()
	defer globalSymbolList.Unlock()
	m1 := globalSymbolList.Load().(map[string]BitmexSymbol)
	m2 := make(map[string]BitmexSymbol)
	for k, v := range m1 {
		m2[k] = v
	}
	newVal := m2[key]

	for i, ch := range newVal.Clients {
		if channel == ch {
			//if len(newVal.Clients) == 1 {
			//	newVal.Clients = []chan<- []byte{}
			//}
			newVal.Clients = append(newVal.Clients[:i], newVal.Clients[i+1:]...)
		}
	}

	m2[key] = newVal
	logger.Logf(logrus.InfoLevel, "Subscriber removed to %s. Number of subscribers - %d", key, len(m2[key].Clients))
	globalSymbolList.Store(m2)
}

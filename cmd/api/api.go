package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/sirupsen/logrus"

	"github.com/error2215/bitmex_ws_api/config"
	"github.com/error2215/bitmex_ws_api/server"
	"github.com/error2215/bitmex_ws_api/server/http"
	"github.com/error2215/bitmex_ws_api/server/ws"
	"github.com/error2215/bitmex_ws_api/storage"
)

const bitmexURL = "testnet.bitmex.com"

type BitmexInstrumentAPI struct {
	*logrus.Logger
	servers map[config.ServerType]server.Server
	conf    *config.AppConfig
}

func defaultApiApp() *BitmexInstrumentAPI {
	return &BitmexInstrumentAPI{
		Logger:  logrus.New(),
		servers: map[config.ServerType]server.Server{},
	}
}

func main() {
	var wg sync.WaitGroup

	app := defaultApiApp()
	appConfig, err := config.ParseAppConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	app.conf = appConfig

	app.connectToBitmexInstrumentWS()

	for _, serverType := range appConfig.ServerTypes() {
		wg.Add(1)
		if serverType == config.HTTP {
			app.servers[serverType] = http.NewHTTPServer(app.conf.HTTPServerPort)
		}
		if serverType == config.WS {
			app.servers[serverType] = ws.NewWSServer(app.conf.WSServerPort)
		}

		app.servers[serverType].Start(&wg)
	}

	wg.Wait()
}

func (a *BitmexInstrumentAPI) connectToBitmexInstrumentWS() {
	a.Info("Connecting to Bitmex WS")
	bWS, err := websocket.Dial(fmt.Sprintf("wss://%s/realtime?subscribe=instrument", bitmexURL), "", fmt.Sprintf("http://%s/", bitmexURL))
	if err != nil {
		a.Fatalln(err.Error())
	}
	go a.readInstrumentMessages(bWS)
}

func (a *BitmexInstrumentAPI) readInstrumentMessages(ws *websocket.Conn) {
	for {
		var message interface{}
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			a.Errorln(err.Error())
			continue
		}
		if _, ok := message.(map[string]interface{})["data"]; ok {

			data := message.(map[string]interface{})["data"]

			for _, symbol := range data.([]interface{}) {
				if price, ok := symbol.(map[string]interface{})["lastPrice"].(float64); ok {
					name := symbol.(map[string]interface{})["symbol"].(string)

					symb := storage.ReadSymbol(name)
					symb.Symbol = name
					symb.Timestamp = symbol.(map[string]interface{})["timestamp"].(string)
					symb.Price = price

					storage.UpdateSymbol(name, symb, a.Logger)
					a.SendFreshData(name)
				}
			}
		}
	}
}

func (a *BitmexInstrumentAPI) SendFreshData(key string) {
	symbol := storage.ReadSymbol(key)
	data, err := json.Marshal(symbol)
	if err != nil {
		a.Errorln(err)
	}
	// send fresh information to all subscribers
	for _, client := range symbol.Clients {
		client <- data
	}
}

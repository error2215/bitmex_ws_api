package main

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/error2215/bitmex_ws_api/config"
	"github.com/error2215/bitmex_ws_api/server"
	"github.com/error2215/bitmex_ws_api/server/http"
	"github.com/error2215/bitmex_ws_api/server/ws"
)

type ApiApp struct {
	logrus.Logger
	servers map[config.ServerType]server.Server
	conf    *config.AppConfig
}

func defaultApiApp() *ApiApp {
	return &ApiApp{
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

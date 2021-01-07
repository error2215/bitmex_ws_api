package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	envFileLoadError = fmt.Errorf("Error loading .env file ")
)

type AppConfig struct {
	serverTypes    []ServerType
	HTTPServerPort string
	WSServerPort   string
}

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		serverTypes:    []ServerType{WS},
		WSServerPort:   "8080",
		HTTPServerPort: "8081",
	}
}

func ParseAppConfig() (*AppConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, envFileLoadError
	}

	appConf := DefaultAppConfig()
	if os.Getenv("HTTP_SERVER") == "true" {
		appConf.serverTypes = append(appConf.serverTypes, HTTP)
		appConf.HTTPServerPort = os.Getenv("HTTP_SERVER_PORT")
	}

	if os.Getenv("WS_SERVER") == "true" {
		appConf.serverTypes = append(appConf.serverTypes, WS)
		appConf.WSServerPort = os.Getenv("WS_SERVER_PORT")

	}

	return appConf, nil
}

func (conf *AppConfig) ServerTypes() []ServerType {
	return conf.serverTypes
}

type ServerType int

const (
	HTTP ServerType = iota
	WS
	GRPC
)

func (s ServerType) String() string {
	return [...]string{"Http", "WebSocket", "gRPC"}[s]
}

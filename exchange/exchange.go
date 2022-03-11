package exchange

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"nhooyr.io/websocket"
)

type ExchangeConfig struct {
    ApiKey    string `flag:"key||Exchange API key" env:"TRADEMAN_EXCHANGE_KEY"`
    ApiSecret string `flag:"secret||Exchange API secret" env:"TRADEMAN_EXCHANGE_SECRET"`
    RESTHost  string `flag:"rest_host||Exchange API REST protocol host" env:"TRADEMAN_EXCHANGE_REST_HOST"`
    WSHost    string `flag:"websocket_host||Exchange WebSocket protocol host" env:"TRADEMAN_EXCHANGE_WEBSOCKET_HOST"`

	// TODO: get logfile settings from config file
}

type Exchange struct {
	apiKey    string
	apiSecret string

	restHost      string
	restClient    *http.Client
	websocketHost string
	context       context.Context
	connection    *websocket.Conn

	logger    lumberjack.Logger
	debugMode bool
}

func Connect(c ExchangeConfig) (e *Exchange, err error) {
	e = &Exchange{
		apiKey:        c.ApiKey,
		apiSecret:     c.ApiSecret,
		restHost:      c.RESTHost,
		websocketHost: c.WSHost,
		restClient: &http.Client{
			Timeout: 6 * time.Second, // TODO check with documentation
		},
		context: context.Background(),

		debugMode: false,
	}

	e.connection, _, err = websocket.Dial(e.context, e.websocketHost, nil)
	if err != nil {
		return nil, err
	}

	e.logger.Filename = "./log/trademan.log"
	e.logger.MaxSize = 10000000
	e.logger.Compress = true

	err = e.Authenticate()

	if err == nil {
		e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Successfully connected.\n", time.Now().Format("2006-01-02 15:04:05.000"))))
	} else {
		e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Error connecting to exchange: %v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)))
	}

	return e, err
}

func (e *Exchange) Reconnect() (err error) {
	e.connection, _, err = websocket.Dial(e.context, e.websocketHost, nil)
	err = e.Authenticate()
	if err == nil {
		e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Successfully connected.\n", time.Now().Format("2006-01-02 15:04:05.000"))))
	} else {
		e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Error connecting to exchange: %v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)))
	}

	return err
}

func (e *Exchange) Close() {
	e.context.Done()
	err := e.connection.Close(http.StatusConflict, "weirdness")
	if err == nil {
		e.logger.Write([]byte(fmt.Sprintf("%s [server] Succesfully disconnected.\n", time.Now().Format("2006-01-02 15:04:05.000"))))
	} else {
		e.logger.Write([]byte(fmt.Sprintf("%s [server] Unsuccessfully disconnected?\n", time.Now().Format("2006-01-02 15:04:05.000"))))
	}
	e.logger.Close()
}

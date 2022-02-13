package exchange

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"nhooyr.io/websocket"
)

// TODO: set defaults and descriptions
type ExchangeConfig struct {
	ApiKey    string
	ApiSecret string
	RESTHost  string
	WSHost    string
}

type Exchange struct {
	apiKey    string
	apiSecret string

	restHost   string       // TODO: rename to lowercaps to avoid exporting
	restClient *http.Client // TODO: rename to lowercaps to avoid exporting

	websocketHost string
	context       context.Context
	connection    *websocket.Conn

	logger lumberjack.Logger

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
	}

	e.connection, _, err = websocket.Dial(e.context, e.websocketHost, nil)
	if err != nil {
		return nil, err
	}

	// TODO: get these (and other) settings from config file
	e.logger.Filename = "./logs/trademan.log"
	e.logger.MaxSize = 10000000
	e.logger.Compress = true

	err = e.Authenticate()

	if err == nil {
		e.logger.Write([]byte(fmt.Sprintf("%s [server] Succesfully connected\n", time.Now().Format("2006-01-02 15:04:05.000"))))
	} else {
		e.logger.Write([]byte(fmt.Sprintf("%s [server] Error connecting to exchange: %v\n", time.Now().Format("2006-01-02 15:04:05.000"), err)))
	}

	return e, err
}

func (e *Exchange) Close() {
	e.context.Done()
	err := e.connection.Close(http.StatusConflict, "weirdness")
	if err == nil { // TODO: but what if it isn't
		e.logger.Write([]byte(fmt.Sprintf("%s [server] Succesfully disconnected\n", time.Now().Format("2006-01-02 15:04:05.000"))))
	}
	e.logger.Close()
}

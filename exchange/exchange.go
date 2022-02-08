package exchange

import (
	"context"
	"net/http"
	"time"

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

	RESTHost   string       // TODO: rename to lowercaps to avoid exporting
	RESTClient *http.Client // TODO: rename to lowercaps to avoid exporting

	websocketHost string
	context       context.Context
	connection    *websocket.Conn

	debugMode bool
}

func Connect(c ExchangeConfig) (e Exchange, err error) {
	e = Exchange{
		apiKey:        c.ApiKey,
		apiSecret:     c.ApiSecret,
		RESTHost:      c.RESTHost,
		websocketHost: c.WSHost,
		RESTClient: &http.Client{
			Timeout: 6 * time.Second, // TODO check with documentation
		},
        context: context.Background(),
	}

	e.connection, _, err = websocket.Dial(e.context, e.websocketHost, nil)
	if err != nil {
		return Exchange{}, err
	}

	err = e.Authenticate()
    if err != nil {
		return Exchange{}, err
    }

	err = e.Ping()

	return e, err
}

package main

import (
	"context"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type ByBit struct {
	apiKey    string
	apiSecret string

	RESTHost   string
	RESTClient *http.Client

	websocketHost string
	context       context.Context
	connection    *websocket.Conn

	debugMode bool
}

func NewByBit(RESTHost, websocketHost string, apiKey, apiSecret string, debugMode bool) *ByBit {
	return &ByBit{
		apiKey:    apiKey,
		apiSecret: apiSecret,

		RESTHost: RESTHost,
		RESTClient: &http.Client{
			Timeout: 6 * time.Second, // TODO check with documentation
		},

		websocketHost: websocketHost,

		debugMode: debugMode,
	}
}

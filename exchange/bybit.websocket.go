package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"nhooyr.io/websocket/wsjson"
)

type websocketCmd struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type websocketResponse struct {
	Topic         string      `json:"topic"`
	Success       bool        `json:"success,omitempty"`
	ReturnMessage string      `json:"ret_msg,omitempty"`
	Action        string      `json:"action,omitempty"`
	Data          interface{} `json:"data"`
}

func (e *Exchange) Authenticate() error {
	expiresTime := time.Now().Unix()*1000 + 10000
	req := "GET/realtime" + strconv.FormatInt(expiresTime, 10)
	sig := hmac.New(sha256.New, []byte(e.apiSecret))
	sig.Write([]byte(req))
	signature := hex.EncodeToString(sig.Sum(nil))

	cmd := websocketCmd{Op: "auth",
		Args: []interface{}{e.apiKey,
			expiresTime,
			signature,
		},
	}
	err := wsjson.Write(e.context, e.connection, cmd)
	if err != nil {
		return err
	}
	_, data, err := e.connection.Read(e.context)
	if err != nil {
		return err
	}
	var response websocketResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	if response.Success != true {
        return err
	}
	return nil
}

func (e *Exchange) Ping() error {
	cmd := websocketCmd{Op: "ping"}
	err := wsjson.Write(e.context, e.connection, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exchange) Subscribe(topic string) error {
	cmd := websocketCmd{Op: "subscribe",
		Args: []interface{}{topic},
	}

	err := wsjson.Write(e.context, e.connection, cmd)
	if err != nil {
		return err
	}
	_, data, err := e.connection.Read(e.context)
	if err != nil {
		return err
	}
	var response websocketResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	if response.Success != true {
        return errors.New(response.ReturnMessage)
	}
	return nil
}

func (e *Exchange) Close() {
	e.context.Done()
	e.connection.Close(http.StatusConflict, "weirdness")
}

package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dez11de/cryptodb"
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

func (b *Exchange) Authenticate() error {
	expiresTime := time.Now().Unix()*1000 + 10000
	req := "GET/realtime" + strconv.FormatInt(expiresTime, 10)
	sig := hmac.New(sha256.New, []byte(b.apiSecret))
	sig.Write([]byte(req))
	signature := hex.EncodeToString(sig.Sum(nil))

	cmd := websocketCmd{Op: "auth",
		Args: []interface{}{b.apiKey,
			expiresTime,
			signature,
		},
	}
	err := wsjson.Write(b.context, b.connection, cmd)
	if err != nil {
		return err
	}
	_, data, err := b.connection.Read(b.context)
	if err != nil {
		return err
	}
	var response websocketResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	if response.Success != true {
		return fmt.Errorf("Unable to authenticate")
	}
	return nil
}

func (b *Exchange) Ping() error {
	cmd := websocketCmd{Op: "ping"}
	err := wsjson.Write(b.context, b.connection, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (b *Exchange) Subscribe(topic string) error {
	cmd := websocketCmd{Op: "subscribe",
		Args: []interface{}{topic},
	}

	err := wsjson.Write(b.context, b.connection, cmd)
	if err != nil {
		return err
	}
	_, data, err := b.connection.Read(b.context)
	if err != nil {
		return err
	}
	var response websocketResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	if response.Success != true {
		return fmt.Errorf("Unable to subscribe to topic %s", topic)
	}
	return nil
}

func (b *Exchange) ProcessMessages(positionsChannel chan<- Position, executionChannel chan<- Execution, orderChannel chan<- cryptodb.Order) {
	for {
		_, data, err := b.connection.Read(b.context)
		if err != nil {
			fmt.Printf("Error reading data from socket (%v)\n", err)
		}
		var rawData json.RawMessage
		wsresp := websocketResponse{
			Data: &rawData,
		}

		err = json.Unmarshal(data, &wsresp)
		if err != nil {
			fmt.Printf("Error reading response (%v)\n", err)
		}

		switch wsresp.Success {
		case true:
			switch wsresp.ReturnMessage {
			case "pong":
			}
		}

        if !wsresp.Success { // TODO: this doesn't seem right does it?
			var p []Position
			var e []Execution
			var o []cryptodb.Order
			switch wsresp.Topic {
			case "position":
				err = json.Unmarshal(rawData, &p)
				if err != nil {
					fmt.Printf("Error unmarshalling position from response (%+v)", err)
				}
				for _, positions := range p {
					positionsChannel <- positions
				}
			case "execution":
				err = json.Unmarshal(rawData, &e)
				if err != nil {
					fmt.Printf("Error unmarshalling execution from response (%+v)", err)
				}
				for _, executions := range e {
					executionChannel <- executions
				}
			case "order":
				err = json.Unmarshal(rawData, &o)
				if err != nil {
					fmt.Printf("Error unmarshalling execution from response (%+v)", err)
				}
				for _, orders := range o {
					orderChannel <- orders
				}
			default:
				fmt.Printf("Unknown message received (%s)\n", string(data))
			}
		}
	}
}

func (b *Exchange) Close() {
	b.context.Done()
	b.connection.Close(http.StatusConflict, "weirdness")
}

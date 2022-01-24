package exchange

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
    "github.com/dez11de/cryptodb"
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

func (b *ByBit) Connect() error {
	b.context = context.Background()
	var err error
	b.connection, _, err = websocket.Dial(b.context, b.websocketHost, nil)
	if err != nil {
		return err
	}
	return b.Authenticate()
}

func (b *ByBit) Authenticate() error {
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

func (b *ByBit) Ping() error {
	cmd := websocketCmd{Op: "ping"}
	err := wsjson.Write(b.context, b.connection, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (b *ByBit) Subscribe(topic string) error {
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

func (b *ByBit) ProcessMessages(positionsChannel chan<- cryptoDB.Plan, executionChannel chan<- cryptoDB.Execution, orderChannel chan<- cryptoDB.Order) {
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

		if !wsresp.Success {
			var p []cryptoDB.Plan
			var e []cryptoDB.Execution
			var o []cryptoDB.Order
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

func (b *ByBit) Close() {
	b.context.Done()
	b.connection.Close(http.StatusConflict, "weirdness")
}

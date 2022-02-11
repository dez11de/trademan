package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (e *Exchange) ProcessMessages(positionChannel chan<- Position, executionChannel chan<- Execution, orderChannel chan<- Order, errorChannel chan<- error) {
	for {
		_, data, err := e.connection.Read(e.context)
		if err != nil {
			errorChannel <- err
		}
		var rawData json.RawMessage
		wsresp := websocketResponse{
			Data: &rawData,
		}

		err = json.Unmarshal(data, &wsresp)
		if err != nil {
			errorChannel <- err
		}

		switch wsresp.Success {
		case true:
			switch wsresp.ReturnMessage {
			case "pong":
			}
		}

		if !wsresp.Success {
			var positions []Position
			var executions []Execution
			var orders []Order
			switch wsresp.Topic {
			case "position":
				err = json.Unmarshal(rawData, &positions)
				if err != nil {
					errorChannel <- err
				}
				for _, position := range positions {
					positionChannel <- position
				}
			case "execution":
				err = json.Unmarshal(rawData, &executions)
				if err != nil {
					errorChannel <- err
				}
				for _, execution := range executions {
					executionChannel <- execution
				}
			case "order", "stop_order":
				err = json.Unmarshal(rawData, &orders)
				if err != nil {
					errorChannel <- err
				}
				for _, order := range orders {
					orderChannel <- order
				}
			default:
				errorChannel <- errors.New(fmt.Sprintf("This should NEVER happen: Error: %s or unknown topic: %s\n", err, wsresp.Topic))
			}
		}
	}
}

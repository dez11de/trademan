package exchange

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (e *Exchange) ProcessMessages(positionChannel chan<- Position, executionChannel chan<- Execution, orderChannel chan<- Order) {
	for {
		_, data, err := e.connection.Read(e.context)
		if e.debugMode {
			e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Raw:       %v\n", time.Now().Format("2006-01-02 15:04:05.000"), string(data))))
		}

		if err != nil {
			log.Printf("Error reading websocket.")
		}

		var rawData json.RawMessage
		wsresp := websocketResponse{Data: &rawData}
		err = json.Unmarshal(data, &wsresp)
		if err != nil {
			log.Printf("Error unmarshalling response %s", err)
		}

		switch wsresp.Success {
		case true:
			switch wsresp.ReturnMessage {
			case "pong":
				// TODO: IF haven't received pong in 2 minutes (re)connect, ELSE Reset timer?
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
					log.Printf("Error unmarshalling position.")
				}
				for _, position := range positions {
					e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Position:  %v\n", time.Now().Format("2006-01-02 15:04:05.000"), position)))
					positionChannel <- position
				}
			case "execution":
				err = json.Unmarshal(rawData, &executions)
				if err != nil {
					log.Printf("Error unmarshalling execution.")
				}
				for _, execution := range executions {
					e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Execution: %v\n", time.Now().Format("2006-01-02 15:04:05.000"), execution)))
					executionChannel <- execution
				}
			case "order", "stop_order":
				err = json.Unmarshal(rawData, &orders)
				if err != nil {
					log.Printf("Error unmarshalling order %s", err)
                    log.Printf("rawData: %s", string(rawData))
				}
				for _, order := range orders {
					e.logger.Write([]byte(fmt.Sprintf("%s [exchange] Order:     %v\n", time.Now().Format("2006-01-02 15:04:05.000"), order)))
					orderChannel <- order
				}
			}
		}
	}
}

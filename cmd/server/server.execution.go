package main

import (
	"log"

	"github.com/dez11de/exchange"
)

func processExecution(e exchange.Execution) (err error) {
    log.Printf("Processing execution %+v", e)
	return err
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

// TODO: rewrite as gRPC
func allPlansHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var allPlans []cryptodb.Plan
	result := db.Find(&allPlans)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}
	jsonResp, err := json.Marshal(allPlans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

// TODO: rewrite as gRPC
func executePlanHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("ID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var plan cryptodb.Plan
	var pair cryptodb.Pair
	var orders []cryptodb.Order
	var balance cryptodb.Balance
	db.Where("id = ?", id).First(&plan)
	db.Where("id = ?", plan.PairID).First(&pair)
	db.Where("plan_id = ?", plan.ID).Find(&orders)
	db.Where("symbol = ?", pair.QuoteCurrency).Order("created_at DESC").First(&balance)
	ticker, _ := e.GetTicker(pair.Name)

	tx := db.Begin()
	tx.Create(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: "Finalized orders."})
	// TODO: this should handle an error
	plan.FinalizeOrders(balance.Available, pair, orders)
    result := tx.Save(&orders)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, result.Error)
		tx.Rollback()
		return
	}

	plan.Status = cryptodb.Ordered
    result = tx.Save(&plan)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, result.Error)
		tx.Rollback()
		return
	}

	jsonResp, err := json.Marshal(plan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		tx.Rollback()
		return
	}
    // TODO: check for error
	tx.Create(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: "Sent plan to exchange."})
	tx.Commit()

    log.Printf("Pausing main routine")
    pause <- struct{}{}
	err = PlaceOrders(plan, pair, ticker, orders)
    play <- struct{}{}
    log.Printf("Continue main routine")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		db.Create(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: fmt.Sprintf("Exchange did not accept plan. %s", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)

	return
}

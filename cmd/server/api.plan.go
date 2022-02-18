package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

func allPlansHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allPlans, err := db.GetPlans()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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

func executePlanHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("ID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	plan, _ := db.GetPlan(uint(id))
	pair, _ := db.GetPair(plan.PairID)
	orders, _ := db.GetOrders(plan.ID)
	balance, _ := db.GetCurrentBalance(pair.QuoteCurrency)
	ticker, _ := e.GetTicker(pair.Name)

	tx := db.Begin()
	db.CreateLog(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: "Finalized orders."})
	// TODO: this should handle an error
	plan.FinalizeOrders(balance.Available, pair, orders)
	err = db.SaveOrders(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		tx.Rollback()
		return
	}

	plan.Status = cryptodb.Ordered
	err = db.SavePlan(&plan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
	db.CreateLog(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: "Sent plan to exchange."})
	tx.Commit()

	err = e.PlaceOrders(plan, pair, ticker, orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		db.CreateLog(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: fmt.Sprintf("Exchange did not accept plan. %s", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)

	return
}

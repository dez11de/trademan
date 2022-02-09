package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

// TODO: should only return 'open' plans. See GORM api documentation.
func allPlansHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allPlans, err := db.GetPlans()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResp, err := json.Marshal(allPlans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func executePlanHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, _ := strconv.Atoi(params.ByName("ID"))

	plan, _ := db.GetPlan(uint(id))
	pair, _ := db.GetPair(plan.PairID)
	orders, _ := db.GetOrders(plan.ID)
	balance, _ := db.GetCurrentBalance(pair.QuoteCurrency)

	plan.FinalizeOrders(balance.Available, pair, orders)
    // TODO: this should be in an transaction so that when sending to exchange fails the database gets rolledback

	err := e.PlaceOrders(plan, pair, orders)
	if err != nil {
	}
    plan.Status = cryptodb.StatusOrdered
	db.SaveOrders(orders)
    // TODO: if no errors occured -> Commit
}

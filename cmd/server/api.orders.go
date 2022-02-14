package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getOrdersHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	planID, err := strconv.Atoi(p.ByName("PlanID"))
	orders, err := db.GetOrders(uint(planID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	jsonResp, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

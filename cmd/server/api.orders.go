package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

func getOrdersHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	planID, err := strconv.Atoi(p.ByName("PlanID"))
    var orders []cryptodb.Order

	result := db.Where("plan_id = ?", planID).Find(&orders)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}
	jsonResp, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

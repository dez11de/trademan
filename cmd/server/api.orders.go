package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
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

func saveOrdersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var orders []cryptodb.Order
	err := json.NewDecoder(r.Body).Decode(&orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	var result *gorm.DB
	if orders[cryptodb.MarketStopLoss].ID == 0 {
		result = db.Create(&orders)
	} else {
		result = db.Save(&orders)
	}
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func getPlansHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func getPlanHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("ID"))
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var plan cryptodb.Plan
	result := db.Where("id = ?", id).First(&plan)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}
	jsonResp, err := json.Marshal(plan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func savePlanHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var plan cryptodb.Plan
	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

    log.Printf("saving plan: %v", plan)

	var result *gorm.DB
	if plan.ID == 0 {
		result = db.Create(&plan)
	} else {
		result = db.Save(&plan)
	}
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}

	jsonResp, err := json.Marshal(plan)
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

	var plan cryptodb.Plan
	var pair cryptodb.Pair
	var orders []cryptodb.Order
	var balance cryptodb.Balance
	db.Where("id = ?", id).First(&plan)
	db.Where("id = ?", plan.PairID).First(&pair)
	db.Where("plan_id = ?", plan.ID).Find(&orders)
	db.Where("symbol = ?", pair.QuoteCurrency).Order("created_at DESC").First(&balance)

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

	result = tx.Create(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: "Sending plan."})
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		tx.Rollback()
		return
	}
	tx.Commit()

	pause <- struct{}{}
	err = placeOrders(plan, pair, orders)
	play <- struct{}{}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		db.Create(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: fmt.Sprintf("Exchange did not accept plan. %s", err)})
		return
	}
	db.Create(&cryptodb.Log{PlanID: plan.ID, Source: cryptodb.Server, Text: "Sending plan successful."})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

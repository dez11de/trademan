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

func getReviewHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("PlanID"))
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var Review cryptodb.Review
	result := db.Where("plan_id = ?", id).First(&Review)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}
	jsonResp, err := json.Marshal(Review)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func saveReviewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var review cryptodb.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	var result *gorm.DB
	if review.ID == 0 {
		result = db.Create(&review)
	} else {
		result = db.Save(&review)
	}
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}

	jsonResp, err := json.Marshal(review)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func getReviewOptionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	colNames := []string{"risk", "timing", "stop_loss", "entry", "emotion", "follow_plan", "order_management", "move_stop_loss_in_profit", "take_profit_strategy", "take_profit_count"}
	var colOptions []string
	options := make(map[string][]string)

	for _, colName := range colNames {
		result := db.Table("reviews").Distinct(colName).Find(&colOptions)
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, result.Error)
			return
		}
		options[colName] = colOptions
	}

	jsonResp, err := json.Marshal(options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

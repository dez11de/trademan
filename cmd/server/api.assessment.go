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

func getAssessmentHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("PlanID"))
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var assessment cryptodb.Assessment
	result := db.Where("plan_id = ?", id).First(&assessment)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}
	jsonResp, err := json.Marshal(assessment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func saveAssessmentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var assessment cryptodb.Assessment
	err := json.NewDecoder(r.Body).Decode(&assessment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	var result *gorm.DB
	if assessment.ID == 0 {
		result = db.Create(&assessment)
	} else {
		result = db.Save(&assessment)
	}
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}

	if assessment.Status == "Completed" {
		var plan cryptodb.Plan
		db.Where("id = ?", assessment.PlanID).First(&plan)
		plan.Status = cryptodb.Logged
		db.Save(&plan)
	}

	jsonResp, err := json.Marshal(assessment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func getAssessmentOptionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	colNames := []string{"risk", "timing", "stop_loss", "entry", "emotion", "follow_plan", "order_management", "move_stop_loss_in_profit", "take_profit_strategy", "take_profit_count"}
	var colOptions []string
	options := make(map[string][]string)

	for _, colName := range colNames {
		result := db.Table("assessments").Distinct(colName).Find(&colOptions)
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

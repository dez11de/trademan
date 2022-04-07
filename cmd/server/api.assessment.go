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

func getAssessmentHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("ID"))
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

	log.Printf("Saving assesment %+v", assessment)

	var result *gorm.DB
	if assessment.ID == 0 {
		result = db.Create(&assessment)
	} else {
		result = db.Save(&assessment)
	}
	if result.Error != nil {
		log.Printf("Error occured here %+v", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, result.Error)
		return
	}

	if assessment.Status == "Completed" {
        log.Print("Archiving plan")
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

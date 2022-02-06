package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO: should only return 'open' plans. See bybit api documentation.
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

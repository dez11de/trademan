package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// TODO: should only return Active pairs. See GORM api documentation.
func allPairsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allPairs, err := db.GetPairs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResp, err := json.Marshal(allPairs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func pairHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pairID, err := strconv.Atoi(p.ByName("PairID"))
	pair, err := db.GetPair(uint(pairID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResp, err := json.Marshal(pair)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func searchPairsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	part := p.ByName("part")
	pair, err := db.FindPairNames(part)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResp, err := json.Marshal(pair)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

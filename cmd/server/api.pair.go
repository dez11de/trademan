package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO: rewrite as gRPC
func allPairsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allPairs, err := db.GetPairs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
	}
	jsonResp, err := json.Marshal(allPairs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

// func pairHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	pairID, err := strconv.Atoi(p.ByName("PairID"))
// 	pair, err := db.GetPair(uint(pairID))
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
//         w.Write([]byte(err.Error()))
//         return
// 	}
// 	jsonResp, err := json.Marshal(pair)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
//         w.Write([]byte(err.Error()))
//         return
// 	}
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.Write(jsonResp)
// }

package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func dbNameHandler (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    dbName := trademanCfg.Database.Database

	jsonResp, err := json.Marshal(dbName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

package main

import (
	"encoding/json"
	"net/http"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

func setupHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setup := cryptodb.NewSetup()

	err := json.NewDecoder(r.Body).Decode(&setup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if setup.Plan.ID == 0 {
		err = db.CreateSetup(&setup)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		err = db.SaveSetup(cryptodb.User, &setup)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	jsonResp, err := json.Marshal(setup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

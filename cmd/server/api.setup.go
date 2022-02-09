package main

import (
	"encoding/json"
	"net/http"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

// TODO: error handling
func setupHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setup := cryptodb.NewSetup()

	err := json.NewDecoder(r.Body).Decode(&setup)
	if err != nil {
	}

	if setup.Plan.ID == 0 {
		db.CreateSetup(&setup)
	} else {
		db.SaveSetup(cryptodb.SourceUser, &setup)
	}

	// TODO: return statusOK and updated Setup or something
}


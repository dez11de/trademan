package cryptodb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (db *api) performanceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var pReq Performance
	err := json.NewDecoder(r.Body).Decode(&pReq)
	if err != nil {
		log.Printf("error decoding body: %v", err)
	}

	performance, err := db.GetPerformance(pReq.Symbol, pReq.Since)
	if err != nil {
		log.Printf("error getting performance from database: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        // TODO: shouldn't i write something else?
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
        jsonResp, err := json.Marshal(Performance{Performance: performance})
        if err != nil {
            log.Printf("Error marshalling PerformanceResponse %v", err)
        }
        w.Write(jsonResp)
	}
}

func (db *api) planHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Returning plans")
	log.Println("Endpoint hit: plans")
}

func (db *api) HandleRequests() {
	log.Printf("Routing HTTP handler functions")
	http.HandleFunc("/performance", db.performanceHandler)
	http.HandleFunc("/plans", db.planHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

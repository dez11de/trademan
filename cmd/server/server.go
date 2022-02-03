package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
)

/*
func performanceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
    var symbolPerformance string
	err := json.NewDecoder(r.Body).Decode(&symbolPerformance)
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
*/

func  planHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: plan")
	switch r.Method {
	case "GET":
		log.Print("GET request received, returning plan")
		PlanIDs, ok := r.URL.Query()["PlanID"]
		if !ok {
			// TODO: Rewrite this so that it also can return all-Plans, all-Openplans(default?), all-LoggedPlans
			log.Print("no PlanID found, don't know what to do, returning nothing.")
			// TODO set http error
			return
		}
		PlanID, err := strconv.Atoi(PlanIDs[0])
		PlanID64 := uint(PlanID)
		if err != nil {
			log.Printf("Invalid ID? %s", err)
			// TODO: set http error
			return
		}
		plan, err := db.GetPlan(PlanID64)
		if err != nil {
			log.Printf("Error getting plan %s", err)
		}
		log.Printf("Plan found! %v", plan)
		// TODO: write plan to w

	case "POST":
		log.Print("POST request received, add/updating plan")
	}
}

func pairsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Maybe use Gorilla for this kind of stuff?
	var err error
	var pairs map[uint]cryptodb.Pair
	var pairNames []string
	switch r.Method {
	case "GET":
		search, ok := r.URL.Query()["search"]
		if ok {
			log.Printf("Searching for pair like %s", search[0])
			pairNames, err = db.FindPairNames(search[0])
			if err != nil {
				log.Printf("Error searching pairsNames: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			jsonResp, err := json.Marshal(pairNames)
			if err != nil {
				log.Printf("Error marshalling Pairs %v", err)
			}
			w.Write(jsonResp)
		} else {
			_, ok := r.URL.Query()["PairID"]
			if !ok {
				log.Print("No PairID specified, returning all Pairs")
				pairs, err = db.GetPairs()
				if err != nil {
					log.Printf("Error getting pairs %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				jsonResp, err := json.Marshal(pairs)
				if err != nil {
					log.Printf("Error marshalling Pairs %v", err)
				}
				w.Write(jsonResp)
			}
		}
		if err != nil {
			log.Printf("error getting pairs from database: %v", err)
			// TODO: shouldn't i write something else?
		}

	case "POST":
		// TODO: this shouldn't be even supported?
		log.Print("POST request received, add/updating pair?")
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: list")
	plans, err := db.GetPlans()
	if err != nil {
		log.Printf("error getting plans from database %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(plans)
	if err != nil {
		log.Printf("Error marshalling Pairs %v", err)
	}
	w.Write(jsonResp)
}


func HandleRequests() {
	log.Printf("Routing HTTP handler functions")
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/plan", planHandler)
	http.HandleFunc("/pairs", pairsHandler)

	// TODO make at least port configurable
	log.Fatal(http.ListenAndServe(":8888", nil))
}

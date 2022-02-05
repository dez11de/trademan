package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/julienschmidt/httprouter"
)

// TODO: should only return Active pairs. See bybit api documentation.
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

func getOrdersHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	planID, err := strconv.Atoi(p.ByName("PlanID"))
	orders, err := db.GetOrders(uint(planID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResp, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

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

	// TODO: return statusOK or something
}

const APIv1Base = "/api/v1/"

func HandleRequests() {
	log.Printf("Routing HTTP handler functions")
	router := httprouter.New()

	router.GET(APIv1Base+"pairs", allPairsHandler)
	router.GET(APIv1Base+"pair/:PairID", pairHandler)
	router.GET(APIv1Base+"pairs_search/:part", searchPairsHandler)

	router.GET(APIv1Base+"plans", allPlansHandler)

	router.GET(APIv1Base+"orders/:PlanID", getOrdersHandler)

	router.POST(APIv1Base+"setup", setupHandler)

	// TODO make at least port configurable
	log.Fatal(http.ListenAndServe(":8888", router))
}

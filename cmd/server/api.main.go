package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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

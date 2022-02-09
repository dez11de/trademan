package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const APIv1Base = "/api/v1/"

func HandleRequests(c RESTServerConfig) {
	log.Printf("HTTP Server ready.\n")
	router := httprouter.New()

	router.GET(APIv1Base+"pairs", allPairsHandler)
	router.GET(APIv1Base+"pair/:PairID", pairHandler)
	router.GET(APIv1Base+"pairs_search/:part", searchPairsHandler)

	router.GET(APIv1Base+"plans", allPlansHandler)
	router.GET(APIv1Base+"plan/execute/:ID", executePlanHandler)

	router.GET(APIv1Base+"orders/:PlanID", getOrdersHandler)

	router.POST(APIv1Base+"setup", setupHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", c.Host, c.Port), router))
}
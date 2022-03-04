package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const APIv1Base = "/api/v1/"

func HandleRequests(c RESTServerConfig) {
	router := httprouter.New()

	router.GET(APIv1Base+"pairs", allPairsHandler)

	router.GET(APIv1Base+"plans", allPlansHandler)
	router.GET(APIv1Base+"plan/execute/:ID", executePlanHandler)

	router.GET(APIv1Base+"orders/:PlanID", getOrdersHandler)

	router.POST(APIv1Base+"setup", setupHandler)

	log.Printf("==========================[ API Server Ready ]==========================\n")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", c.Host, c.Port), router))
}

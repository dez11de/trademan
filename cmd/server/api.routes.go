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

	router.GET(APIv1Base+"databaseName", dbNameHandler)

	router.GET(APIv1Base+"pairs", allPairsHandler)

	router.GET(APIv1Base+"execute/:ID", executePlanHandler)
    router.GET(APIv1Base+"plans/:Archived", getPlansHandler)
	router.GET(APIv1Base+"plan/:ID", getPlanHandler)
	router.POST(APIv1Base+"plan", savePlanHandler)

	router.GET(APIv1Base+"orders/:PlanID", getOrdersHandler)
	router.POST(APIv1Base+"orders", saveOrdersHandler)

	router.GET(APIv1Base+"review/:PlanID", getReviewHandler)
	router.POST(APIv1Base+"review", saveReviewHandler)

	router.GET(APIv1Base+"reviewOptions", getReviewOptionsHandler)

	router.GET(APIv1Base+"logs/:PlanID", getLogsHandler)

	log.Println("=========================[ API Server Ready ]=========================")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", c.Host, c.Port), router))
}

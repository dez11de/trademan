package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
)

// TODO: make host:port configurable in env/param/file
const BaseURL = "http://localhost:8888/api/v1/"

func getPairs() (pairs []cryptodb.Pair, err error) {
	resp, err := http.Get(BaseURL + "pairs")
	if err != nil {
		return pairs, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &pairs)

	return pairs, err
}

func searchPairs(s string) (pairs []string, err error) {
	resp, err := http.Get(BaseURL + "pairs_search/" + s)
	if err != nil {
		return pairs, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &pairs)

	return pairs, err
}

func getPair(id uint) (pair cryptodb.Pair, err error) {
	resp, err := http.Get(BaseURL + "pair/" + strconv.Itoa(int(id)))
	if err != nil {
		return pair, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &pair)

	return pair, err
}

func getPlans() (plans []cryptodb.Plan, err error) {
	resp, err := http.Get(BaseURL + "plans")
	if err != nil {
		return plans, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &plans)

	return plans, err
}

func getOrders(PlanID uint) (orders []cryptodb.Order, err error) {
	orders = cryptodb.NewOrders(PlanID) // TODO: is this really necessary here?
	resp, err := http.Get(BaseURL + "orders/" + strconv.Itoa(int(PlanID)))
	if err != nil {
		return orders, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &orders)

	return orders, err
}

func sendSetup(p cryptodb.Plan, o []cryptodb.Order) (err error) {
	var setup cryptodb.Setup
	setup.Plan = p
	setup.Orders = o

	setupJSON, _ := json.Marshal(setup)
    resp, err := http.Post(BaseURL+"setup", "", bytes.NewBuffer(setupJSON)) // TODO: should this include "application/json or Content-Type/json"
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

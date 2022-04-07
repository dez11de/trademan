package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
)

func getDBName() (name string, err error) {
	resp, err := http.Get(BaseURL + "databaseName")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &name)
	return name, err
}

func getPairs() (pairs []cryptodb.Pair, err error) {
	resp, err := http.Get(BaseURL + "pairs")
	if err != nil {
		return pairs, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return pairs, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pairs, err
	}

	err = json.Unmarshal(body, &pairs)
	return pairs, err
}

func getPair(id uint) (pair cryptodb.Pair, err error) {
	resp, err := http.Get(BaseURL + "pair/" + strconv.Itoa(int(id)))
	if err != nil {
		return pair, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return pair, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pair, err
	}

	err = json.Unmarshal(body, &pair)
	return pair, err
}

func getPlans() (plans []cryptodb.Plan, err error) {
	resp, err := http.Get(BaseURL + "plans")
	if err != nil {
		return plans, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return plans, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return plans, err
	}

	err = json.Unmarshal(body, &plans)
	return plans, err
}

func getPlan(id uint64) (plan cryptodb.Plan, err error) {
	resp, err := http.Get(BaseURL + "plan/" +strconv.Itoa(int(id)))
	if err != nil {
		return cryptodb.Plan{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return cryptodb.Plan{}, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cryptodb.Plan{}, err
	}

	err = json.Unmarshal(body, &plan)
	return plan, err
}

func savePlan(p cryptodb.Plan) (plan cryptodb.Plan, err error) {
	setupJSON, _ := json.Marshal(p)
	resp, err := http.Post(BaseURL+"plan", "application/json", bytes.NewBuffer(setupJSON))
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return p, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(body, &plan)
	return plan, err
}


func executePlan(id uint64) (plan cryptodb.Plan, err error) {
	resp, err := http.Get(BaseURL + "execute/" + strconv.Itoa(int(id)))
	if err != nil {
		return plan, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return plan, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return plan, err
	}

	err = json.Unmarshal(body, &plan)
	return plan, err
}

func getOrders(PlanID uint64) (orders []cryptodb.Order, err error) {
	resp, err := http.Get(BaseURL + "orders/" + strconv.Itoa(int(PlanID)))
	if err != nil {
		return orders, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return orders, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return orders, err
	}

	err = json.Unmarshal(body, &orders)
	return orders, err
}

func saveOrders(o []cryptodb.Order) (orders []cryptodb.Order, err error) {
	setupJSON, _ := json.Marshal(o)
	resp, err := http.Post(BaseURL+"orders", "application/json", bytes.NewBuffer(setupJSON))
	if err != nil {
		return o, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return o, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}

	err = json.Unmarshal(body, &orders)
	return orders, err
}


func getLogs(PlanID uint64) (entries []cryptodb.Log, err error) {
	resp, err := http.Get(BaseURL + "logs/" + strconv.Itoa(int(PlanID)))
	if err != nil {
		return entries, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return entries, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entries, err
	}

	err = json.Unmarshal(body, &entries)
	return entries, err
}

func getAssessment(PlanID uint64) (assessment cryptodb.Assessment , err error) {
	resp, err := http.Get(BaseURL + "assessment/" + strconv.Itoa(int(PlanID)))
	if err != nil {
		return assessment, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return assessment, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return assessment, err
	}

	err = json.Unmarshal(body, &assessment)
	return assessment, err
}

func saveAssessment(a cryptodb.Assessment) (assessment cryptodb.Assessment, err error) {
	assessmentJSON, _ := json.Marshal(a)
	resp, err := http.Post(BaseURL+"assessment", "application/json", bytes.NewBuffer(assessmentJSON))
	if err != nil {
		return a, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return a, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(body, &assessment)
	return assessment, err
}

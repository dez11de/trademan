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
    resp, err := http.Get(BaseURL + "plans/" + strconv.FormatBool(ui.showArchived))
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
	resp, err := http.Get(BaseURL + "plan/" + strconv.Itoa(int(id)))
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

func getReview(PlanID uint64) (review cryptodb.Review, err error) {
	resp, err := http.Get(BaseURL + "review/" + strconv.Itoa(int(PlanID)))
	if err != nil {
		return review, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return review, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return review, err
	}

	err = json.Unmarshal(body, &review)
	return review, err
}

func saveReview(r cryptodb.Review) (review cryptodb.Review, err error) {
	reviewJSON, _ := json.Marshal(r)
	resp, err := http.Post(BaseURL+"review", "application/json", bytes.NewBuffer(reviewJSON))
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return r, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &review)
	return review, err
}

func getReviewOptions() (options map[string][]string, err error) {
	resp, err := http.Get(BaseURL + "reviewOptions")
	if err != nil {
		return options, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(resp.Body)
		return options, errors.New(string(errorMessage))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return options, err
	}

	err = json.Unmarshal(body, &options)
	return options, err
}

package exchange

// Copied from  https://github.com/frankrap/bybit-api/
// All credits to him.

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type RequestParameters map[string]interface{}

func (e *Exchange) PublicRequest(method string, apiURL string, params map[string]interface{}, result interface{}) (resp []byte, err error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%v=%v", k, params[k]))
	}

	param := strings.Join(p, "&")
	fullURL := e.restHost + apiURL
	if param != "" {
		fullURL += "?" + param
	}
	if e.debugMode {
		log.Printf("PublicRequest: %v", fullURL)
	}
	var binBody = bytes.NewReader(make([]byte, 0))

	var request *http.Request
	request, err = http.NewRequest(method, fullURL, binBody)
	if err != nil {
		return
	}

	var response *http.Response
	response, err = e.restClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if e.debugMode {
		log.Printf("PublicRequest: %v", string(resp))
	}

	err = json.Unmarshal(resp, result)
	return
}

func (b *Exchange) PrivateRequest(method string, apiURL string, params map[string]interface{}, result interface{}) (resp []byte, err error) {
	timestamp := time.Now().UnixNano() / 1e6

	params["api_key"] = b.apiKey
	params["timestamp"] = timestamp

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%v=%v", k, params[k]))
	}

	param := strings.Join(p, "&")
	signature := b.getSigned(param)
	param += "&sign=" + signature

	fullURL := b.restHost + apiURL + "?" + param
	if b.debugMode {
		log.Printf("SignedRequest: %v", fullURL)
	}
	var binBody = bytes.NewReader(make([]byte, 0))

	// get a http request
	var request *http.Request
	request, err = http.NewRequest(method, fullURL, binBody)
	if err != nil {
		return
	}

	var response *http.Response
	response, err = b.restClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if b.debugMode {
		log.Printf("SignedRequest: %v", string(resp))
	}

	err = json.Unmarshal(resp, result)
	return
}

func (b *Exchange) SignedRequest(method string, apiURL string, params map[string]interface{}, result interface{}) (fullURL string, resp []byte, err error) {
	timestamp := time.Now().UnixNano() / 1e6

	params["api_key"] = b.apiKey
	params["timestamp"] = timestamp

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%v=%v", k, params[k]))
	}

	param := strings.Join(p, "&")
	signature := b.getSigned(param)
	param += "&sign=" + signature

	fullURL = b.restHost + apiURL + "?" + param
	if b.debugMode {
		log.Printf("SignedRequest: %v", fullURL)
	}
	var binBody = bytes.NewReader(make([]byte, 0))

	// get a http request
	var request *http.Request
	request, err = http.NewRequest(method, fullURL, binBody)
	if err != nil {
		return
	}

	var response *http.Response
	response, err = b.restClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if b.debugMode {
		log.Printf("SignedRequest: %v", string(resp))
	}
	err = json.Unmarshal(resp, result)
	return
}

func (b *Exchange) getSigned(param string) string {
	sig := hmac.New(sha256.New, []byte(b.apiSecret))
	sig.Write([]byte(param))
	signature := hex.EncodeToString(sig.Sum(nil))
	return signature
}

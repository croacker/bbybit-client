package main

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
	"net/http/httputil"
	"runtime"
	"strconv"
	"time"

	"github.com/croacker/bybit-client/internal/client"
	"github.com/croacker/bybit-client/internal/config"
)

const ONE_MINUTE = 960000

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

var url string = "https://api-testnet.bybit.com"
var api_key = ""
var apiSecret = ""
var recv_window = "5000"
var signature = ""

func main() {
	appConfig := config.LoadConfig()
	endMs := getEndMilis()
	startMs := endMs - ONE_MINUTE
	for _, symbol := range appConfig.Symbols {
		go client.MarkPriceKline(symbol, startMs, endMs)
	}

	runtime.Goexit()
}

func getEndMilis() int64 {
	now := time.Now()
	return now.UnixMilli()
}

func requestSymbols() {
	now := time.Now()
	client := httpClient()
	endpoint := "/v5/market/instruments-info"
	params := "category=linear"
	request, error := http.NewRequest("GET", url+endpoint+"?"+params, nil)
	if error != nil {
		panic(error)
	}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	elapsed := time.Since(now).Seconds()
	fmt.Printf("\n%s took %v seconds \n", endpoint, elapsed)
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func getRequest(client *http.Client, method string, params string, endPoint string) []byte {
	now := time.Now()
	unixNano := now.UnixNano()
	time_stamp := unixNano / 1000000
	hmac256 := hmac.New(sha256.New, []byte(apiSecret))
	hmac256.Write([]byte(strconv.FormatInt(time_stamp, 10) + api_key + recv_window + params))
	signature = hex.EncodeToString(hmac256.Sum(nil))
	request, error := http.NewRequest("GET", url+endPoint+"?"+params, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BAPI-API-KEY", api_key)
	request.Header.Set("X-BAPI-SIGN", signature)
	request.Header.Set("X-BAPI-TIMESTAMP", strconv.FormatInt(time_stamp, 10))
	request.Header.Set("X-BAPI-SIGN-TYPE", "2")
	request.Header.Set("X-BAPI-RECV-WINDOW", recv_window)
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Request Dump:\n%s", string(reqDump))
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	elapsed := time.Since(now).Seconds()
	fmt.Printf("\n%s took %v seconds \n", endPoint, elapsed)
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return body
}

func postRequest(client *http.Client, method string, params interface{}, endPoint string) []byte {
	now := time.Now()
	unixNano := now.UnixNano()
	time_stamp := unixNano / 1000000
	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	hmac256 := hmac.New(sha256.New, []byte(apiSecret))
	hmac256.Write([]byte(strconv.FormatInt(time_stamp, 10) + api_key + recv_window + string(jsonData[:])))
	signature = hex.EncodeToString(hmac256.Sum(nil))
	request, error := http.NewRequest("POST", url+endPoint, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BAPI-API-KEY", api_key)
	request.Header.Set("X-BAPI-SIGN", signature)
	request.Header.Set("X-BAPI-TIMESTAMP", strconv.FormatInt(time_stamp, 10))
	request.Header.Set("X-BAPI-SIGN-TYPE", "2")
	request.Header.Set("X-BAPI-RECV-WINDOW", recv_window)
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Request Dump:\n%s", string(reqDump))
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	elapsed := time.Since(now).Seconds()
	fmt.Printf("\n%s took %v seconds \n", endPoint, elapsed)
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return body
}

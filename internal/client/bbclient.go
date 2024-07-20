package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/croacker/bybit-client/internal/dto"
)

const BB_URL string = "https://api-testnet.bybit.com"
const INTERVAL = 5

type MarkPriceKlineClient struct {
	Symbol string
}

func MarkPriceKline(symbol string, start int64, end int64) {
	client := httpClient()

	for {
		url := toUrl(symbol, start, end)
		log.Println("request mark-price-kline for", symbol)

		request, error := http.NewRequest("GET", url, nil)
		if error != nil {
			panic(error)
		}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}

		body, _ := ioutil.ReadAll(response.Body)
		dto := dto.MarkPriceKlineResponseDto{}
		err := json.Unmarshal(body, &dto)
		if err != nil {
			log.Fatal("error unmarshal MarkPriceKline:", err)
		}
		log.Println("response body:", dto)

		response.Body.Close()

		time.Sleep(INTERVAL * time.Second)
		start += INTERVAL // TODO
		end += INTERVAL   // TODO
	}
}

func toUrl(symbol string, start int64, end int64) string {
	url := BB_URL + "/v5/market/mark-price-kline?"
	url += "category=linear"
	url += "&symbol=" + symbol
	url += "&interval=15"
	url += "&start=" + strconv.FormatInt(start, 10)
	url += "&end=" + strconv.FormatInt(end, 10)
	return url
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

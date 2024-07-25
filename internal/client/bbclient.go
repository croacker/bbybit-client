package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/croacker/bybit-client/internal/dto"
	"github.com/croacker/bybit-client/internal/store"
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

		body, _ := io.ReadAll(response.Body)
		responseDto := dto.MarkPriceKlineResponseDto{}
		err := json.Unmarshal(body, &responseDto)
		if err != nil {
			log.Fatal("error unmarshal mark-price-kline response:", err)
		}

		for _, item := range responseDto.Result.List {
			candle := dto.NewMarkPriceKlineCandleDto(responseDto.Result.Symbol, item[0], item[1], item[2], item[3], item[4])
			log.Println("candle:", candle)
			processCandle(candle)
		}

		response.Body.Close()

		time.Sleep(INTERVAL * time.Second)
		start += INTERVAL // TODO
		end += INTERVAL   // TODO
	}
}

func processCandle(candle *dto.MarkPriceKlineCandleDto) {
	storedItem := store.GetStoreItem(candle.Symbol)
	log.Println("stored item:", storedItem)
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

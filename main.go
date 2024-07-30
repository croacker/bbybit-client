package main

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/croacker/bybit-client/internal/client"
	"github.com/croacker/bybit-client/internal/config"
	"github.com/croacker/bybit-client/internal/dto"
)

const ONE_MINUTE = 960000

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func main() {
	appConfig := config.LoadConfig()

	bbClient := client.NewBbClient(appConfig)
	candlesCh := bbClient.GetOutgoingChannel()
	bbClient.Start()

	go readCandles(candlesCh)

	runtime.Goexit()
}

func readCandles(candlesCh chan *dto.MarkPriceKlineCandleDto) {
	for candle := range candlesCh {
		log.Println("receive candle: ", candle)
	}
}

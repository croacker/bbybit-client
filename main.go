package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/croacker/bybit-client/internal/client"
	"github.com/croacker/bybit-client/internal/config"
	"github.com/croacker/bybit-client/internal/db"
	"github.com/croacker/bybit-client/internal/dto"
)

const ONE_MINUTE = 960000

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func main() {
	appConfig := config.LoadConfig()

	db.SetupDb(appConfig)
	db.AllChats()

	bbClient := client.NewBbClient(appConfig)
	candlesCh := bbClient.GetOutgoingChannel()
	bbClient.Start()

	tgClient := client.NewTgClient(appConfig)
	tgOutgoingCh := tgClient.GetOutgoingChannel()
	tgClient.Start()

	go readCandles(candlesCh, tgOutgoingCh)

	go writeMessages(tgOutgoingCh)

	runtime.Goexit()
}

func readCandles(candlesCh chan *dto.MarkPriceKlineCandleDto, tgOutgoingCh chan string) {
	for candle := range candlesCh {
		msg := fmt.Sprintf("%v", candle)
		tgOutgoingCh <- msg
		log.Println("receive candle: ", candle)
	}
}

func writeMessages(tgOutgoingCh chan string) {

}

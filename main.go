package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/croacker/bybit-client/internal/client"
	"github.com/croacker/bybit-client/internal/config"
	"github.com/croacker/bybit-client/internal/db"
	"github.com/croacker/bybit-client/internal/dto"
	"github.com/croacker/bybit-client/internal/service"
)

const ONE_MINUTE = 960000

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func main() {
	appConfig := config.LoadConfig()

	db.SetupDb(appConfig)

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
		if service.NeedSendAlert(candle) {
			msg := fmt.Sprintf("%v", candle)
			tgOutgoingCh <- msg
		}
	}
}

func writeMessages(tgOutgoingCh chan string) {

}

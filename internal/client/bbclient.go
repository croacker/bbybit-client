package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/croacker/bybit-client/internal/config"
	"github.com/croacker/bybit-client/internal/dto"
)

const BB_URL string = "https://api-testnet.bybit.com"
const INTERVAL = 15
const ONE_MINUTE = 960000

var bbClient *MarkPriceKlineClient

type MarkPriceKlineClient struct {
	Symbol       string
	symbols      []string
	bbUrl        string
	interval     int32
	outgoingChan chan *dto.MarkPriceKlineCandleDto // TODO incoming????
}

type ReqDetails struct {
	symbol string
	start  int64
	end    int64
}

func NewBbClient(cfg *config.AppConfig) *MarkPriceKlineClient {
	bbClient = &MarkPriceKlineClient{
		Symbol:       "",
		symbols:      cfg.Symbols,
		bbUrl:        cfg.BbClient.Url,
		interval:     cfg.BbClient.Interval,
		outgoingChan: make(chan *dto.MarkPriceKlineCandleDto),
	}
	return bbClient
}

func (c *MarkPriceKlineClient) Start() {
	log.Println("start bb-client...")

	go loopRequests(c.symbols, c.interval)
	log.Println("bb-client started")
}

func (c *MarkPriceKlineClient) GetOutgoingChannel() chan *dto.MarkPriceKlineCandleDto {
	return c.outgoingChan
}

func loopRequests(symbols []string, interval int32) {
	end := getEndMilis()
	start := end - ONE_MINUTE
	for {
		for _, symbol := range symbols {
			body := requestMarkPriceKline(symbol, start, end)
			candles := getCandles(body)
			for _, candle := range candles {
				bbClient.outgoingChan <- candle
			}
		}
		end += INTERVAL
		start += INTERVAL
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func requestMarkPriceKline(symbol string, start int64, end int64) []byte {
	client := httpClient() // TODO
	url := toUrl(symbol, start, end)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func getCandles(body []byte) []*dto.MarkPriceKlineCandleDto {
	var candles []*dto.MarkPriceKlineCandleDto
	responseDto := unmarshalBody(body)
	for _, item := range responseDto.Result.List {
		candle := dto.NewMarkPriceKlineCandleDto(responseDto.Result.Symbol, item[0], item[1], item[2], item[3], item[4])
		candles = append(candles, candle)
	}
	return candles
}

func unmarshalBody(body []byte) *dto.MarkPriceKlineResponseDto {
	responseDto := dto.MarkPriceKlineResponseDto{}
	err := json.Unmarshal(body, &responseDto)
	if err != nil {
		log.Fatal("error unmarshal mark-price-kline response:", err)
	}
	return &responseDto
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

func getEndMilis() int64 {
	now := time.Now()
	return now.UnixMilli()
}

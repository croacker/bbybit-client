package dto

import (
	"log"
	"strconv"
)

type MarkPriceKlineResponseDto struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Symbol   string     `json:"symbol"`
		Category string     `json:"category"`
		List     [][]string `json:"list"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type MarkPriceKlineCandleDto struct {
	Symbol     string
	StartTime  int64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
}

func NewMarkPriceKlineCandleDto(symbol string, startTime string, openPrice string, highPrice string, lowPrice string, closePrice string) *MarkPriceKlineCandleDto {
	startTimeV, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		log.Println("ERROR:", err)
	}

	openPriceV, err := strconv.ParseFloat(openPrice, 64)
	if err != nil {
		log.Println("ERROR:", err)
	}

	highPriceV, err := strconv.ParseFloat(highPrice, 64)
	if err != nil {
		log.Println("ERROR:", err)
	}

	lowPriceV, err := strconv.ParseFloat(lowPrice, 64)
	if err != nil {
		log.Println("ERROR:", err)
	}

	closePriceV, err := strconv.ParseFloat(closePrice, 64)
	if err != nil {
		log.Println("ERROR:", err)
	}

	return &MarkPriceKlineCandleDto{
		Symbol:     symbol,
		StartTime:  startTimeV,
		OpenPrice:  openPriceV,
		HighPrice:  highPriceV,
		LowPrice:   lowPriceV,
		ClosePrice: closePriceV,
	}

}

package processor

import (
	"log"
	"math"

	"github.com/croacker/bybit-client/internal/dto"
	"github.com/croacker/bybit-client/internal/store"
)

func NeedSendAlert(candle *dto.MarkPriceKlineCandleDto) bool {
	result := false
	storedItem := store.GetStoredItem(candle.Symbol)
	log.Println("stored item:", storedItem)
	if storedItem.StartTime != 0 {
		if isAlert(storedItem.OpenPrice, candle.OpenPrice) {
			log.Println("!!!! new price:", candle.OpenPrice)
			result = true
		}
		if isAlert(storedItem.HighPrice, candle.HighPrice) {
			log.Println("!!!! new price:", candle.HighPrice)
			result = true
		}
		if isAlert(storedItem.LowPrice, candle.LowPrice) {
			log.Println("!!!! new price:", candle.LowPrice)
			result = true
		}
		if isAlert(storedItem.ClosePrice, candle.ClosePrice) {
			log.Println("!!!! new price:", candle.ClosePrice)
			result = true
		}
	}
	storedItem.StartTime = candle.StartTime
	storedItem.OpenPrice = candle.OpenPrice
	storedItem.HighPrice = candle.HighPrice
	storedItem.LowPrice = candle.LowPrice
	storedItem.ClosePrice = candle.ClosePrice
	store.StoreItem(storedItem)
	return result
}

func isAlert(oldPrice float64, newPrice float64) bool {
	percents15 := oldPrice * 0.15
	delta := math.Abs(oldPrice - newPrice)
	return delta > percents15
}

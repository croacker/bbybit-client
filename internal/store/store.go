package store

import (
	"fmt"
)

type CandleStore struct {
	Items map[string]CandleStoreItem
}

type CandleStoreItem struct {
	Symbol     string
	StartTime  int64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
}

var store *CandleStore

func (c CandleStoreItem) String() string {
	return fmt.Sprintf("symbol: %v, startTime:%v, openPrice:%v, highPrice:%v, lowPrice:%v, closePrice:%v", c.Symbol, c.StartTime, c.OpenPrice, c.HighPrice, c.LowPrice, c.ClosePrice)
}

func GetStore() *CandleStore {
	if store == nil {
		store = &CandleStore{}
		store.Items = make(map[string]CandleStoreItem)
	}
	return store
}

func GetStoredItem(symbol string) *CandleStoreItem {
	s := GetStore()
	item, found := s.Items[symbol]
	if !found {
		item = CandleStoreItem{Symbol: symbol}
		s.Items[symbol] = item
	}
	return &item
}

func StoreItem(item *CandleStoreItem) {
	s := GetStore()
	s.Items[item.Symbol] = *item
}

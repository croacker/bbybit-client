package store

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

func GetStore() *CandleStore {
	if store == nil {
		store = &CandleStore{}
		store.Items = make(map[string]CandleStoreItem)
	}
	return store
}

func GetStoreItem(symbol string) *CandleStoreItem {
	s := GetStore()
	item, found := s.Items[symbol]
	if !found {
		item = CandleStoreItem{Symbol: symbol}
		s.Items[symbol] = item
	}
	return &item
}

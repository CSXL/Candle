// Contains wrappers for the FinHub API.
package api

import (
	"context"
	"fmt"
	"log"

	finhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	websocket "github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type FinHubClient struct {
	*finhub.DefaultApiService
	apiKey string
}

type Quote struct {
	Symbol        string
	Price         float32
	Change        float32
	PercentChange float32
	DayHigh       float32
	DayLow        float32
	PreviousOpen  float32
	PreviousClose float32
}

type WebsocketResponse struct {
	Type string  `json:"type"`
	Data []Trade `json:"data"`
}

type Trade struct {
	Symbol    string  `json:"s"`
	Price     float32 `json:"p"`
	Volume    int     `json:"v"`
	Timestamp int64   `json:"t"`
}

type Candles struct {
	Open      []float32
	High      []float32
	Low       []float32
	Close     []float32
	Volume    []float32
	Timestamp []int64
}

func NewFinHubClient(apiKey string) *FinHubClient {
	config := finhub.NewConfiguration()
	config.AddDefaultHeader("X-Finnhub-Token", apiKey)
	DefaultApiService := finhub.NewAPIClient(config).DefaultApi
	return &FinHubClient{
		DefaultApiService: DefaultApiService,
		apiKey:            apiKey,
	}
}

func (c *FinHubClient) GetCurrentPrice(symbol string) (Quote, error) {
	// Documentations: https://finnhub.io/docs/api/quote
	data, _, err := c.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return Quote{}, err
	}
	if slices.Contains([]bool{data.HasC(), data.HasD(), data.HasDp(), data.HasH(), data.HasL(), data.HasO(), data.HasPc()}, false) {
		return Quote{}, fmt.Errorf("Missing data for symbol %s", symbol)
	}
	quote := Quote{
		Symbol:        symbol,
		Price:         *data.C,
		Change:        *data.D,
		PercentChange: *data.Dp,
		DayHigh:       *data.H,
		DayLow:        *data.L,
		PreviousOpen:  *data.O,
		PreviousClose: *data.Pc,
	}
	return quote, err
}

func (c *FinHubClient) GetCandles(symbol string, resolution string, from int64, to int64) (Candles, error) {
	// Documentation: https://finnhub.io/docs/api/stock-candles
	// resolution: 1, 5, 15, 30, 60, D, W, M
	// from: Unix timestamp
	// to: Unix timestamp
	data, _, err := c.StockCandles(context.Background()).Symbol(symbol).Resolution(resolution).From(from).To(to).Execute()
	if err != nil {
		return Candles{}, err
	}
	if len(*data.O) != len(*data.H) || len(*data.O) != len(*data.L) || len(*data.O) != len(*data.C) || len(*data.O) != len(*data.V) || len(*data.O) != len(*data.T) {
		return Candles{}, fmt.Errorf("length of data is not equal")
	}
	if *data.S != "ok" {
		return Candles{}, fmt.Errorf("status is %s", *data.S)
	}
	candles := Candles{
		Open:      *data.O,
		High:      *data.H,
		Low:       *data.L,
		Close:     *data.C,
		Volume:    *data.V,
		Timestamp: *data.T,
	}
	return candles, nil
}

func (c *FinHubClient) OpenRealtimeStream(symbols []string) (*websocket.Conn, error) {
	// Documentation: https://finnhub.io/docs/api/websocket-trades
	apiKey := c.apiKey
	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", apiKey), nil)
	if err != nil {
		return nil, err
	}
	// Subscribe to the symbols
	for _, symbol := range symbols {
		subscribe := map[string]string{
			"type":   "subscribe",
			"symbol": symbol,
		}
		if err := w.WriteJSON(subscribe); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (c *FinHubClient) ReceiveRealtimeData(w *websocket.Conn) (chan Trade, chan int) {
	// Documentation: https://finnhub.io/docs/api/websocket-trades
	ch := make(chan Trade, 100)
	stop := make(chan int, 1)
	go func() {
		defer w.Close()
		for {
			// Check if stop signal is received
			select {
			case <-stop:
				return
			default:
				var res WebsocketResponse
				if err := w.ReadJSON(&res); err != nil {
					log.Fatal(err)
				} else {
					for _, trade := range res.Data {
						ch <- trade
					}
				}
			}
		}
	}()
	return ch, stop
}

func (c *FinHubClient) CloseRealtimeStream(stop chan int) {
	stop <- 1
}

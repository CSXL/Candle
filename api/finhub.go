// Contains wrappers for the FinHub API.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	finhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	websocket "github.com/gorilla/websocket"
)

type FinHubClient struct {
	*finhub.DefaultApiService
	apiKey string
}

type WebsocketResponse struct {
	Type string  `json:"type"`
	Data []Trade `json:"data"`
}

type Trade struct {
	Symbol    string  `json:"s"`
	Price     float64 `json:"p"`
	Volume    int     `json:"v"`
	Timestamp int64   `json:"t"`
}

type Candles struct {
	Open      []float64 `json:"o"`
	High      []float64 `json:"h"`
	Low       []float64 `json:"l"`
	Close     []float64 `json:"c"`
	Volume    []int     `json:"v"`
	Timestamp []int64   `json:"t"`
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

func (c *FinHubClient) GetCurrentPrice(symbol string) (finhub.Quote, error) {
	// Documentations: https://finnhub.io/docs/api/quote
	data, _, err := c.Quote(context.Background()).Symbol(symbol).Execute()
	return data, err
}

func (c *FinHubClient) GetCandles(symbol string, resolution string, from int64, to int64) (Candles, error) {
	// Documentation: https://finnhub.io/docs/api/stock-candles
	// resolution: 1, 5, 15, 30, 60, D, W, M
	// from: Unix timestamp
	// to: Unix timestamp
	_, res, err := c.StockCandles(context.Background()).Symbol(symbol).Resolution(resolution).From(from).To(to).Execute()
	if err != nil {
		return Candles{}, err
	}
	// Serialize the response
	var candles Candles
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return Candles{}, err
	}
	if err := json.Unmarshal(res_body, &candles); err != nil {
		return Candles{}, err
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

func (c *FinHubClient) ReceiveRealtimeData(w *websocket.Conn) (chan interface{}, chan int) {
	// Documentation: https://finnhub.io/docs/api/websocket-trades
	ch := make(chan interface{}, 100)
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
					ch <- err
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

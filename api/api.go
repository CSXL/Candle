package api

import (
	"encoding/json"
	"io"

	"log"
	"net/http"
	"time"
)

type MetaData struct {
	Information string `json:"1. Information"`
	Symbol      string `json:"2. Symbol"`
	LastRefresh string `json:"3. Last Refreshed"`
	Interval    string `json:"4. Interval"`
	OutputSize  string `json:"5. Output Size"`
	TimeZone    string `json:"6. Time Zone"`
}

type Stock struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type Response struct {
	MetaData MetaData         `json:"Meta Data"`
	Stocks   map[string]Stock `json:"Time Series (5min)"`
}

func Get(url string, params map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func RequestStockData(symbol string) []byte {
	api_url := "https://www.alphavantage.co/query"
	params := map[string]string{
		"function": "TIME_SERIES_INTRADAY",
		"symbol":   symbol,
		"interval": "5min",
		"apikey":   "demo",
	}
	body, err := Get(api_url, params)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func ParseStockData(data []byte) Response {
	var stock_data Response
	var err = json.Unmarshal([]byte(data), &stock_data)
	if err != nil {
		log.Fatal(err)
	}
	return stock_data
}

func GetStockData(symbol string) map[string]Stock {
	raw_data := RequestStockData(symbol)
	parsed_data := ParseStockData(raw_data)
	return parsed_data.Stocks
}

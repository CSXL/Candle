package api

// This program should request data from a stock API and return the data in an Object.
// The program should be able to handle errors and return a message if the API is down.
// This program should request from the alpha vantage API.

import (
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"time"
)

// Write a struct that resembles the metadata of the response.
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

// This is the function that requests data from the API. It should take in a stock symbol as a parameter.
func RequestStockData(symbol string) Response {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=%s&interval=5min&apikey=demo", symbol)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObj Response
	err = json.Unmarshal([]byte(body), &responseObj)
	if err != nil {
		log.Fatal(err)
	}
	return responseObj
}

// This is the function that returns the data from the API.
func GetStockData(symbol string) map[string]Stock {
	responseObj := RequestStockData(symbol)
	return responseObj.Stocks
}

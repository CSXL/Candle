package api

// This program should request data from a stock API and return the data in an Object.
// The program should be able to handle errors and return a message if the API is down.
// This program should request from the alpha vantage API.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// This is the struct that holds the data from the API.
type Stock struct {
	Open     float64 `json:"1. open"`
	High     float64 `json:"2. high"`
	Low      float64 `json:"3. low"`
	Close    float64 `json:"4. close"`
	Volume   float64 `json:"5. volume"`
	Adjusted float64 `json:"6. adjusted close"`
}

// This is the struct that holds the data from the API.
type TimeSeries struct {
	Stocks map[string]Stock `json:"Time Series (5min)"`
}

// This is the struct that holds the data from the API.
type Response struct {
	TimeSeries TimeSeries `json:"Time Series (5min)"`
}

// This is the function that requests data from the API.
func RequestStockData() Response {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// This is the URL to the API.
	url := "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=IBM&interval=5min&apikey=demo"
	// This is the request to the API.
	request, err := http.NewRequest("GET", url, nil)
	// This handles errors.
	if err != nil {
		log.Fatal(err)
	}
	// This sends the request to the API.
	response, err := client.Do(request)
	// This handles errors.
	if err != nil {
		log.Fatal(err)
	}
	// This closes the response.
	defer response.Body.Close()
	// This reads the response.
	body, err := ioutil.ReadAll(response.Body)
	// This handles errors.
	if err != nil {
		log.Fatal(err)
	}
	// This unmarshals the response.
	var responseObj Response
	var result map[string]any
	err = json.Unmarshal(body, &result)
	// This handles errors.
	if err != nil {
		log.Fatal(err)
	}
	// This returns the response.
	return responseObj
}

// This is the function that returns the data from the API.
func GetStockData() {
	// This requests data from the API.
	responseObj := RequestStockData()
	// This iterates through the response.
	for k, v := range responseObj.TimeSeries.Stocks {
		// This prints the data.
		fmt.Println(k, v)
	}
}

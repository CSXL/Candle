// Contains API for the application
package api

import (
	"time"

	"github.com/CSXL/Candle/api/finnhub"
)

type ApiClient struct {
	ApiKey        string
	FinnhubClient *finnhub.FinnhubClient
}

type Candles struct {
	Symbol     string
	Resolution string
	StartDate  string
	EndDate    string
	Open       []float32
	High       []float32
	Low        []float32
	Close      []float32
	Volume     []float32
	Timestamp  []int64
}

func NewApiClient(apiKey string) *ApiClient {
	finnhubClient := finnhub.NewFinnhubClient(apiKey)
	return &ApiClient{ApiKey: apiKey, FinnhubClient: finnhubClient}
}

func (c *ApiClient) GetCandles(symbol string, resolution string, startDate string, endDate string) (Candles, error) {
	// Time format: YYYY-MM-DD:HH-MM-SS
	// Example: 2022-01-01:00-00-00
	startDateUnix, err := dateTimeToUnix(startDate)
	if err != nil {
		return Candles{}, err
	}
	endDateUnix, err := dateTimeToUnix(endDate)
	if err != nil {
		return Candles{}, err
	}
	fetched_candles, err := c.FinnhubClient.GetCandles(symbol, resolution, startDateUnix, endDateUnix)
	if err != nil {
		return Candles{}, err
	}
	candles := Candles{
		Symbol:     symbol,
		Resolution: resolution,
		StartDate:  startDate,
		EndDate:    endDate,
		Open:       fetched_candles.Open,
		High:       fetched_candles.High,
		Low:        fetched_candles.Low,
		Close:      fetched_candles.Close,
		Volume:     fetched_candles.Volume,
		Timestamp:  fetched_candles.Timestamp,
	}
	return candles, nil
}

func dateTimeToUnix(dateTime string) (int64, error) {
	// Time format: YYYY-MM-DD:HH-MM-SS
	// Example: 2022-01-01:00-00-00
	reference_time := "2006-01-02:15-04-05"
	parsedTime, err := time.Parse(reference_time, dateTime)
	if err != nil {
		return 0, err
	}
	return parsedTime.Unix(), nil
}

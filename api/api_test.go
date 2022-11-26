package api

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	test_data, err := os.ReadFile("test_data.json")
	if err != nil {
		panic(err)
	}
	httpmock.RegisterResponder("GET", "https://www.alphavantage.co/query?apikey=demo&function=TIME_SERIES_INTRADAY&interval=5min&symbol=IBM",
		httpmock.NewStringResponder(200, string(test_data)))
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	url := "https://example.com/hello_world?apikey=demo"
	expected_message := "CSX Labs Rules!"
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, expected_message))

	res, err := Get(url, nil)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(res), expected_message)
}

func TestRequestStockData(t *testing.T) {
	res := RequestStockData("IBM")
	if string(res) == "" {
		t.Errorf("RequestStockData returned an empty string")
	}
}

func TestParseStockData(t *testing.T) {
	data, err := os.ReadFile("test_data.json")
	if err != nil {
		panic(err)
	}
	stock_data := ParseStockData(data)
	assert.Equal(t, stock_data.MetaData.Information, "Intraday (5min) open, high, low, close prices and volume")
	assert.Equal(t, stock_data.MetaData.Symbol, "IBM")
	assert.Equal(t, stock_data.MetaData.LastRefresh, "2096-09-10 16:00:00")
	assert.Equal(t, stock_data.MetaData.Interval, "5min")
	assert.Equal(t, stock_data.MetaData.OutputSize, "Compact")
	assert.Equal(t, stock_data.MetaData.TimeZone, "US/Eastern")
	assert.Equal(t, stock_data.Stocks["2096-09-10 16:00:00"].Open, "143.0000")
	assert.Equal(t, stock_data.Stocks["2096-09-10 16:00:00"].High, "143.0000")
	assert.Equal(t, stock_data.Stocks["2096-09-10 16:00:00"].Low, "143.0000")
	assert.Equal(t, stock_data.Stocks["2096-09-10 16:00:00"].Close, "143.0000")
	assert.Equal(t, stock_data.Stocks["2096-09-10 16:00:00"].Volume, "0")
}

func TestGetStockData(t *testing.T) {
	stock_data := GetStockData("IBM")
	if len(stock_data) == 0 {
		t.Errorf("GetStockData returned an empty map")
	}
}

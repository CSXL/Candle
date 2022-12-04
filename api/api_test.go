package api

import (
	"bytes"
	"io"
	"net/http"
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

	assert.Equal(t, expected_message, string(res))
}

// Test _BuildRequest with a valid URL
func TestBuildRequest(t *testing.T) {
	url := "https://example.com/hello_world?apikey=demo"
	req, err := _BuildRequest("GET", url, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, url, req.URL.String())
}

// Test _BuildRequest with an invalid URL
func TestBuildRequestInvalidMethod(t *testing.T) {
	url := "https://example.com/"
	_, err := _BuildRequest("INVALID METHOD", url, nil)
	assert.Error(t, err)
}

func TestGetUnparsableURL(t *testing.T) {
	bad_url := "this is not a URL"
	_, err := Get(bad_url, nil)
	assert.Error(t, err)
}

func TestGetConnectionRefused(t *testing.T) {
	url := "https://example.com/hello_world?apikey=demo"
	ErrConnectionRefused := &os.SyscallError{Syscall: "connect", Err: os.ErrDeadlineExceeded}
	httpmock.RegisterResponder("GET", url,
		httpmock.NewErrorResponder(ErrConnectionRefused))

	_, err := Get(url, nil)
	assert.Error(t, err)
}

func TestGetInvalidUTF8(t *testing.T) {
	url := "https://example.com/hello_world?apikey=demo"
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{ // Custom response to bypass httpmock's UTF-8 validation
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("\xff")),
			}, nil
		},
	)
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
	assert.Equal(t, "Intraday (5min) open, high, low, close prices and volume", stock_data.MetaData.Information)
	assert.Equal(t, "IBM", stock_data.MetaData.Symbol)
	assert.Equal(t, "2096-09-10 16:00:00", stock_data.MetaData.LastRefresh)
	assert.Equal(t, "5min", stock_data.MetaData.Interval)
	assert.Equal(t, "Compact", stock_data.MetaData.OutputSize)
	assert.Equal(t, "US/Eastern", stock_data.MetaData.TimeZone)
	assert.Equal(t, "143.0000", stock_data.Stocks["2096-09-10 16:00:00"].Open)
	assert.Equal(t, "143.0000", stock_data.Stocks["2096-09-10 16:00:00"].High)
	assert.Equal(t, "143.0000", stock_data.Stocks["2096-09-10 16:00:00"].Low)
	assert.Equal(t, "143.0000", stock_data.Stocks["2096-09-10 16:00:00"].Close)
	assert.Equal(t, "0", stock_data.Stocks["2096-09-10 16:00:00"].Volume)
}

func TestGetStockData(t *testing.T) {
	stock_data := GetStockData("IBM")
	if len(stock_data) == 0 {
		t.Errorf("GetStockData returned an empty map")
	}
}

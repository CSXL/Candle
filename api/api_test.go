package api

import (
	"log"
	"os"
	"testing"

	testutils "github.com/CSXL/Candle/api/testutils"
	gock "github.com/h2non/gock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var apiKey string

func TestMain(m *testing.M) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No test.env file found, running tests with environment variables.")
	}
	apiKey = os.Getenv("FINHUB_API_KEY")
	m.Run()
}

func TestGetStockCandles(t *testing.T) {
	SYMBOL := "AAPL"
	START_DATE := "2022-01-01:00-00-00" // YYYY-MM-DD:HH-MM-SS
	END_DATE := "2022-01-05:23-59-59"
	START_DATE_UNIX := "1640995200"
	END_DATE_UNIX := "1641427199"
	DATE_RANGE := 3 // Stock market is closed on weekends
	RESOLUTION := "D"

	type ResponseSchema struct {
		Open       []float32 `json:"o"`
		High       []float32 `json:"h"`
		Low        []float32 `json:"l"`
		Close      []float32 `json:"c"`
		Volume     []float32 `json:"v"`
		Timestamps []int64   `json:"t"`
		Status     string    `json:"s"`
	}

	sample_response, err := testutils.LoadSampleResponse[ResponseSchema]("candles")
	if err != nil {
		t.Error(err)
	}

	defer gock.Off()
	gock.New("https://finnhub.io/api/v1/stock/candle").
		MatchHeader("X-Finnhub-Token", apiKey).
		MatchParam("symbol", SYMBOL).
		MatchParam("resolution", RESOLUTION).
		MatchParam("from", START_DATE_UNIX).
		MatchParam("to", END_DATE_UNIX).
		Reply(200).
		JSON(sample_response)

	client := NewApiClient(apiKey)
	stockPrices, err := client.GetCandles(SYMBOL, RESOLUTION, START_DATE, END_DATE)
	assert.Nil(t, err)
	assert.Equal(t, SYMBOL, stockPrices.Symbol)
	assert.Equal(t, RESOLUTION, stockPrices.Resolution)
	assert.Equal(t, START_DATE, stockPrices.StartDate)
	assert.Equal(t, END_DATE, stockPrices.EndDate)
	assert.Equal(t, DATE_RANGE, len(stockPrices.Open))
	assert.Equal(t, DATE_RANGE, len(stockPrices.High))
	assert.Equal(t, DATE_RANGE, len(stockPrices.Low))
	assert.Equal(t, DATE_RANGE, len(stockPrices.Close))
	assert.Equal(t, DATE_RANGE, len(stockPrices.Volume))
}

func TestDateTimeToUnix(t *testing.T) {
	DATE_TIME := "2022-01-01:00-00-00"
	UNIX := int64(1640995200)
	unix, err := dateTimeToUnix(DATE_TIME)
	assert.Nil(t, err)
	assert.Equal(t, UNIX, unix)
}

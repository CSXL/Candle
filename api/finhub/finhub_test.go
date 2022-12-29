package api

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"testing"

	gock "github.com/h2non/gock"
	godotenv "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func LoadSampleResponse[T any](endpoint string) (T, error) {
	var empty T // Equivalent to nil when we return errors
	var sample_response T
	file, err := os.Open("sample_responses/" + endpoint + ".json")
	if err != nil {
		return empty, err
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		return empty, err
	}
	err = json.Unmarshal(contents, &sample_response)
	if err != nil {
		return empty, err
	}
	return sample_response, nil
}

var apiKey string

func TestMain(m *testing.M) {
	// Load .env file
	err := godotenv.Load("test.env")
	if err != nil {
		log.Println("No test.env file found, running tests with environment variables.")
	}
	apiKey = os.Getenv("FINHUB_API_KEY")
	m.Run()
}

func TestGetQuote(t *testing.T) {
	defer gock.Off()
	type QuoteSchema struct {
		Price         float32 `json:"c"`
		Change        float32 `json:"dp"`
		PercentChange float32 `json:"d"`
		DayHigh       float32 `json:"h"`
		DayLow        float32 `json:"l"`
		PreviousOpen  float32 `json:"o"`
		PreviousClose float32 `json:"pc"`
	}

	sample_response, err := LoadSampleResponse[QuoteSchema]("realtime_quote")
	if err != nil {
		t.Error(err)
	}

	SYMBOL := "AAPL"
	gock.New("https://finnhub.io/api/v1/quote").
		MatchHeader("X-Finnhub-Token", apiKey).
		MatchParam("symbol", SYMBOL).
		Reply(200).
		JSON(sample_response)

	client := NewFinHubClient(apiKey)
	quote, err := client.GetQuote(SYMBOL)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, SYMBOL, quote.Symbol)
	assert.Equal(t, sample_response.Price, quote.Price)
	assert.Equal(t, sample_response.Change, quote.Change)
	assert.Equal(t, sample_response.PreviousClose, quote.PreviousClose)
	assert.Equal(t, sample_response.PreviousOpen, quote.PreviousOpen)
	assert.Equal(t, sample_response.DayHigh, quote.DayHigh)
	assert.Equal(t, sample_response.DayLow, quote.DayLow)
	assert.Equal(t, sample_response.PercentChange, quote.PercentChange)
}

func TestGetCandles(t *testing.T) {
	type ResponseSchema struct {
		Close     []float32 `json:"c"`
		Open      []float32 `json:"o"`
		High      []float32 `json:"h"`
		Low       []float32 `json:"l"`
		Volume    []float32 `json:"v"`
		Timestamp []int64   `json:"t"`
		Status    string    `json:"s"`
	}

	sample_response, err := LoadSampleResponse[ResponseSchema]("stock_candles")
	if err != nil {
		t.Error(err)
	}

	START_TIME := int64(1569297600)
	END_TIME := int64(1569470400)
	RESOLUTION := "D"
	SYMBOL := "AAPL"

	gock.New("https://finnhub.io/api/v1/stock/candle").
		MatchHeader("X-Finnhub-Token", apiKey).
		MatchParam("symbol", SYMBOL).
		MatchParam("resolution", RESOLUTION).
		MatchParam("from", strconv.FormatInt(START_TIME, 10)).
		MatchParam("to", strconv.FormatInt(END_TIME, 10)).
		Reply(200).
		JSON(sample_response)

	client := NewFinHubClient(apiKey)
	candles, err := client.GetCandles(SYMBOL, RESOLUTION, START_TIME, END_TIME)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, SYMBOL, candles.Symbol)
	assert.Equal(t, sample_response.Open, candles.Open)
	assert.Equal(t, sample_response.High, candles.High)
	assert.Equal(t, sample_response.Low, candles.Low)
	assert.Equal(t, sample_response.Close, candles.Close)
	assert.Equal(t, sample_response.Volume, candles.Volume)
	assert.Equal(t, sample_response.Timestamp, candles.Timestamp)
}

func TestOpenRealtimeStream(t *testing.T) {
	// TODO: Mock gorilla websocket
	t.Skip("This test requires a websocket connection, which is not supported by gock.")
	client := NewFinHubClient(apiKey)
	w, err := client.OpenRealtimeStream([]string{"COINBASE:BTC-USD", "MSFT"})
	if err != nil {
		t.Error(err)
	}
	defer w.Close()
	for i := 0; i < 3; i++ {
		_, message, err := w.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		t.Log(string(message))
	}
}

func TestReceiveRealtimeData(t *testing.T) {
	// TODO: Mock gorilla websocket
	t.Skip("This test requires a websocket connection, which is not supported by gock.")
	client := NewFinHubClient(apiKey)
	w, err := client.OpenRealtimeStream([]string{"COINBASE:BTC-USD", "MSFT"})
	if err != nil {
		t.Error(err)
	}
	ch, stop := client.ReceiveRealtimeData(w)
	quote := <-ch
	t.Log("From Websocket: Current Price of", quote.Symbol, quote.Price)
	client.CloseRealtimeStream(stop)
}

// Path: api/finhub_test.go
package api

import (
	"log"
	"os"
	"testing"

	godotenv "github.com/joho/godotenv"
)

var apiKey string

func TestMain(m *testing.M) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey = os.Getenv("FINHUB_API_KEY")
	m.Run()
}

func TestGetCurrentPrice(t *testing.T) {
	client := NewFinHubClient(apiKey)
	quote, err := client.GetCurrentPrice("AAPL")
	if err != nil {
		t.Error(err)
	}
	t.Log(*quote.C)
}

func TestGetCandles(t *testing.T) {
	client := NewFinHubClient(apiKey)
	candles, err := client.GetCandles("AAPL", "D", 1671774551, 1672264493)
	if err != nil {
		t.Error(err)
	}
	t.Log("Candle First Opening Price: ", candles.Open[0])
}

func TestOpenRealtimeStream(t *testing.T) {
	client := NewFinHubClient(apiKey)
	w, err := client.OpenRealtimeStream([]string{"AAPL", "MSFT"})
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
	client := NewFinHubClient(apiKey)
	w, err := client.OpenRealtimeStream([]string{"AAPL", "MSFT"})
	if err != nil {
		t.Error(err)
	}
	ch, stop := client.ReceiveRealtimeData(w)
	quote := <-ch
	t.Log(quote)
	stop <- 1
}

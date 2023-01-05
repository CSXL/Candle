package chart

import (
	"testing"

	api "github.com/CSXL/Candle/api"
	testutils "github.com/CSXL/Candle/api/testutils"
	"github.com/charmbracelet/lipgloss"
)

var sampleCandles api.Candles

func loadCandles() {
	type ObjectSchema struct {
		Symbol     string    `json:"Symbol"`
		StartDate  string    `json:"StartDate"`
		EndDate    string    `json:"EndDate"`
		Resolution string    `json:"Resolution"`
		Close      []float32 `json:"Close"`
		High       []float32 `json:"High"`
		Low        []float32 `json:"Low"`
		Open       []float32 `json:"Open"`
		Volume     []float32 `json:"Volume"`
		Timestamp  []int64   `json:"Timestamp"`
	}
	data, err := testutils.LoadSampleResponse[ObjectSchema]("stock_candles")
	if err != nil {
		panic(err)
	}
	sampleCandles = api.Candles{
		Symbol:     data.Symbol,
		Resolution: data.Resolution,
		StartDate:  data.StartDate,
		EndDate:    data.EndDate,
		Close:      data.Close,
		High:       data.High,
		Low:        data.Low,
		Open:       data.Open,
		Volume:     data.Volume,
		Timestamp:  data.Timestamp,
	}
}

func TestMain(m *testing.M) {
	loadCandles()
	m.Run()
}

func TestNew(t *testing.T) {
	WIDTH := 80
	HEIGHT := 20
	TITLE := "Test Chart"
	COLORS := Colors{
		Primary: "#2c3e50",
		Text:    "#ecf0f1",
		Bearish: "#e74c3c",
		Bullish: "#1aa260",
	}
	chart := New(WIDTH, HEIGHT, TITLE, COLORS, sampleCandles, Range{Min: 0, Max: 0})
	if chart.Width != WIDTH {
		t.Errorf("Width was incorrect, got: %d, want: %d.", chart.Width, WIDTH)
	}
	if chart.Height != HEIGHT {
		t.Errorf("Height was incorrect, got: %d, want: %d.", chart.Height, HEIGHT)
	}
	if chart.Title != TITLE {
		t.Errorf("Title was incorrect, got: %s, want: %s.", chart.Title, TITLE)
	}
	if chart.Colors != COLORS {
		t.Errorf("Colors were incorrect, got: %v, want: %v.", chart.Colors, COLORS)
	}
}

func TestDraw(t *testing.T) {
	WIDTH := 80
	HEIGHT := 20
	TITLE := "Test Chart"
	COLORS := Colors{
		Primary: lipgloss.Color("#2c3e50"),
		Text:    lipgloss.Color("#ecf0f1"),
		Bearish: lipgloss.Color("#e74c3c"),
		Bullish: lipgloss.Color("#1aa260"),
	}
	chart := New(WIDTH, HEIGHT, TITLE, COLORS, sampleCandles, Range{Min: 0, Max: 0})
	chart.Draw()
	t.Logf("\n%s\n", chart.View())
}

package chart

import (
	"strings"

	api "github.com/CSXL/Candle/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Range struct {
	Min float32
	Max float32
}

type Model struct {
	Width          int
	Height         int
	content        string
	Title          string
	Candles        api.Candles
	Range          Range
	PrimaryColor   string
	SecondaryColor string
	TextColor      string
	BullishColor   string
	BearishColor   string
}

func New(width int, height int, title string, bullish_color string, bearish_color string, candles api.Candles, range_ Range) Model {
	return Model{
		Width:        width,
		Height:       height,
		Title:        title,
		Candles:      candles,
		content:      "",
		Range:        range_,
		BullishColor: bullish_color,
		BearishColor: bearish_color,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.content
}

func (m *Model) Draw() {
	var builder strings.Builder
	colors := map[string]lipgloss.Color{
		"bearish":   lipgloss.Color(m.BearishColor),
		"bullish":   lipgloss.Color(m.BullishColor),
		"primary":   lipgloss.Color(m.PrimaryColor),
		"secondary": lipgloss.Color(m.SecondaryColor),
		"text":      lipgloss.Color(m.TextColor),
	}
	title := lipgloss.NewStyle().
		Foreground(colors["text"]).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Center).
		Width(m.Width / 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors["primary"])
	builder.WriteString(title.Render(m.Title))
	builder.WriteString("\n")
	m.content = builder.String()
}

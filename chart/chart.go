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
	Width   int
	Height  int
	content string
	Title   string
	Candles api.Candles
	Range   Range
	Colors  Colors
}

type Colors struct {
	Primary lipgloss.Color
	Text    lipgloss.Color
	Bearish lipgloss.Color
	Bullish lipgloss.Color
}

func New(width int, height int, title string, colors Colors, candles api.Candles, range_ Range) Model {
	return Model{
		Width:   width,
		Height:  height,
		Title:   title,
		Candles: candles,
		content: "",
		Range:   range_,
		Colors:  colors,
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

func (m Model) titleView() string {
	title := lipgloss.NewStyle().
		Foreground(m.Colors.Text).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Colors.Primary)
	return title.Render(m.Title)
}

func (m *Model) Draw() {
	var builder strings.Builder
	title := m.titleView()
	builder.WriteString(title)
	builder.WriteString("\n")
	m.content = builder.String()
}

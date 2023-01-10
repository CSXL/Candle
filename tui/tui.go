package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type Frame int

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Align(lipgloss.Top).
			Width(80).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffffff"))
	itemFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1aa260")).
				Align(lipgloss.Left)
	itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

const (
	FrameQuote Frame = iota
	FrameChart
	FrameREPL
)

func (f Frame) String() string {
	switch f {
	case FrameQuote:
		return "Quote"
	case FrameChart:
		return "Chart"
	case FrameREPL:
		return "REPL"
	default:
		return "Unknown"
	}
}

type model struct {
	Frame Frame
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Frame > FrameQuote {
				m.Frame--
			}
		case "down":
			if m.Frame < FrameREPL {
				m.Frame++
			}
		case "1":
			m.Frame = FrameQuote
		case "2":
			m.Frame = FrameChart
		case "3":
			m.Frame = FrameREPL
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	s += titleStyle.Render("Candle - Bubbletea Boilerplate")
	s += "\n"
	for _, frame := range []Frame{FrameQuote, FrameChart, FrameREPL} {
		
		if m.Frame == frame {
			s += itemFocusedStyle.Render(" > " + frame.String())
		} else {
			s += itemStyle.Render(" > " + frame.String())
		}
		s += "\n"
	}
	s += helpStyle.Render(fmt.Sprintf("Current frame: %s • Use the arrow keys to navigate, or press 1, 2, or 3 to select a frame • Press Ctrl+C to quit", m.Frame.String()))
	return s
}

func Run() error {
	p := tea.NewProgram(
		model{},
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	return err
}

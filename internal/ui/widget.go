package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Chloezhu010/Interlude/internal/fun"
)

// Model holds the widget state
type Model struct {
	startTime   time.Time
	quitting    bool
	currentJoke string
	showingJoke bool
}

// New creates a new widget model
func New() Model {
	return Model{
		startTime: time.Now(),
	}
}

// Run launches the widget TUI (called by daemon)
func Run(startTime time.Time) error {
	model := Model{startTime: startTime}
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key presses and messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "j":
			m.currentJoke = fun.GetRandomJoke()
			m.showingJoke = true
			return m, nil
		case "n":
			if m.showingJoke {
				m.currentJoke = fun.GetRandomJoke()
				return m, nil
			}
		case "b", "esc":
			if m.showingJoke {
				m.showingJoke = false
				return m, nil
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	elapsed := time.Since(m.startTime).Truncate(time.Second)

	if m.showingJoke {
		b.WriteString(fmt.Sprintf("Claude working... %s\n\n", elapsed))
		b.WriteString("+-------------------------------------------+\n")
		wrapped := wordWrap(m.currentJoke, 39)
		for _, line := range wrapped {
			b.WriteString(fmt.Sprintf("| %-41s |\n", line))
		}
		b.WriteString("+-------------------------------------------+\n")
		b.WriteString("\n[n] Next joke   [b] Back   [q] Quit\n")
	} else {
		b.WriteString(fmt.Sprintf("Claude working... %s\n\n", elapsed))
		b.WriteString("[j] Tell me a joke\n")
		b.WriteString("[q] Quit\n")
	}

	return b.String()
}

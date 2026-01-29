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
	// trivia state
	currentQuestion *fun.Question
	showingTrivia   bool
	answered        bool
	wasCorrect      bool
	// flashcard state
	currentFlashcard *fun.Flashcard
	showingFlashcard bool
	revealed         bool
}

// New creates a new widget model
func New() Model {
	return Model{
		startTime: time.Now(),
	}
}

// Run launches the widget TUI (called by daemon)
func Run(startTime time.Time) error {
	// Load trivia in background (don't block startup)
	go fun.LoadTrivia()

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
			if m.showingTrivia && m.answered {
				m.currentQuestion = fun.GetRandomQuestion()
				m.answered = false
				m.wasCorrect = false
				return m, nil
			}
			if m.showingFlashcard && m.revealed {
				m.currentFlashcard = fun.GetRandomFlashcard()
				m.revealed = false
				return m, nil
			}
		case "b", "esc":
			if m.showingJoke {
				m.showingJoke = false
				return m, nil
			}
			if m.showingTrivia {
				m.showingTrivia = false
				m.answered = false
				return m, nil
			}
			if m.showingFlashcard {
				m.showingFlashcard = false
				m.revealed = false
				return m, nil
			}
		case "t":
			if !m.showingJoke && !m.showingTrivia && !m.showingFlashcard {
				m.currentQuestion = fun.GetRandomQuestion()
				m.showingTrivia = true
				m.answered = false
				return m, nil
			}
		case "f":
			if !m.showingJoke && !m.showingTrivia && !m.showingFlashcard {
				m.currentFlashcard = fun.GetRandomFlashcard()
				m.showingFlashcard = true
				m.revealed = false
				return m, nil
			}
		case " ":
			if m.showingFlashcard && !m.revealed {
				m.revealed = true
				return m, nil
			}
		case "1":
			if m.showingTrivia && !m.answered && m.currentQuestion != nil && len(m.currentQuestion.AllAnswers) > 0 {
				m.wasCorrect = m.currentQuestion.AllAnswers[0] == m.currentQuestion.CorrectAnswer
				m.answered = true
				return m, nil
			}
		case "2":
			if m.showingTrivia && !m.answered && m.currentQuestion != nil && len(m.currentQuestion.AllAnswers) > 1 {
				m.wasCorrect = m.currentQuestion.AllAnswers[1] == m.currentQuestion.CorrectAnswer
				m.answered = true
				return m, nil
			}
		case "3":
			if m.showingTrivia && !m.answered && m.currentQuestion != nil && len(m.currentQuestion.AllAnswers) > 2 {
				m.wasCorrect = m.currentQuestion.AllAnswers[2] == m.currentQuestion.CorrectAnswer
				m.answered = true
				return m, nil
			}
		case "4":
			if m.showingTrivia && !m.answered && m.currentQuestion != nil && len(m.currentQuestion.AllAnswers) > 3 {
				m.wasCorrect = m.currentQuestion.AllAnswers[3] == m.currentQuestion.CorrectAnswer
				m.answered = true
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
	} else if m.showingTrivia {
		b.WriteString(fmt.Sprintf("Claude working... %s\n\n", elapsed))
		if m.currentQuestion == nil {
			b.WriteString("No trivia questions loaded.\n")
			b.WriteString("\n[b] Back   [q] Quit\n")
		} else if m.answered {
			// Show result
			if m.wasCorrect {
				b.WriteString("✅ Correct!\n\n")
			} else {
				b.WriteString("❌ Wrong!\n")
				b.WriteString(fmt.Sprintf("Answer: %s\n\n", m.currentQuestion.CorrectAnswer))
			}
			b.WriteString("[n] Next question   [b] Back   [q] Quit\n")
		} else {
			// Show question
			b.WriteString("🧠 Quick Quiz\n\n")
			wrapped := wordWrap(m.currentQuestion.Question, 45)
			for _, line := range wrapped {
				b.WriteString(line + "\n")
			}
			b.WriteString("\n")
			labels := []string{"[1]", "[2]", "[3]", "[4]"}
			for i, ans := range m.currentQuestion.AllAnswers {
				if i < len(labels) {
					b.WriteString(fmt.Sprintf("%s %s\n", labels[i], ans))
				}
			}
			b.WriteString("\n[b] Back   [q] Quit\n")
		}
	} else if m.showingFlashcard {
		b.WriteString(fmt.Sprintf("Claude working... %s\n\n", elapsed))
		if m.currentFlashcard == nil {
			b.WriteString("No flashcards loaded.\n")
			b.WriteString("\n[b] Back   [q] Quit\n")
		} else if m.revealed {
			// Show answer
			b.WriteString("📖 Answer\n\n")
			b.WriteString(fmt.Sprintf("%s\n\n", m.currentFlashcard.Answer))
			if m.currentFlashcard.Explanation != "" {
				wrapped := wordWrap(m.currentFlashcard.Explanation, 45)
				for _, line := range wrapped {
					b.WriteString(line + "\n")
				}
				b.WriteString("\n")
			}
			b.WriteString("[n] Next card   [b] Back   [q] Quit\n")
		} else {
			// Show question
			b.WriteString("🧠 CS Flashcard\n\n")
			wrapped := wordWrap(m.currentFlashcard.Question, 45)
			for _, line := range wrapped {
				b.WriteString(line + "\n")
			}
			b.WriteString("\n[space] Reveal answer   [b] Back   [q] Quit\n")
		}
	} else {
		b.WriteString(fmt.Sprintf("Claude working... %s\n\n", elapsed))
		b.WriteString("[j] Tell me a dev joke\n")
		b.WriteString("[t] Trivia quiz\n")
		b.WriteString("[f] CS flashcard\n")
		b.WriteString("[q] Quit\n")
	}

	return b.String()
}

package fun

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const triviaURL = "https://opentdb.com/api.php?amount=50&category=18&type=multiple"

// Question represents a trivia question
type Question struct {
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
	AllAnswers       []string `json:"all_answers"` // shuffled options
}

var (
	questions  []Question
	questionMu sync.RWMutex
	triviaRng  *rand.Rand
)

func init() {
	triviaRng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GetRandomQuestion returns a random trivia question
func GetRandomQuestion() *Question {
	questionMu.RLock()
	defer questionMu.RUnlock()
	if len(questions) == 0 {
		return nil
	}
	q := questions[triviaRng.Intn(len(questions))]
	// Reshuffle answers each time for variety
	shuffleAnswers(&q)
	return &q
}

// LoadTrivia loads from cache or fetches from API
func LoadTrivia() error {
	questionMu.Lock()
	defer questionMu.Unlock()

	// Try cache first
	data, err := os.ReadFile(triviaCachePath())
	if err == nil {
		if err := json.Unmarshal(data, &questions); err == nil && len(questions) > 0 {
			return nil
		}
	}

	// Fetch from API
	return fetchTrivia()
}

func triviaCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".interlude", "trivia.json")
}

func fetchTrivia() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(triviaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("trivia API returned %d", resp.StatusCode)
	}

	// Limit response size to 1MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}

	var response struct {
		Results []Question `json:"results"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	// Decode HTML entities (API returns &quot; etc.)
	for i := range response.Results {
		response.Results[i].Question = html.UnescapeString(response.Results[i].Question)
		response.Results[i].CorrectAnswer = html.UnescapeString(response.Results[i].CorrectAnswer)
		for j := range response.Results[i].IncorrectAnswers {
			response.Results[i].IncorrectAnswers[j] = html.UnescapeString(response.Results[i].IncorrectAnswers[j])
		}
		// Build AllAnswers
		response.Results[i].AllAnswers = append(
			[]string{response.Results[i].CorrectAnswer},
			response.Results[i].IncorrectAnswers...,
		)
	}

	questions = response.Results

	// Save cache
	return saveTrivia()
}

func shuffleAnswers(q *Question) {
	triviaRng.Shuffle(len(q.AllAnswers), func(a, b int) {
		q.AllAnswers[a], q.AllAnswers[b] = q.AllAnswers[b], q.AllAnswers[a]
	})
}

func saveTrivia() error {
	dir := filepath.Dir(triviaCachePath())
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write
	tmp := triviaCachePath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, triviaCachePath())
}

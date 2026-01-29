package fun

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"sync"
	"time"
)

//go:embed flashcards_data.json
var flashcardsJSON []byte

// Flashcard represents a CS flashcard
type Flashcard struct {
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
}

var (
	flashcards   []Flashcard
	flashcardsMu sync.RWMutex
	flashcardRng *rand.Rand
)

func init() {
	flashcardRng = rand.New(rand.NewSource(time.Now().UnixNano()))
	json.Unmarshal(flashcardsJSON, &flashcards)
}

// GetRandomFlashcard returns a random CS flashcard
func GetRandomFlashcard() *Flashcard {
	flashcardsMu.RLock()
	defer flashcardsMu.RUnlock()
	if len(flashcards) == 0 {
		return nil
	}
	card := flashcards[flashcardRng.Intn(len(flashcards))]
	return &card
}

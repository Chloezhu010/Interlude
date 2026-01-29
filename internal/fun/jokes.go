package fun

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"sync"
	"time"
)

//go:embed jokes_data.json
var jokesJSON []byte

var (
	jokes   []string
	jokesMu sync.RWMutex
	rng     *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	json.Unmarshal(jokesJSON, &jokes)
}

// GetRandomJoke returns a random joke (thread-safe)
func GetRandomJoke() string {
	jokesMu.RLock()
	defer jokesMu.RUnlock()
	if len(jokes) == 0 {
		return "Why do programmers prefer dark mode? Because light attracts bugs."
	}
	return jokes[rng.Intn(len(jokes))]
}

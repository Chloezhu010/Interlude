package fun

import "testing"

func TestGetRandomJoke(t *testing.T) {
	joke := GetRandomJoke()
	if joke == "" {
		t.Error("GetRandomJoke returned empty string")
	}
}

func TestJokesLoaded(t *testing.T) {
	if len(jokes) == 0 {
		t.Error("No jokes loaded from embedded JSON")
	}
}

func TestRandomnessVaries(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 20; i++ {
		seen[GetRandomJoke()] = true
	}
	if len(seen) < 2 {
		t.Error("GetRandomJoke always returns same joke")
	}
}

func TestShowJokes(t *testing.T) {
	for i := 0; i < 5; i++ {
		joke := GetRandomJoke()
		t.Logf("Joke %d: %s", i+1, joke)
	}
}
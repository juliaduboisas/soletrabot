package game

import (
	"testing"
)

func TestGameAddWords(t *testing.T) {
	game := Game{Words: make(map[string]struct{})}
	game.AddWords("hello", "world")
	if _, exists := game.Words["hello"]; !exists {
		t.Errorf("Expected word 'hello' not found in game words")
	}
	if _, exists := game.Words["world"]; !exists {
		t.Errorf("Expected word 'world' not found in game words")
	}
}

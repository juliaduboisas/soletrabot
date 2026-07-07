package game

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestGameAddWords(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o')}

	// act
	game.AddWords([]string{"hello", "world"})

	// assert
	if !game.Words.Contains("hello") {
		t.Errorf("Expected word 'hello' not found in game words")
	}
	if game.Words.Contains("world") {
		t.Errorf("'world' should not be in words list")
	}
}

func TestWordValidation(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o')}

	// assert
	if !game.isValidWord("hello") {
		t.Errorf("'hello' should be a valid word")
	}
	if game.isValidWord("he") {
		t.Errorf("'he' shouldn't be a valid word, contains less than 3 letters")
	}
	if !game.isValidWord("hell") {
		t.Errorf("'hell' should be a valid word")
	}
}

func TestSetup(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o')}

	// act
	game.Setup([]rune{'a', 'b'})

	// assert
	if !game.Letters.Contains('a', 'b') {
		t.Errorf("Setup should have included 'a' and 'b'")
	}
	if game.Letters.Contains('h') {
		t.Errorf("Setup should have removed 'h'")
	}
}

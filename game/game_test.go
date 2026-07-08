package game

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestGameAddWordsGlobal(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "world"}, "teka")

	// assert
	if !game.Words.Contains("hello") {
		t.Errorf("Expected word 'hello' not found in game words")
	}
	if game.Words.Contains("world") {
		t.Errorf("'world' should not be in words list")
	}
}

func TestGameAddWordsUser(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "world"}, "teka")

	// assert
	if !game.PlayerWords["teka"].Contains("hello") {
		t.Errorf("Expected word 'hello' not found in player words")
	}
	if game.PlayerWords["teka"].Contains("world") {
		t.Errorf("'world' should not be in player words list")
	}
}

func TestGameDiff(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")
	game.AddWords([]string{"hello"}, "veter")

	difference := mapset.NewSet(game.GetDifference("veter")...)

	// assert
	if !difference.Contains("hollow", "howl") {
		t.Errorf("Difference between 'veter' words and global should be 'hollow' and 'howl'. Current diff: %v", difference)
	}
	if difference.Contains("hello") {
		t.Errorf("'hello' should not be in 'veter' difference. Current diff: %v", difference)
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
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello"}, "veter")
	game.Setup([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'})

	// assert
	if !game.Letters.Contains('a', 'b', 'c', 'd', 'e', 'f', 'g') {
		t.Errorf("Setup should have included 'a', 'b', 'c', 'd', 'e', 'f', 'g'")
	}
	if game.Letters.Contains('h') {
		t.Errorf("Setup should have removed 'h'")
	}
	if len(game.PlayerWords) > 0 {
		t.Errorf("Setup should have removed all player words")
	}
}

func TestSetupWithLessThanSevenLettersThrowsError(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello"}, "veter")
	_, err1 := game.Setup([]rune{'a', 'b', 'c', 'd', 'e', 'f'})
	_, err2 := game.Setup([]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})

	// assert
	if err1 == nil || err2 == nil {
		t.Errorf("Setup should not be done with number of letters that is not 7")
	}
}

func TestSync(t *testing.T) {
	// arrange
	game := Game{Words: mapset.NewSet[string](), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")
	game.AddWords([]string{"hello"}, "veter")

	game.SyncUser("veter")

	// assert
	if !game.PlayerWords["veter"].Contains("hollow", "howl", "hello") {
		t.Errorf("Sync between 'veter' words and global failed. Current 'veter' words: %v", game.PlayerWords["veter"])
	}
}

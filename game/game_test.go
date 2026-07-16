package game

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestGameAddWordsGlobal(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "world"}, "teka")

	// assert
	if _, exists := game.Words["hello"]; !exists {
		t.Errorf("Expected word 'hello' not found in game words")
	}
	if _, exists := game.Words["world"]; exists {
		t.Errorf("'world' should not be in words list")
	}
}

func TestGameAddWordsUser(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

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

func TestAddWordsWithWhitespace(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello ", " hell"}, "teka")

	// assert
	if _, exists := game.Words["hello"]; !exists {
		t.Errorf("Expected word 'hello' not found in words")
	}
	if _, exists := game.Words["hell"]; !exists {
		t.Errorf("Expected word 'hell' not found in words")
	}
}

func TestAddWordsWithUppercase(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"Hello ", " heLl"}, "teka")

	// assert
	if _, exists := game.Words["hello"]; !exists {
		t.Errorf("Expected word 'hello' not found in words")
	}
	if _, exists := game.Words["hell"]; !exists {
		t.Errorf("Expected word 'hell' not found in words")
	}
}

func TestGameDiff(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

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

func TestGameDiffForUserWithoutWordsReturnsAllWords(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow"}, "teka")
	difference := mapset.NewSet(game.GetDifference("veter")...)

	// assert
	if !difference.Contains("hello", "hollow") {
		t.Errorf("Difference for a user without any words should include all global words. Current diff: %v", difference)
	}
}

func TestWordValidation(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o')}

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
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

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
	if len(game.Words) > 0 {
		t.Errorf("Setup should have removed all words")
	}
}

func TestSetupWithLessThanSevenLettersThrowsError(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o'), PlayerWords: make(map[string]mapset.Set[string])}

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
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")
	game.AddWords([]string{"hello"}, "veter")

	game.SyncUser("veter")

	// assert
	if !game.PlayerWords["veter"].Contains("hollow", "howl", "hello") {
		t.Errorf("Sync between 'veter' words and global failed. Current 'veter' words: %v", game.PlayerWords["veter"])
	}
}

func TestBlameValid(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")
	game.AddWords([]string{"hello"}, "veter")

	player, err := game.Blame("hollow")

	// assert
	if player != "teka" {
		t.Errorf("Error in blame command: expected 'teka' got '%s'", player)
	}
	if err != nil {
		t.Errorf("Error in blame command: %v", err)
	}
}

func TestBlameInvalid(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")
	game.AddWords([]string{"hello"}, "veter")

	_, err := game.Blame("uau")

	// assert
	if err == nil {
		t.Errorf("Blame command should have thrown error")
	}
}

func TestBlameOrder(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")
	game.AddWords([]string{"hello"}, "veter")

	player, _ := game.Blame("hello")

	// assert
	if player != "teka" {
		t.Errorf("Error in blame command: expected 'teka' got '%s'", player)
	}
}

func TestSyncNewUser(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")

	ok := game.SyncUser("veter")

	// assert
	if !ok {
		t.Errorf("Error syncing new user")
	}
	if len(game.PlayerWords["veter"].ToSlice()) != 3 {
		t.Errorf("Sync had wrong result. Expected %v, got %v", []string{"hello", "hollow", "howl"}, game.PlayerWords["veter"].ToSlice())
	}
}

func TestGlobalWordCountValid(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	game.AddWords([]string{"hello", "hollow", "howl"}, "teka")

	globalWordCount, error := game.GlobalWordCount()

	// assert
	if globalWordCount != len(game.GetWords()) {
		t.Errorf("Word count had wrong result. Expected %d, got %d", len(game.GetWords()), globalWordCount)
	}
	if error != nil {
		t.Errorf("Error in count command error handling. Expected nil, got %v", error)
	}
}

func TestGlobalWordCountZero(t *testing.T) {
	// arrange
	game := Game{Words: make(map[string]string), Letters: mapset.NewSet('h', 'e', 'l', 'o', 'w'), PlayerWords: make(map[string]mapset.Set[string])}

	// act
	globalWordCount, error := game.GlobalWordCount()

	// assert
	if globalWordCount != 0 {
		t.Errorf("Word count had wrong result. Expected 0, got %d", globalWordCount)
	}
	if error == nil {
		t.Errorf("Error in count command error handling. Expected error, got nil.")
	}
}

func TestShowLeaderboardCountsWordsPerPlayer(t *testing.T) {
	// arrange
	game := Game{
		Words: map[string]string{
			"hello":  "teka",
			"hollow": "teka",
			"howl":   "veter",
		},
		Letters:     mapset.NewSet('h', 'e', 'l', 'o', 'w'),
		PlayerWords: make(map[string]mapset.Set[string]),
	}

	// act
	leaderboard := game.ShowLeaderboard()

	// assert
	var orderedLeaderboard []string

	for player := range leaderboard {
		orderedLeaderboard = append(orderedLeaderboard, player)
	}

	if len(leaderboard) != 2 {
		t.Fatalf("Expected 2 leaderboard entries, got %d", len(leaderboard))
	}
	if orderedLeaderboard[0] != "teka" {
		t.Errorf("Expected teka to be first with 2 words, got %+v", orderedLeaderboard[0])
	}
	if leaderboard[orderedLeaderboard[0]] != 2 {
		t.Errorf("Expected teka to have 2 words, got %+v", leaderboard[orderedLeaderboard[0]])
	}
	if orderedLeaderboard[1] != "veter" {
		t.Errorf("Expected veter to be second with 1 word, got %+v", orderedLeaderboard[1])
	}
	if leaderboard[orderedLeaderboard[1]] != 1 {
		t.Errorf("Expected veter to have 1 word, got %+v", leaderboard[orderedLeaderboard[1]])
	}
}

func TestOrderMapByValueSortsDescendingByCount(t *testing.T) {
	// arrange
	input := map[string]int{"alice": 2, "bob": 5, "carol": 1}

	// act
	ordered := orderMapByValue(input)

	// assert
	var orderedPlayers []string

	for player := range ordered {
		orderedPlayers = append(orderedPlayers, player)
	}

	if len(ordered) != 3 {
		t.Fatalf("Expected 3 ordered entries, got %d", len(ordered))
	}
	if orderedPlayers[0] != "bob" {
		t.Errorf("Expected bob to be first with value 5, got %+v", orderedPlayers[0])
	}
	if ordered[orderedPlayers[0]] != 5 {
		t.Errorf("Expected bob to have value 5, got %+v", ordered[orderedPlayers[0]])
	}
	if orderedPlayers[1] != "alice" {
		t.Errorf("Expected alice to be second with value 2, got %+v", orderedPlayers[1])
	}
	if ordered[orderedPlayers[1]] != 2 {
		t.Errorf("Expected alice to have value 2, got %+v", ordered[orderedPlayers[1]])
	}
	if orderedPlayers[2] != "carol" {
		t.Errorf("Expected carol to be third with value 1, got %+v", orderedPlayers[2])
	}
	if ordered[orderedPlayers[2]] != 1 {
		t.Errorf("Expected carol to have value 1, got %+v", ordered[orderedPlayers[2]])
	}
}

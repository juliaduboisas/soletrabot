package persistence

import (
	"sort"
	"strings"
	"testing"

	g "example.com/soletrabot/game"
	mapset "github.com/deckarep/golang-set/v2"
)

func TestCreateGameDataStringIncludesLettersWordsAndPlayerWords(t *testing.T) {
	// arrange
	game := g.Game{
		Letters: mapset.NewSet('a', 'b', 'c', 'd', 'e', 'f', 'g'),
		Words: map[string]string{
			"hello": "teka",
			"world": "veter",
		},
		PlayerWords: map[string]mapset.Set[string]{
			"teka":  mapset.NewSet("hello"),
			"veter": mapset.NewSet("world"),
		},
	}

	// act
	data := createGameDataString(game)

	// assert
	lines := strings.Split(strings.TrimSuffix(data, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("Expected serialized game data to contain sections, got %q", data)
	}
	if lines[0] != ":Letters" {
		t.Errorf("Expected first line to be the letters section header, got %q", lines[0])
	}
	if lines[len(lines)-1] != "hello" && lines[len(lines)-1] != "world" {
		t.Errorf("Expected last line to be one of the player words, got %q", lines[len(lines)-1])
	}

	lettersSectionEnd := -1
	for i, line := range lines {
		if line == ":GameWords" {
			lettersSectionEnd = i
			break
		}
	}
	if lettersSectionEnd == -1 {
		t.Fatalf("Expected letters section to end before the game words section")
	}

	letters := lines[1:lettersSectionEnd]
	sort.Strings(letters)
	expectedLetters := []string{"a", "b", "c", "d", "e", "f", "g"}
	if strings.Join(letters, "|") != strings.Join(expectedLetters, "|") {
		t.Errorf("Expected letters section to contain %v, got %v", expectedLetters, letters)
	}

	if lines[lettersSectionEnd] != ":GameWords" {
		t.Errorf("Expected game words section header at line %d, got %q", lettersSectionEnd, lines[lettersSectionEnd])
	}

	gameWordEntries := lines[lettersSectionEnd+1 : lettersSectionEnd+3]
	sort.Strings(gameWordEntries)
	expectedGameWords := []string{"hello,teka", "world,veter"}
	if strings.Join(gameWordEntries, "|") != strings.Join(expectedGameWords, "|") {
		t.Errorf("Expected game word entries to contain %v, got %v", expectedGameWords, gameWordEntries)
	}

	playerSectionEntries := lines[lettersSectionEnd+3:]
	if !containsAll(playerSectionEntries, []string{":teka", "hello", ":veter", "world"}) {
		t.Errorf("Expected player sections to contain the expected entries, got %v", playerSectionEntries)
	}
}

func containsAll(values []string, expected []string) bool {
	for _, expectedValue := range expected {
		found := false
		for _, value := range values {
			if value == expectedValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestConvertGameDataStringToGameParsesSerializedState(t *testing.T) {
	// arrange
	data := strings.Join([]string{
		":Letters",
		"a",
		"b",
		"c",
		":GameWords",
		"hello,teka",
		"world,veter",
		":teka",
		"hello",
		":veter",
		"world",
	}, "\n")

	// act
	loadedGame := convertGamaDataStringToGame([]byte(data))

	// assert
	if !loadedGame.Letters.Contains('a', 'b', 'c') {
		t.Errorf("Expected letters to be parsed from serialized state")
	}
	if loadedGame.Letters.Cardinality() != 3 {
		t.Errorf("Expected exactly 3 letters in the parsed collection, got %d", loadedGame.Letters.Cardinality())
	}
	if loadedGame.Words["hello"] != "teka" {
		t.Errorf("Expected 'hello' to be attributed to 'teka', got '%s'", loadedGame.Words["hello"])
	}
	if loadedGame.Words["world"] != "veter" {
		t.Errorf("Expected 'world' to be attributed to 'veter', got '%s'", loadedGame.Words["world"])
	}
	if len(loadedGame.Words) != 2 {
		t.Errorf("Expected exactly 2 game words, got %d", len(loadedGame.Words))
	}
	if !loadedGame.PlayerWords["teka"].Contains("hello") {
		t.Errorf("Expected 'hello' to be present for 'teka'")
	}
	if !loadedGame.PlayerWords["veter"].Contains("world") {
		t.Errorf("Expected 'world' to be present for 'veter'")
	}
	if loadedGame.PlayerWords["teka"].Cardinality() != 1 || loadedGame.PlayerWords["veter"].Cardinality() != 1 {
		t.Errorf("Expected exactly one word per player after parsing, got teka=%d, veter=%d", loadedGame.PlayerWords["teka"].Cardinality(), loadedGame.PlayerWords["veter"].Cardinality())
	}
}

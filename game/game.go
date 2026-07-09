package game

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type Game struct {
	Words       mapset.Set[string]
	Letters     mapset.Set[rune]
	PlayerWords map[string]mapset.Set[string]
}

func (G *Game) AddWords(words []string, player string) int {
	_, exists := G.PlayerWords[player]
	if !exists {
		G.PlayerWords[player] = mapset.NewSet[string]()
	}

	count := 0
	for _, word := range words {
		// add word for player
		if exists := G.PlayerWords[player].Contains(word); !exists && G.isValidWord(word) {
			G.PlayerWords[player].Add(word)
		}

		// add word globally
		if exists := G.Words.Contains(word); !exists && G.isValidWord(word) {
			G.Words.Add(word)
			G.PlayerWords[player].Add(word)
			count++
		}
	}
	return count
}

func (G *Game) GetWords() []string {
	return G.Words.ToSlice()
}

func (G *Game) Setup(letters []rune) ([]rune, error) {
	if len(letters) != 7 {
		return []rune{}, errors.New("Setup should be done with exactly 7 letters")
	}

	G.Letters.Clear()
	for k := range G.PlayerWords {
		delete(G.PlayerWords, k)
	}

	for _, letter := range letters {
		G.Letters.Add(letter)
	}
	return letters, nil
}

func (G *Game) GetDifference(user string) []string {
	playerWords, exists := G.PlayerWords[user]
	if !exists {
		return G.GetWords()
	}

	difference := G.Words.Difference(playerWords)
	return difference.ToSlice()
}

func (G *Game) SyncUser(user string) bool {
	playerWords, exists := G.PlayerWords[user]
	if !exists {
		return false
	}

	G.PlayerWords[user] = playerWords.Union(G.Words)
	return true
}

func (G Game) SaveGameState() (string, error) {
	dir, dirErr := os.Getwd()
	targetPath := filepath.Join(dir, "gameData.txt")

	if dirErr != nil {
		return "", dirErr
	}

	gameData := G.createGameDataString()

	// write with create or override
	writeErr := os.WriteFile(targetPath, []byte(gameData), 0644)

	if writeErr != nil {
		return "", writeErr
	}

	return targetPath, nil
}

func (G Game) createGameDataString() string {
	var sb strings.Builder

	sb.WriteString(":Letters\n")

	for _, letter := range G.Letters.ToSlice() {
		sb.WriteString(string(rune(letter)))
		sb.WriteString("\n")
	}

	sb.WriteString(":GameWords\n")

	for _, word := range G.GetWords() {
		sb.WriteString(word)
		sb.WriteString("\n")
	}

	for player := range G.PlayerWords {
		sb.WriteString(":")
		sb.WriteString(player)
		sb.WriteString("\n")

		for _, word := range G.GetWords() {
			sb.WriteString(word)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (G *Game) isValidWord(word string) bool {
	if len(word) < 4 {
		return false
	}

	for _, letter := range word {
		if !G.Letters.Contains(letter) {
			return false
		}
	}

	return true
}

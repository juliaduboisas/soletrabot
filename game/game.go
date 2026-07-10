package game

import (
	"cmp"
	"errors"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
)

type Game struct {
	Letters     mapset.Set[rune]
	Words       mapset.Set[string]
	PlayerWords map[string]mapset.Set[string]
}

func NewGame(letters mapset.Set[rune], words mapset.Set[string], playerWords map[string]mapset.Set[string]) *Game {
	return &Game{
		Letters:     letters,
		Words:       words,
		PlayerWords: playerWords,
	}
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
	return sortSizeFirst(G.Words.ToSlice())
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
	return sortSizeFirst(difference.ToSlice())
}

func (G *Game) SyncUser(user string) bool {
	playerWords, exists := G.PlayerWords[user]
	if !exists {
		return false
	}

	G.PlayerWords[user] = playerWords.Union(G.Words)
	return true
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

func sortSizeFirst(s []string) []string {
	slices.SortFunc(s, func(a, b string) int {
		if len(a) != len(b) {
			return cmp.Compare(len(a), len(b))
		}
		return cmp.Compare(a, b)
	})

	return s
}

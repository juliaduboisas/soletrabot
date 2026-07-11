package game

import (
	"cmp"
	"errors"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
)

type Game struct {
	Letters     mapset.Set[rune]
	Words       map[string]string
	PlayerWords map[string]mapset.Set[string]
}

func NewGame(letters mapset.Set[rune], words map[string]string, playerWords map[string]mapset.Set[string]) *Game {
	return &Game{
		Letters:     letters,
		Words:       words,
		PlayerWords: playerWords,
	}
}

func (G *Game) AddWords(wordsToAdd []string, player string) int {
	_, exists := G.PlayerWords[player]
	if !exists {
		G.PlayerWords[player] = mapset.NewSet[string]()
	}

	words := mapset.NewSet[string]()
	for word := range G.Words {
		words.Add(word)
	}

	count := 0
	for _, word := range wordsToAdd {
		// add word for player
		if exists := G.PlayerWords[player].Contains(word); !exists && G.isValidWord(word) {
			G.PlayerWords[player].Add(word)
		}

		// add word globally
		if exists := words.Contains(word); !exists && G.isValidWord(word) {
			G.Words[word] = player
			G.PlayerWords[player].Add(word)
			count++
		}
	}
	return count
}

func (G *Game) GetWords() []string {
	var words []string
	for word := range G.Words {
		words = append(words, word)
	}
	return sortSizeFirst(words)
}

func (G *Game) Setup(letters []rune) ([]rune, error) {
	if len(letters) != 7 {
		return []rune{}, errors.New("Setup should be done with exactly 7 letters")
	}

	G.Letters.Clear()
	for k := range G.Words {
		delete(G.Words, k)
	}
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

	words := mapset.NewSet[string]()
	for word := range G.Words {
		words.Add(word)
	}

	difference := words.Difference(playerWords)
	return sortSizeFirst(difference.ToSlice())
}

func (G *Game) SyncUser(user string) bool {
	playerWords, exists := G.PlayerWords[user]
	if !exists {
		G.PlayerWords[user] = mapset.NewSet[string]()
		playerWords, _ = G.PlayerWords[user]
	}

	words := mapset.NewSet[string]()
	for word := range G.Words {
		words.Add(word)
	}

	G.PlayerWords[user] = playerWords.Union(words)
	return true
}

func (G *Game) Blame(word string) (string, error) {
	user, exists := G.Words[word]
	if !exists {
		return "", errors.New("This word was not added dummy :)")
	}

	return user, nil
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

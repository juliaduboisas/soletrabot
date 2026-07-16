package game

import (
	"cmp"
	"errors"
	"slices"
	"sort"
	"strings"

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
		word = strings.ToLower(strings.TrimSpace(word))
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

func (G *Game) GlobalWordCount() (int, error) {

	totalWords := len(G.GetWords())

	if totalWords <= 0 {
		return 0, errors.New("There are no words added yet! You can be the first one ;D")
	}

	return totalWords, nil
}

func (G *Game) ShowLeaderboard() map[string]int {
	leaderboard := map[string]int{}

	for _, player := range G.Words {
		_, exists := leaderboard[player]
		if !exists {
			leaderboard[player] = 0
		}

		leaderboard[player]++
	}

	orderedLeaderboard := orderMapByValue(leaderboard)

	return orderedLeaderboard
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

func orderMapByValue(disordered map[string]int) map[string]int {
	keys := make([]string, 0, len(disordered))
	for k := range disordered {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return disordered[keys[i]] > disordered[keys[j]]
	})

	ordered := map[string]int{}

	for _, k := range keys {
		ordered[k] = disordered[k]
	}

	return ordered
}

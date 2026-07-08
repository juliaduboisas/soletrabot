package game

import mapset "github.com/deckarep/golang-set/v2"

type Game struct {
	Words       mapset.Set[string]
	Letters     mapset.Set[rune]
	PlayerWords map[string]mapset.Set[string]
}

func (G Game) AddWords(words []string, player string) int {
	playerWords, exists := G.PlayerWords[player]
	if !exists {
		playerWords = mapset.NewSet[string]()
		G.PlayerWords[player] = playerWords
	}

	count := 0
	for _, word := range words {
		if exists := G.Words.Contains(word); !exists && G.isValidWord(word) {
			G.Words.Add(word)
			playerWords.Add(word)
			count++
		}
	}
	return count
}

func (G Game) GetWords() []string {
	return G.Words.ToSlice()
}

func (G Game) Setup(letters []rune) []rune {
	G.Letters.Clear()
	for _, letter := range letters {
		G.Letters.Add(letter)
	}
	return letters
}

func (G Game) GetDifference(user string) []string {
	playerWords, exists := G.PlayerWords[user]
	if !exists {
		return G.GetWords()
	}

	difference := G.Words.Difference(playerWords)
	return difference.ToSlice()
}

func (G Game) isValidWord(word string) bool {
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

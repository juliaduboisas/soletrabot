package game

import mapset "github.com/deckarep/golang-set/v2"

type Game struct {
	Words   mapset.Set[string]
	Letters mapset.Set[rune]
}

func (G Game) AddWords(words []string) int {
	count := 0
	for _, word := range words {
		if exists := G.Words.Contains(word); !exists && G.isValidWord(word) {
			G.Words.Add(word)
			count++
		}
	}
	return count
}

func (G Game) GetWords() []string {
	return G.Words.ToSlice()
}

func (G Game) Setup(letters []rune) {
	G.Letters.Clear()
	for _, letter := range letters {
		G.Letters.Add(letter)
	}
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

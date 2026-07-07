package game

type Game struct {
	Words map[string]struct{}
}

type GameModifier interface {
	AddWords(words []string)
	GetWords() []string
}

func (G Game) AddWords(words []string) int {
	count := 0
	for _, word := range words {
		if _, exists := G.Words[word]; !exists {
			G.Words[word] = struct{}{}
			count++
		}
	}
	return count
}

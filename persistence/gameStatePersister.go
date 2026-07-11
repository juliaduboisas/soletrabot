package persistence

import (
	"os"
	"strings"

	g "example.com/soletrabot/game"
	mapset "github.com/deckarep/golang-set/v2"
)

type GameStatePersister struct {
	gameStateFilePath string
}

func NewGameStatePersister(gameStateFilePath string) *GameStatePersister {
	return &GameStatePersister{gameStateFilePath: gameStateFilePath}
}

func (P GameStatePersister) SaveGameState(game g.Game) (string, error) {
	gameData := createGameDataString(game)

	// write with create or override
	writeErr := os.WriteFile(P.gameStateFilePath, []byte(gameData), 0644)

	if writeErr != nil {
		return "", writeErr
	}

	return P.gameStateFilePath, nil
}

func (P GameStatePersister) LoadGameState() (*g.Game, error) {
	data, err := os.ReadFile(P.gameStateFilePath)

	if err != nil {
		return newEmptyGame(), err
	}

	textLines := strings.Split(string(data), "\n")

	letters := mapset.NewSet[rune]()
	words := make(map[string]string)
	playerWords := make(map[string]mapset.Set[string])

	nowAdding := -1
	user := ""

	for _, word := range textLines {
		switch word {
		case ":Letters":
			nowAdding = 0
			continue
		case ":GameWords":
			nowAdding = 1
			continue
		default:
			if strings.HasPrefix(word, ":") {
				user = strings.TrimSpace(strings.Replace(word, ":", "", 1))
				playerWords[user] = mapset.NewSet[string]()
				nowAdding = 2
				continue
			}
		}

		switch nowAdding {
		case 0:
			letters.Add([]rune(word)[0])
		case 1:
			wordAndAdder := strings.Split(word, ",")
			words[wordAndAdder[0]] = wordAndAdder[1]
		case 2:
			playerWords[user].Add(word)
		}
	}

	game := g.Game{Letters: letters, Words: words, PlayerWords: playerWords}

	return &game, nil
}

func newEmptyGame() *g.Game {

	return g.NewGame(mapset.NewSet[rune](), make(map[string]string), make(map[string]mapset.Set[string]))
}

func createGameDataString(game g.Game) string {
	var sb strings.Builder

	sb.WriteString(":Letters\n")

	for _, letter := range game.Letters.ToSlice() {
		sb.WriteString(string(rune(letter)))
		sb.WriteString("\n")
	}

	sb.WriteString(":GameWords\n")

	for word, player := range game.Words {
		sb.WriteString(word)
		sb.WriteString(",")
		sb.WriteString(player)
		sb.WriteString("\n")
	}

	for player := range game.PlayerWords {
		sb.WriteString(":")
		sb.WriteString(player)
		sb.WriteString("\n")

		playerWords := game.PlayerWords[player].ToSlice()
		for _, word := range playerWords {
			sb.WriteString(word)
			sb.WriteString("\n")
		}
	}

	result := sb.String()
	return strings.TrimSuffix(result, "\n")
}

func convertGamaDataStringToGame(data []byte) *g.Game {
	textLines := strings.Split(string(data), "\n")

	letters := mapset.NewSet[rune]()
	words := make(map[string]string)
	playerWords := make(map[string]mapset.Set[string])

	nowAdding := -1
	user := ""

	for _, word := range textLines {
		switch word {
		case ":Letters":
			nowAdding = 0
			continue
		case ":GameWords":
			nowAdding = 1
			continue
		default:
			if strings.HasPrefix(word, ":") {
				user = strings.TrimSpace(strings.Replace(word, ":", "", 1))
				playerWords[user] = mapset.NewSet[string]()
				nowAdding = 2
				continue
			}
		}

		switch nowAdding {
		case 0:
			letters.Add([]rune(word)[0])
		case 1:
			wordAndAdder := strings.Split(word, ",")
			words[wordAndAdder[0]] = wordAndAdder[1]
		case 2:
			playerWords[user].Add(word)
		}
	}

	game := g.Game{Letters: letters, Words: words, PlayerWords: playerWords}

	return &game
}

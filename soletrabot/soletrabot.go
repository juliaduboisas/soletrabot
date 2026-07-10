package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	g "example.com/soletrabot/game"
	p "example.com/soletrabot/persistence"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botToken := os.Getenv("TOKEN")
	webhookURL := os.Getenv("WEBHOOK_URL")

	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set up a webhook on Telegram side
	_ = bot.SetWebhook(ctx, &telego.SetWebhookParams{
		URL:         webhookURL,
		SecretToken: bot.SecretToken(),
	})

	// Receive information about webhook
	info, _ := bot.GetWebhookInfo(ctx)
	fmt.Printf("Webhook Info: %+v\n", info)

	// Create http serve mux
	mux := http.NewServeMux()

	// Get an update channel from webhook.
	// (more on configuration in examples/updates_webhook/main.go)
	updates, _ := bot.UpdatesViaWebhook(ctx, telego.WebhookHTTPServeMux(mux, "/bot", bot.SecretToken()))

	bh, _ := th.NewBotHandler(bot, updates)

	dir, _ := os.Getwd()
	gameStateFilePath := filepath.Join(dir, "gameData.txt")

	gameStatePersister := p.NewGameStatePersister(gameStateFilePath)

	game := g.NewGame(mapset.NewSet[rune](), mapset.NewSet[string](), make(map[string]mapset.Set[string]))

	if loadedGame, loadError := gameStatePersister.LoadGameState(); loadError == nil {
		game = loadedGame
	}

	// '/start' handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		_, _ = bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"Hello %s!", update.Message.From.FirstName,
		))
		return nil
	}, th.CommandEqual("start"))

	// '/add' handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		wordsInText := strings.Split(strings.ToLower(update.Message.Text), "\n")
		if len(wordsInText) < 2 {
			_, _ = bot.SendMessage(ctx, tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				"the /add command should contain at least one word.\n"+
					"Example:\n/add\n<word1>\n<word2>\n...",
			))
			return nil
		}

		var wordsToAdd []string
		for i := 1; i < len(wordsInText); i++ {
			wordsToAdd = append(wordsToAdd, wordsInText[i])
		}
		addedCount := game.AddWords(wordsToAdd, update.Message.From.Username)

		// Send message
		_, _ = bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%s added %v new words", update.Message.From.FirstName, addedCount,
		))
		return nil
	}, th.CommandEqual("add"))

	// '/get' handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		words := game.GetWords()

		// Send message
		_, _ = bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%s", strings.Join(words, "\n"),
		))
		return nil
	}, th.CommandEqual("get"))

	// '/setup' handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		commandMessageLines := strings.Split(strings.ToLower(update.Message.Text), "\n")
		if len(commandMessageLines) < 2 {
			_, _ = bot.SendMessage(ctx, tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				"The /setup command should have the letters in the second line.\n"+
					"Example:\n/setup\nabcdefg",
			))
			return nil
		}

		letters := []rune(commandMessageLines[1])
		setupLetters, err := game.Setup(letters)

		if err != nil {
			// Send message
			_, _ = bot.SendMessage(ctx, tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				"%v", err,
			))
			return nil
		}

		// Send message
		_, _ = bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%v", string(setupLetters),
		))
		return nil
	}, th.CommandEqual("setup"))

	// '/diff' handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		diffSlice := game.GetDifference(update.Message.From.Username)

		diff := strings.Join(diffSlice, "\n")

		// Send message
		_, _ = bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%s", diff,
		))
		return nil
	}, th.CommandEqual("diff"))

	// '/sync' handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		success := game.SyncUser(update.Message.From.Username)
		var message string
		if success {
			message = "Synced!"
		} else {
			message = "An error occurred during sync, try again later"
		}

		// Send message
		_, _ = bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%s", message,
		))
		return nil
	}, th.CommandEqual("sync"))

	// Start server for receiving requests from the Telegram
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Stop handling updates
	defer func() { _ = bh.Stop() }()

	// Start handling updates
	go func() {
		if err := bh.Start(); err != nil {
			log.Println(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	log.Println("Shutting down...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	_ = bh.Stop()

	webhook_del_err := bot.DeleteWebhook(shutdownCtx, &telego.DeleteWebhookParams{
		DropPendingUpdates: false,
	})
	if webhook_del_err != nil {
		log.Println("failed to delete webhook:", err)
	}

	_ = server.Shutdown(shutdownCtx)

	if gameDataPath, saveError := gameStatePersister.SaveGameState(*game); saveError != nil {
		log.Println("failed to save game data: ", saveError)
	} else {
		log.Println("saved game data to ", gameDataPath)
	}

	log.Println("Shutdown complete.")
}

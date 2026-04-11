package main

import (
	"SplitCore/internal/delivery/telegram"
	"SplitCore/internal/repository/postgres"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		slog.Error("BOT_TOKEN env var is missing")
		os.Exit(1)
	}

	settings := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	pool, err := postgres.NewPostgresPool(ctx, os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	userRepository := postgres.NewUserRepository(pool)
	fundRepository := postgres.NewFundRepository(pool)

	b, err := tele.NewBot(settings)
	if err != nil {
		slog.Error("Error creating bot", "error", err)
		os.Exit(1)
	}

	h := telegram.NewBotHandler(userRepository, fundRepository)
	h.SetupRegister(b)

	slog.Info("Starting bot", "version", "1.0.0", "env", "dev")
	b.Start()
}

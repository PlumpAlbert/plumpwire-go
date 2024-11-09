package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

var config Config

func main() {
	parse_json()
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)

	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(default_handler),
		bot.WithCallbackQueryDataHandler("devices", bot.MatchTypePrefix, device_list_handler),
	}

	bot_token := os.Getenv("TELEGRAM_TOKEN")
	b, err := bot.New(bot_token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

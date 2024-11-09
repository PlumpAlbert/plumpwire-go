package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		bot.WithDefaultHandler(handler),
		bot.WithCallbackQueryDataHandler("devices", bot.MatchTypePrefix, data_handler),
	}

	bot_token := os.Getenv("TELEGRAM_TOKEN")
	b, err := bot.New(bot_token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Devices", CallbackData: "devices"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Hi there, " + update.Message.From.Username + "!",
		ReplyMarkup: kb,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
			ChatID:    update.Message.Chat.ID,
		},
	})
}

func data_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("bot: ", b)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	username := update.CallbackQuery.From.Username
	fmt.Printf("Received message from %s", username)

	var clients []Client
	for _, c := range config.Clients {
		match, _ := regexp.MatchString(username+"\\s\\[(.+)\\]", c.Name)
		if match {
			clients = append(clients, c)
		}
		fmt.Printf("\n")
	}

	var lol []string
	for _, s := range clients {
		lol = append(lol, s.Name)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   strings.Join(lol, ", "),
	})
}

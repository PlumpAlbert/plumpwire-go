package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func default_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func device_list_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

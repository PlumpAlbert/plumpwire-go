package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CallbackQuery int

const (
	DEVICES CallbackQuery = iota
)

var Callbacks = map[CallbackQuery]string{
	DEVICES: "devices",
}

func default_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Список моих устройств", CallbackData: Callbacks[DEVICES]},
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

	var buttons []models.InlineKeyboardButton
	for _, c := range config.Clients {
		re := regexp.MustCompile(username + `\s\[(?P<Device>.+)\]`)

		if re == nil {
			return
		}

		matches := re.FindStringSubmatch(c.Name)
		if matches != nil {
			fmt.Printf("Matches: %s\n", matches)
			buttons = append(buttons, models.InlineKeyboardButton{
				Text: matches[re.SubexpIndex("Device")],
				CopyText: models.CopyTextButton{
					Text: c.IPAddress,
				},
			})
		}
	}

	if len(buttons) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "К сожалению, Вы не являетесь нашим пользователем",
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.CallbackQuery.Message.Message.ID,
				ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			},
		})
		return
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			buttons,
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Выбери конфигурацию, " + update.CallbackQuery.From.Username + "!",
		ReplyMarkup: kb,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.CallbackQuery.Message.Message.ID,
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		},
	})
}

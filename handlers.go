package main

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CallbackQuery int

const (
	DEVICES CallbackQuery = iota
	GET_CONFIG
)

var Callbacks = map[CallbackQuery]string{
	DEVICES:    "devices",
	GET_CONFIG: "config",
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
				Text:         matches[re.SubexpIndex("Device")],
				CallbackData: Callbacks[GET_CONFIG] + "/" + c.ID,
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
		ProtectContent: true,
		ChatID:         update.CallbackQuery.Message.Message.Chat.ID,
		Text:           "Выбери конфигурацию, " + update.CallbackQuery.From.Username + "!",
		ReplyMarkup:    kb,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.CallbackQuery.Message.Message.ID,
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		},
	})
}

func config_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	re := regexp.MustCompile(Callbacks[GET_CONFIG] + `/(?P<ID>.+)`)
	matches := re.FindStringSubmatch(update.CallbackQuery.Data)
	if matches == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			Text: "Что-то пошло не так...",
		})
		return
	}

	id := matches[re.SubexpIndex("ID")]
	doc, err := get_config(id)
	if err != nil {
		fmt.Println("Could not get config: " + err.Error())
		return
	}

	b.SendDocument(ctx, &bot.SendDocumentParams{
		Caption: "Ваша конфигурация подключения!",
		ChatID:  update.CallbackQuery.Message.Message.Chat.ID,
		Document: &models.InputFileUpload{
			Filename: "Wireguard." + update.CallbackQuery.From.Username + ".conf",
			Data:     bytes.NewReader(doc),
		},
	})
}

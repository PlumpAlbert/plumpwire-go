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
	fmt.Printf("Received message from %s\n", username)

	err := wg.GetClients()
	if err != nil {
		return
	}

	devices, err := wg.GetDevices(username)

	if len(devices) == 0 {
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

	var buttons []models.InlineKeyboardButton
	for _, c := range devices {
		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         c.Name,
			CallbackData: Callbacks[GET_CONFIG] + "/" + c.ID,
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ProtectContent: true,
		ChatID:         update.CallbackQuery.Message.Message.Chat.ID,
		Text:           "Выберите конфигурацию, которую хотите скачать",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				buttons,
			},
		},
		ReplyParameters: &models.ReplyParameters{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
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

	doc, err := wg.GetClientConfig(matches[re.SubexpIndex("ID")])
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

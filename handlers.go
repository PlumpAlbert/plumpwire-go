package main

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CallbackQuery int

const (
	DEVICES CallbackQuery = iota
	GET_CONFIG
	INVOICE
)

var Callbacks = map[CallbackQuery]string{
	DEVICES:    "devices",
	GET_CONFIG: "config",
	INVOICE:    "invoice",
}

func GetDefaultActions() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Список устройств", CallbackData: Callbacks[DEVICES]},
			},
			{
				{Text: "Статус подписки", CallbackData: Callbacks[INVOICE]},
			},
		},
	}
}

func default_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Hi there, " + update.Message.From.Username + "!",
		ReplyMarkup: GetDefaultActions(),
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

	devices, err := wg.GetDevices(username)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Что-то пошло не по плану...",
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.CallbackQuery.Message.Message.ID,
				ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			},
		})
		fmt.Printf("Error: %s\n", err)
		return
	}

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

	var buttons [][]models.InlineKeyboardButton
	for _, c := range devices {
		buttons = append(buttons, []models.InlineKeyboardButton{{
			Text:         c.DeviceName,
			CallbackData: Callbacks[GET_CONFIG] + "/" + c.ID,
		}})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ProtectContent: true,
		ChatID:         update.CallbackQuery.Message.Message.Chat.ID,
		Text:           "Выберите конфигурацию, которую хотите скачать",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
		ReplyParameters: &models.ReplyParameters{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		},
	})

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
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

func invoice_handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	client, err := im.GetClient(update.CallbackQuery.From.Username)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Кажется, Вы не являетесь нашим пользователем",
		})
		return
	}

	text := fmt.Sprintf("*Переплата*: %.2f ₽", client.PaymentBalance)

	invoices, err := im.GetRecurringInvoices(client)
	var invoiceTexts []string
	for _, i := range invoices {
		t := fmt.Sprintf(
			"📅 *Дата следующей оплаты*: %s\n💰 *Сумма к оплате* - %.2f ₽",
			time.Time(i.NextSendDate).Format("02.01.2006"),
			i.Amount,
		)
		invoiceTexts = append(invoiceTexts, t)
	}

	message := text + "\n\n" + strings.Join(invoiceTexts, "\n")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: GetDefaultActions(),
		Text: strings.ReplaceAll(
			strings.ReplaceAll(message, ".", "\\."),
			"-",
			"\\-",
		),
	})

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
}

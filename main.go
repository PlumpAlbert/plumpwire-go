package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"plumpalbert.xyz/plumpwire/invoice"
	"plumpalbert.xyz/plumpwire/models"
	"plumpalbert.xyz/plumpwire/wgez"
)

var config models.AppConfig
var wg wgez.WGEasy
var im *invoice.InvoiceManager

func LoadConfig() error {
	config.WG_HOST = os.Getenv("WG_HOST")
	if config.WG_HOST == "" {
		return errors.New("please provide WG_HOST environment variable")
	}

	config.TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	if config.TELEGRAM_TOKEN == "" {
		return errors.New("please provide TELEGRAM_TOKEN environment variable")
	}

	config.INVOICE_HOST = os.Getenv("INVOICE_HOST")
	if config.INVOICE_HOST == "" {
		return errors.New("please provide INVOICE_HOST environment variable")
	}

	return nil
}

func main() {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}

	wg = wgez.New(config.WG_HOST)
	im, err = invoice.New(config.INVOICE_HOST, os.Getenv("INVOICE_API_KEY"))

	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)

	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(default_handler),
		bot.WithCallbackQueryDataHandler(
			Callbacks[DEVICES],
			bot.MatchTypeExact,
			device_list_handler,
		),
		bot.WithCallbackQueryDataHandler(
			Callbacks[GET_CONFIG],
			bot.MatchTypePrefix,
			config_handler,
		),
	}

	b, err := bot.New(config.TELEGRAM_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

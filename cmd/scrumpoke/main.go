package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TelegramController struct {
	updates chan tgbotapi.Update
}

func (t *TelegramController) handleTelegramWebhook(c echo.Context) error {
	update := &tgbotapi.Update{}
	err := c.Bind(update)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	t.updates <- *update

	return c.String(http.StatusNoContent, "")
}

type Config struct {
	Telegram struct {
		Secret      string `env:"TELEGRAM_SECRET,notEmpty"`
		CallbackUrl string `env:"TELEGRAM_CALLBACK_URL" envDefault:""`
	}
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)

		panic(err)
	}

	var bot *tgbotapi.BotAPI
	var err error

	fmt.Println(cfg.Telegram.Secret)
	fmt.Println(cfg.Telegram.CallbackUrl)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	webhooks := e.Group("/webhooks/telegram")
	webhooks.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-telegram-bot-api-secret-token",
		Validator: func(key string, _ echo.Context) (bool, error) {
			return key == "testing-key", nil // todo: add config for secret key and update it via api on init
		},
	}))

	var tcUpdatesChannel tgbotapi.UpdatesChannel
	if cfg.Telegram.CallbackUrl != "" {
		bot, err = tgbotapi.NewBotAPIWithAPIEndpoint(
			cfg.Telegram.Secret,
			cfg.Telegram.CallbackUrl,
		)
		if err != nil {
			panic(err)
		}

		tcUpdatesChannel := make(chan tgbotapi.Update, 100)
		tc := TelegramController{
			updates: tcUpdatesChannel,
		}

		webhooks.POST("", tc.handleTelegramWebhook)

	} else {
		bot, err = tgbotapi.NewBotAPI(cfg.Telegram.Secret)
		if err != nil {
			panic(err)
		}

		bot.MakeRequest("deleteWebhook", tgbotapi.Params{})

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		tcUpdatesChannel = bot.GetUpdatesChan(u)
	}

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		panic(err)
	}

	fmt.Println("webhook set", webhookInfo.IsSet())
	fmt.Println("webhook url", webhookInfo.URL)

	go func() {
		for update := range tcUpdatesChannel {
			fmt.Println(update.UpdateID)

			if update.Message == nil {
				fmt.Println("not a message")
				continue
			}

			if !update.Message.IsCommand() {
				continue
			}

			// Create a new MessageConfig. We don't have text yet,
			// so we leave it empty.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "start":
				msg.Text = "Hello"
			case "poker":
				var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("1", "1"),
						tgbotapi.NewInlineKeyboardButtonData("2", "2"),
						tgbotapi.NewInlineKeyboardButtonData("3", "33"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("4", "4"),
						tgbotapi.NewInlineKeyboardButtonData("5", "5"),
						tgbotapi.NewInlineKeyboardButtonData("6", "6"),
					),
				)

				msg.Text = "Let's play a game!"
				msg.ReplyMarkup = numericKeyboard
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}

		}
	}()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

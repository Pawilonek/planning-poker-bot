package main

import (
	"context"
	"github.com/Pawilonek/scrumpoke/internal/config"
	"github.com/Pawilonek/scrumpoke/internal/poker"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)

		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Secret)
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.MakeRequest("deleteWebhook", tgbotapi.Params{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("webhook set", webhookInfo.IsSet())
	log.Println("webhook url", webhookInfo.URL)

	pokerBot := poker.NewPoker(bot, cfg.Jira)
	go pokerBot.Run()

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

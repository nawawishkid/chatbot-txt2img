package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleTelegramUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	if update.Message != nil {
		log.Printf("Received a message: %s", update.Message.Text)
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		imgUrl, err := createImage(update.Message.Text)

		if err != nil {
			return err
		}

		photoConfig := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(imgUrl))
		photoConfig.ReplyToMessageID = update.Message.MessageID
		// msg.ReplyToMessageID = update.Message.MessageID

		log.Printf("Sending generated image (%s)...", imgUrl)

		sentMsg, err := bot.Send(photoConfig)

		if err != nil {
			return fmt.Errorf("Error sending Telegram message: %w", err)
		}

		log.Printf("Message sent. Message ID was %d", sentMsg.MessageID)
	}

	return nil
}

func handleTelegramCallback(bot *tgbotapi.BotAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func initTelegramBot(botToken string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		return nil, err
	}

	wh, err := tgbotapi.NewWebhook(os.Getenv("APP_URL") + "/platforms/telegram/callback/" + bot.Token)

	if err != nil {
		return nil, err
	}

	_, err = bot.Request(wh)

	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/platforms/telegram/callback/" + bot.Token)

	go func() {
		for update := range updates {
			// log.Printf("%+v\n", update)

			if err := handleTelegramUpdate(bot, update); err != nil {
				log.Printf("Error handling Telegram update: %s", err)

				continue
			}
		}
	}()

	return bot, nil
}

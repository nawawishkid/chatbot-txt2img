package main

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func handleLINECallback(bot *linebot.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)

		if err != nil {
			log.Printf("Bot parse request error: %s", err)

			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(200)
			} else {
				w.WriteHeader(500)
			}

			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					imgUrl, err := createImage(message.Text)

					if err != nil {
						log.Printf("Error creating image: %s", err)

						w.WriteHeader(500)

						return
					}

					// reply to the user with the image
					if res, err := bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(imgUrl, imgUrl)).Do(); err != nil {
						log.Printf("bot reply message error: %s", err)
						// handle error
						w.WriteHeader(500)
						return
					} else {
						log.Printf("Message sent successfully. Request ID: %s", res.RequestID)
					}
				}
			}
		}

		w.WriteHeader(200)
	}
}

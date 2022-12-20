package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	// load environment variables from .env file
	log.Println("Loading environment variables...")

	if err := godotenv.Load(); err != nil {
		// handle error
		log.Printf("Error loading .env file: %s. Continue the process...", err)
	}

	// create LINE bot client
	log.Println("Creating LINE bot client...")

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	lineBot, err := linebot.New(channelSecret, channelAccessToken)

	if err != nil {
		panic(err)
	}

	_, err = initTelegramBot(os.Getenv("TELEGRAM_BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	http.Handle("/images/", loggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images")))))
	http.Handle("/platforms/line/callback", loggingMiddleware(handleLINECallback(lineBot)))

	var port int

	if os.Getenv("PORT") == "" {
		port = 8080
	} else {
		parsedPort, err := strconv.Atoi(os.Getenv("PORT"))

		if err != nil {
			panic(err)
		}

		port = parsedPort
	}

	// start the web server
	log.Printf("Starting HTTP server listening at port %d...", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		// handle error
		panic(err)
	}

}

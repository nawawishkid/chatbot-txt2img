package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"image"
	"image/color"
	"image/jpeg"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	// load environment variables from .env file
	log.Println("Loading environment variables...")

	if err := godotenv.Load(); err != nil {
		// handle error
		panic(err)
	}

	// create LINE bot client
	log.Println("Creating LINE bot client...")

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	bot, err := linebot.New(channelSecret, channelAccessToken)

	if err != nil {
		// handle error
		panic(err)
	}

	http.Handle("/images/", loggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images")))))
	http.Handle("/callback", loggingMiddleware(handleLINECallback(bot)))

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
					// // create image with the message text
					// img := createImage(message.Text)

					// // upload image to LINE server
					// response, err := bot.UploadContent(linebot.UploadContentArgs{
					// 	ContentType: "image/jpeg",
					// 	Content:     img,
					// })
					// if err != nil {
					// 	// handle error
					// 	w.WriteHeader(500)
					// 	return
					// }
					imgUrl, err := createImage(message.Text)

					if err != nil {
						log.Printf("Error creating image: %s", err)

						w.WriteHeader(500)

						return
					}

					// reply to the user with the image
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(imgUrl, imgUrl)).Do(); err != nil {
						log.Printf("bot reply message error: %s", err)
						// handle error
						w.WriteHeader(500)
						return
					}
				}
			}
		}

		w.WriteHeader(200)
	}
}

func createImage(text string) (string, error) {
	// create a new image
	img := image.NewRGBA(image.Rect(0, 0, 240, 240))

	fontBytes, err := ioutil.ReadFile("./fonts/Sarabun-Regular.ttf")

	if err != nil {
		return "", fmt.Errorf("Error reading font file: %w", err)
	}

	f, err := freetype.ParseFont(fontBytes)

	if err != nil {
		return "", err
	}

	fontFace := truetype.NewFace(f, &truetype.Options{})
	// draw the message text on the image
	drawText(img, text, fontFace, color.White)

	// create a buffer to hold the image
	buf := new(bytes.Buffer)

	// encode the image as a JPEG and write it to the buffer
	if err := jpeg.Encode(buf, img, nil); err != nil {
		// handle error
		return "", fmt.Errorf("Error encoding image to buffer: %w", err)
	}

	hash := sha256.New()

	if _, err := io.Copy(hash, buf); err != nil {
		return "", fmt.Errorf("Error hashing image: %w", err)
	}

	hashString := fmt.Sprintf("%x", hash.Sum(nil))
	fileName := fmt.Sprintf("%s.jpeg", hashString)
	relativeFilePath := fmt.Sprintf("public/images/%s", fileName)

	// Check if the image already exists
	_, err = os.Stat(relativeFilePath)

	if os.IsNotExist(err) {
		log.Printf("Image '%s' does not exists, creating...", fileName)

		file, err := os.Create(relativeFilePath)

		if err != nil {
			// handle error
			return "", fmt.Errorf("Error open file: %s", err)
		}

		defer file.Close()

		// encode the image as a JPEG and write it to the buffer
		if err := jpeg.Encode(file, img, nil); err != nil {
			return "", fmt.Errorf("Error writing image to file: %w", err)
		}
	} else {
		log.Printf("Image '%s' already exists", fileName)
	}

	// return the image as a byte slice
	return fmt.Sprintf("%s/images/%s", os.Getenv("APP_URL"), fileName), nil
}

func drawText(img *image.RGBA, text string, face font.Face, clr color.Color) {
	// get the dimensions of the text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(clr),
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	advance := d.MeasureString(text)

	// calculate the x and y coordinates to center the text
	x := (img.Rect.Dx() - advance.Ceil()) / 2
	y := (img.Rect.Dy() - d.Face.Metrics().Height.Ceil()) / 2

	// draw the text on the image
	d.Dot = fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
	d.DrawString(text)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform some processing before calling the next handler
		// log.Println("Received request:", r.URL)
		log.Printf("%s - %s", r.Method, r.URL)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Perform some processing after the next handler
		// log.Println("Completed request:", r.URL)
	})
}

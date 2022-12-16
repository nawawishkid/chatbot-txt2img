package main

import (
	"net/http"
	"os"

	"image"
	"image/color"
	"image/jpeg"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// handle error
		panic(err)
	}

	// create LINE bot client
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	bot, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		// handle error
		panic(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
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
					imgUrl := createImage(message.Text)

					// reply to the user with the image
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(imgUrl, imgUrl)).Do(); err != nil {
						// handle error
						w.WriteHeader(500)
						return
					}
				}
			}
		}
	})

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// handle error
		panic(err)
	}
}

func createImage(text string) string {
	// create a new image
	img := image.NewRGBA(image.Rect(0, 0, 240, 240))

	// draw the message text on the image
	drawText(img, text, basicfont.Face7x13, color.White)

	// create a buffer to hold the image
	// buf := new(bytes.Buffer)
	file, err := os.Create("image.jpeg")
	if err != nil {
		// handle error
	}
	defer file.Close()

	// encode the image as a JPEG and write it to the buffer
	if err := jpeg.Encode(file, img, nil); err != nil {
		// handle error
	}

	// return the image as a byte slice
	return "http://localhost:8080/image.jpeg"
}

func drawText(img *image.RGBA, text string, face font.Face, clr color.Color) {
	// get the dimensions of the text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(clr),
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	width, height := d.MeasureString(text)

	// calculate the x and y coordinates to center the text
	x := (img.Rect.Dx() - width.Ceil()) / 2
	y := (img.Rect.Dy() - height.Ceil()) / 2

	// draw the text on the image
	d.Dot = fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
	d.DrawString(text)
}

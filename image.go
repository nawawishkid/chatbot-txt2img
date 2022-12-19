package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

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

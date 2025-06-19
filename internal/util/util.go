package util

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// Improved "card-style" border with soft background and padding
func BorderedContainer(content fyne.CanvasObject) fyne.CanvasObject {
	// Create the background rectangle
	bg := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 255}) // soft card background
	bg.SetMinSize(fyne.NewSize(1, 1))

	// Set stroke (border) color and width
	bg.StrokeColor = color.NRGBA{R: 50, G: 50, B: 50, A: 25} // Example border color (semi-transparent black)
	bg.StrokeWidth = 2

	bg.CornerRadius = 10 // Use the theme's corner radius

	return container.NewVBox(
		container.NewStack(
			bg,
			container.New(layout.NewPaddedLayout(), content),
		),
	)
}

// randString returns a random alphanumeric string of length n.
// Useful for generating temp filenames or IDs.
func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func FileExists(p string) bool { _, err := os.Stat(p); return err == nil }

func FormatSec(f float64) string {
	sec := int(f + 0.5)
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}

func VSpacer(height float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.Transparent)
	r.SetMinSize(fyne.NewSize(1, height))
	return r
}

func HSpacer(width float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.Transparent)
	r.SetMinSize(fyne.NewSize(width, 1))
	return r
}

// EncodePNG encodes an image.Image into PNG byte slice for use in Fyne static resources.
func EncodePNG(img image.Image) []byte {
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

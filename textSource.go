package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/jpeg"
	"log"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	white  = color.RGBA{255, 255, 255, 255}
	blue   = color.RGBA{0, 0, 255, 255}
	red    = color.RGBA{255, 0, 0, 255}
	green  = color.RGBA{0, 255, 0, 255}
	yellow = color.RGBA{255, 255, 0, 255}
)

type TextBroadcasterSource struct {
	text string
}

func (source *TextBroadcasterSource) Init() error {
	return nil
}

func (source *TextBroadcasterSource) ReadFrame() ([]byte, error) {
	w := 800
	h := 600
	// color := blue
	// img := image.NewRGBA(image.Rect(0, 0, w, h))
	// draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	// animate dots to show activity
	txt := source.text
	dots := time.Now().Second() % 4
	for i := 0; i < dots; i++ {
		txt = fmt.Sprintf("%s.", txt)
	}
	for i := 0; i < 3-dots+1; i++ {
		txt = fmt.Sprintf("%s ", txt)
	}

	// start drawing
	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc := gg.NewContext(w, h)

	// draw blue background
	dc.SetColor(blue)
	dc.Clear()

	// draw text
	dc.SetFontFace(face)
	dc.SetColor(white)
	dc.DrawStringAnchored(txt, 412, 412, 0.5, 0.5)

	img := dc.Image()

	var buff bytes.Buffer
	jpeg.Encode(&buff, img, nil)
	return buff.Bytes(), nil
}

func (source TextBroadcasterSource) GetName() string {
	return "textsource"
}

func (source TextBroadcasterSource) Pause() {}

func (source TextBroadcasterSource) Unpause() {}

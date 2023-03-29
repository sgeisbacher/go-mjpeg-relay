package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
)

var (
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
	color := blue
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(im, im.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	var buff bytes.Buffer
	jpeg.Encode(&buff, im, nil)
	return buff.Bytes(), nil
}

func (source TextBroadcasterSource) GetName() string {
	return "textsource"
}

func (source TextBroadcasterSource) Pause() {}

func (source TextBroadcasterSource) Unpause() {}

package imager

import (
	"image"
	"log"
	"math/rand"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type ImageOptions struct {
	ImageWidth  int
	ImageHeight int
	Round       float64
	Font        font.Face
	LetterSize  float64
	LetterColor string
	BgColor     string
}

var Palette = map[string]string{
	"blue":          "#4674ca",
	"blue_dark":     "#315cac",
	"green":         "#57be8c",
	"green_dark":    "#3fa372",
	"yellow_orange": "#f9a66d",
	"red":           "#ec5e44",
	"red_dark":      "#e63717",
	"pink":          "#f868bc",
	"purple":        "#6c5fc7",
	"purple_dark":   "#4e3fb4",
	"teal":          "#57b1be",
	"gray":          "#847a8c",
}

func Image(letters string, opt *ImageOptions) *image.Image {
	opt.ValidOptions()
	w, h := opt.ImageWidth, opt.ImageHeight
	fw, fh := float64(w), float64(h)

	dc := gg.NewContext(w, h)

	if opt.Font == nil {
		if err := dc.LoadFontFace("./fonts/Roboto-Regular.ttf", opt.LetterSize); err != nil {
			log.Println(err)
			panic(err)
		}
	} else {
		dc.SetFontFace(opt.Font)
	}

	dc.DrawRoundedRectangle(0, 0, fw, fh, opt.Round)
	dc.SetHexColor(opt.BgColor) // bgColor
	dc.Fill()

	dc.SetHexColor(opt.LetterColor) // letters color
	dc.DrawStringAnchored(letters, fw/2, fh/2, 0.5, 0.38)
	image := dc.Image()
	return &image
}

func (op *ImageOptions) ValidOptions() {
	if op.Font == nil && op.LetterSize < 1 {
		op.LetterSize = 100
	}
	if op.ImageWidth < 1 {
		op.ImageWidth = 200
	}
	if op.ImageHeight < 1 {
		op.ImageHeight = 200
	}
	if op.BgColor == "" || string(op.BgColor[0]) != "#" {
		colors := []string{"blue", "blue_dark", "green", "green_dark", "yellow_orange", "red", "red_dark", "pink", "purple", "purple_dark", "teal", "gray"}
		color := colors[rand.Intn(12)]
		op.BgColor = Palette[color]
	}
	if op.LetterColor == "" {
		op.LetterColor = "fff"
	}
}

func min(first, second int) int {
	if first > second {
		return second
	}
	return first
}

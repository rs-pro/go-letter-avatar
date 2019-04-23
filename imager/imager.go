package imager

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var fBold *truetype.Font
var img *image.RGBA
var config Config

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Init() {
	if config.NoRandValues {
		config.ValidDatas()
	} else {
		config.randomLettars()
		config.randomBG()
		config.randomLC()
		config.similarColors()
		config.randomImageSize()
		config.randomLetterSize()
		config.randomRounding()
	}
	fontBytes, err := ioutil.ReadFile("./fonts/kievit-bold.ttf")
	if err != nil {
		panic(err)
	}
	fBold, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
}

func Template() *image.RGBA {
	var (
		dotX   fixed.Int26_6
		dotY   fixed.Int26_6
		height fixed.Int26_6
		width  fixed.Int26_6
		w      fixed.Int26_6
		h      fixed.Int26_6
	)

	bgColor := color.RGBA{config.BackgroundColor[0],
		config.BackgroundColor[1],
		config.BackgroundColor[2],
		config.BackgroundColor[3]}

	textColor := color.RGBA{config.LetterCollor[0],
		config.LetterCollor[1],
		config.LetterCollor[2],
		config.LetterCollor[3]}

	template := image.NewRGBA(image.Rect(0, 0, int(config.ImageWidth), int(config.ImageHeight)))

	draw.Draw(template, template.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

	face := truetype.NewFace(fBold, &truetype.Options{
		Size:    float64(config.LatterSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	d := &font.Drawer{
		Dst:  template,
		Src:  &image.Uniform{textColor},
		Face: face,
	}

	a, _ := d.BoundString(config.Latters)
	w = a.Max.X - a.Min.X
	h = a.Max.Y - a.Min.Y

	if a.Min.Y < 0 {
		height = -a.Min.Y
	}
	if a.Min.X > 0 {
		width = a.Min.X
	} else {
		width = -a.Min.X
	}

	dotX = (fixed.I(config.ImageWidth) - w) / 2
	dotY = (fixed.I(config.ImageHeight) - h) / 2

	d.Dot = fixed.Point26_6{
		X: dotX - width,
		Y: dotY + height,
	}

	d.DrawString(config.Latters)
	rounding(template)

	return template
}

func rounding(figure *image.RGBA) {
	if config.Rounding == 0 {
		return
	}

	r := min(config.ImageWidth, config.ImageHeight)

	if config.Rounding > r/2 {
		config.Rounding = r / 2
	}

	shiftDown := config.ImageHeight - 1 - (config.Rounding * 2)
	shiftRight := config.ImageWidth - 1 - (config.Rounding * 2)

	x0 := config.Rounding
	y0 := config.Rounding

	x, y, dx, dy := config.Rounding, 0, 1, 1
	err := dx - (int(float64(config.Rounding) * 1.35))

	for x >= y {

		figure.Set(x0-y, y0-x, color.Transparent) // левый верх верхней части
		figure.Set(x0-x, y0-y, color.Transparent) // левый низ верхней части

		figure.Set(x0-x, y0+y+shiftDown, color.Transparent) // левый верх нижней части
		figure.Set(x0-y, y0+x+shiftDown, color.Transparent) // левый низ нижней части

		figure.Set(x0+y+shiftRight, y0-x, color.Transparent) // правый угол верхняя половина половины
		figure.Set(x0+x+shiftRight, y0-y, color.Transparent) // правый угол нижняя половина половины

		figure.Set(x0+x+shiftRight, y0+y+shiftDown, color.Transparent) // низ право верхняя часть
		figure.Set(x0+y+shiftRight, y0+x+shiftDown, color.Transparent) // низ право нижняя часть

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (int(float64(config.Rounding) * 1.9))
		}
	}

	clippingLeft(figure)
	clippingRight(figure)
}
func clippingLeft(figure *image.RGBA) {
	var present bool
	for y := 0; y < config.ImageHeight; y++ {
		for x := 0; x < config.ImageWidth/2; x++ {
			fr, fg, fb, fa := figure.At(x, y).RGBA()

			if fr+fg+fb+fa == 0 {
				if x == 0 && y < config.ImageHeight/2 {
					present = true
					break
				} else {
					break
				}
			} else {
				figure.Set(x, y, color.Transparent)
			}
		}
		if present {
			yt := config.ImageHeight - config.Rounding
			for yt < config.ImageHeight {
				fr, fg, fb, fa := figure.At(0, yt).RGBA()

				if fr+fg+fb+fa != 0 {
					break
				}
				yt++
			}
			y = yt - 1
			present = false
		}
	}
}
func clippingRight(figure *image.RGBA) {
	var present bool
	for y := 0; y < config.ImageHeight; y++ {
		for x := config.ImageWidth - 1; x > config.ImageWidth/2; x-- {
			fr, fg, fb, fa := figure.At(x, y).RGBA()
			if fr+fg+fb+fa == 0 {
				if x == config.ImageWidth-1 && y < config.ImageHeight/2 {
					present = true
					break
				} else {
					break
				}
			} else {
				figure.Set(x, y, color.Transparent)
			}
		}
		if present {
			yt := config.ImageHeight - config.Rounding
			for yt < config.ImageHeight {
				fr, fg, fb, fa := figure.At(config.ImageWidth-1, yt).RGBA()
				if fr+fg+fb+fa != 0 {
					break
				}
				yt++
			}
			y = yt - 1
			present = false
		}
	}
}

// func ImageMain(c *gin.Context) {
// 	Init()
// 	img = Template()
// 	c.Writer.Header().Set("Content-Type", "image/png")
// 	err := png.Encode(c.Writer, img)
// 	pan(err)
// 	err = saveToDisk(img, "test.png")
// 	pan(err)
// }

package imager

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
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

func rounding(figure *image.RGBA) {

	if config.Rounding == 0 {
		return
	}
	var r int
	if config.ImageHeight > config.ImageWidth {
		r = config.ImageWidth
	} else {
		r = config.ImageHeight
	}
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

func Base(c *gin.Context) {
	config.ImageWidth, config.ImageHeight, config.Rounding = 500, 500, 25
	basedTemplate := image.NewRGBA(image.Rect(0, 0, int(config.ImageWidth), int(config.ImageHeight)))
	draw.Draw(basedTemplate, basedTemplate.Bounds(), &image.Uniform{image.White}, image.ZP, draw.Src)
	spew.Dump(config)
	rounding(basedTemplate)

	c.Writer.Header().Set("Content-Type", "image/png")
	err := png.Encode(c.Writer, basedTemplate)
	pan(err)
}

func Template() *image.RGBA {
	template := image.NewRGBA(image.Rect(0, 0, int(config.ImageWidth), int(config.ImageHeight)))
	config.LatterSize = 216
	bgColor := color.RGBA{config.BackgroundColor[0],
		config.BackgroundColor[1],
		config.BackgroundColor[2],
		config.BackgroundColor[3]}
	textColor := color.RGBA{config.LetterCollor[0],
		config.LetterCollor[1],
		config.LetterCollor[2],
		config.LetterCollor[3]}

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
	spew.Dump(config)

	a, _ := d.BoundString(config.Latters)
	w := a.Max.X - a.Min.X
	h := a.Max.Y - a.Min.Y
	log.Println("Ширина контейнера", w, "Высота контейнера", h, "BoundString Мин", a.Min, "Макс", a.Max)
	var dotY fixed.Int26_6
	log.Println("DOT X", (fixed.I(config.ImageWidth)-w)/2, "DOT Y", (fixed.I(config.ImageHeight)+h/4+2*h/4)/2)
	if (config.ImageHeight+int(h)/4+2*int(h)/4)/2 > config.ImageHeight/3 {
		dotY = (fixed.I(config.ImageHeight) + 4*h/6) / 2
	} else {
		dotY = (fixed.I(config.ImageHeight) + 5*h/6) / 2
	}

	d.Dot = fixed.Point26_6{
		X: (fixed.I(config.ImageWidth) - w) / 2,
		Y: dotY,
	}
	spew.Dump()
	d.DrawString(config.Latters)
	rounding(template)

	return template
}

func ImageMain(c *gin.Context) {
	Init()
	img = Template()
	c.Writer.Header().Set("Content-Type", "image/png")
	err := png.Encode(c.Writer, img)
	pan(err)
	err = saveToDisk(img, "test.png")
	pan(err)
}

func saveToDisk(img *image.RGBA, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		panic(err)
	}
	return nil
}

func pan(err error) {
	if err != nil {
		panic(err)
	}
}

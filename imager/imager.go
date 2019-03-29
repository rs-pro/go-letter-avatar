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

func Clip(figure *image.RGBA) {
	clipingTop("left", figure)
	clipingTop("right", figure)
	clipingBottom("left", figure)
	clipingBottom("right", figure)
	log.Println("Final")
}
func clipingBottom(flag string, figure *image.RGBA) {
	p := 99
	var x int
	if flag == "left" {
		x = 0
	} else if flag == "right" {
		x = 100
	} else {
		return
	}
	for i := x; i < x+101; i++ {
		for j := 200; j > p+x; j-- {
			// figure.Set(i, j, color.RGBA{255, 0, 155, 255})
			figure.Set(i, j, color.Transparent)
		}
		if flag == "left" {
			p++
		} else {
			p--
		}
	}
}

func clipingTop(flag string, figure *image.RGBA) {
	var p, x int
	if flag == "left" {
		p = 0
		x = 0
	} else if flag == "right" {
		p = 99
		x = 100
	} else {
		return
	}
	for i := x; i < x+101; i++ {
		for j := 100 - p; j > -1; j-- {
			// figure.Set(i, j, color.RGBA{255, 0, 155, 255})
			figure.Set(i, j, color.Transparent)
		}
		if flag == "left" {
			p++
		} else {
			p--
		}

	}
}

func Base(c *gin.Context) {
	basedTemplate := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(basedTemplate, basedTemplate.Bounds(), &image.Uniform{image.White}, image.ZP, draw.Src)
	Clip(basedTemplate)
	// clip := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// draw.Draw(clip, clip.Bounds(), &image.Uniform{color.RGBA{155, 155, 155, 0}}, image.ZP, draw.Src)
	// draw.DrawMask(basedTemplate, basedTemplate.Bounds(), clip, image.Point{-50, -50}, basedTemplate, image.Point{0, 0}, draw.Src)
	c.Writer.Header().Set("Content-Type", "image/png")
	err := png.Encode(c.Writer, basedTemplate)
	pan(err)
}

func Template() *image.RGBA {
	template := image.NewRGBA(image.Rect(0, 0, int(config.ImageWidth), int(config.ImageHeight)))
	bgColor := color.RGBA{config.BackgroundColor[0],
		config.BackgroundColor[1],
		config.BackgroundColor[2],
		config.BackgroundColor[3]}
	textColor := color.RGBA{config.LetterCollor[0],
		config.LetterCollor[1],
		config.LetterCollor[2],
		config.LetterCollor[3]}
	// spew.Dump(bgColor)
	// r, g, b, a := bgColor.RGBA()
	// log.Println("RGBA r:", r, "g:", g, "b:", b, "a:", a)
	// spew.Dump(bgColor)
	// r, g, b, a = textColor.RGBA()
	// log.Println("RGBA r:", r, "g:", g, "b:", b, "a:", a)
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

	y := config.ImageHeight/2 + (config.LatterSize / 3)
	d.Dot = fixed.Point26_6{
		X: (fixed.I(int(config.ImageWidth)) - d.MeasureString(config.Latters)) / 2,
		Y: fixed.I(int(y)),
	}

	d.DrawString(config.Latters)
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

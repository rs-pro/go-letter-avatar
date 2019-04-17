package imager

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
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

func circleRadius(r *int) *int {
	x := float64(config.ImageWidth / 2)
	y := float64(config.ImageHeight / 2)
	*r = int(math.Sqrt(x*x + y*y))
	return r
}

func rounding(figure *image.RGBA) {

	// сolorLeftTopAngle := color.RGBA{11, 188, 254, 255}
	// сolorRightTopAngle := color.RGBA{0, 255, 0, 255}
	// colorLeftBottom := color.RGBA{102, 255, 0, 255}
	// colorRightBottom := color.RGBA{255, 3, 62, 255}
	// colorTest := color.RGBA{255, 3, 62, 255}
	colorTest := color.Transparent
	radius := 110
	if radius > 200/2 {
		radius = 200 / 2
	}

	shiftDown := 200 - 1 - (radius * 2)
	shiftRight := 200 - 1 - (radius * 2)

	x0 := radius
	y0 := radius

	x, y, dx, dy := radius, 0, 1, 1
	err := dx - (int(float64(radius) * 1.35))

	for x >= y {

		figure.Set(x0-y, y0-x, colorTest) // левый верх верхней части
		figure.Set(x0-x, y0-y, colorTest) // левый низ верхней части

		figure.Set(x0-x, y0+y+shiftDown, colorTest) // левый верх нижней части
		figure.Set(x0-y, y0+x+shiftDown, colorTest) // левый низ нижней части

		figure.Set(x0+y+shiftRight, y0-x, colorTest) // правый угол верхняя половина половины
		figure.Set(x0+x+shiftRight, y0-y, colorTest) // правый угол нижняя половина половины

		figure.Set(x0+x+shiftRight, y0+y+shiftDown, colorTest) // низ право верхняя часть
		figure.Set(x0+y+shiftRight, y0+x+shiftDown, colorTest) // низ право нижняя часть

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (int(float64(radius) * 1.9))
		}
	}

	var bottomPoint int
	if radius == 200/2 {
		bottomPoint = radius
	} else {
		bottomPoint = 200 - radius
	}
	log.Println(bottomPoint)

	// clippingLeft(figure, colorTest)
	clippingRight(figure, colorTest)
}

func clippingLeft(figure *image.RGBA, col color.Color) {

	// colorF := color.RGBA{251, 7, 217, 255}
	colorF := color.Transparent

	for y := 0; y < 200; y++ {
		for x := 0; x < 200/2; x++ {
			cr, cg, cb, ca := col.RGBA()
			fr, fg, fb, fa := figure.At(x, y).RGBA()

			if fr == cr && fg == cg && fb == cb && fa == ca && x == 0 {
				figure.Set(x, y, colorF)
				fr, fg, fb, fa = figure.At(x, y+1).RGBA()
				break
			} else if fr == cr && fg == cg && fb == cb && fa == ca {
				figure.Set(x, y, colorF)
				fr, fg, fb, fa = figure.At(x+1, y).RGBA()
				if fr != cr && fg != cg && fb != cb && fa != ca {
					break
				}
			} else {
				figure.Set(x, y, colorF)
			}
		}
	}
}

func clippingRight(figure *image.RGBA, col color.Color) {

	// colorF := color.RGBA{251, 7, 217, 255}
	colorF := color.Transparent

	for y := 0; y < 200; y++ {
		for x := 200; x > 200/2; x-- {

			cr, cg, cb, ca := col.RGBA()
			fr, fg, fb, fa := figure.At(x, y).RGBA()
			if x > 40 && x < 50 {
				log.Println("x", x, "y", y)
				log.Println("cr", cr, "cg", cg, "cb", cb, "ca", ca)
				log.Println("fr", fr, "fg", fg, "fb", fb, "fa", fa)
			}
			if fr == cr && fg == cg && fb == cb && fa == ca && x == 199 {
				figure.Set(x, y, colorF)
				break
			} else if fr == cr && fg == cb && fb == cb && fa == ca {
				figure.Set(x, y, colorF)
				fr, fg, fb, fa = figure.At(x, y+1).RGBA()
				if fr != cr && fg != cg && fb != cb && fa != ca {
					break
				}
			} else {
				figure.Set(x, y, colorF)
			}
		}

	}
}

func Base(c *gin.Context) {
	config.ImageWidth, config.ImageHeight = 200, 200
	basedTemplate := image.NewRGBA(image.Rect(0, 0, int(config.ImageWidth), int(config.ImageHeight)))
	draw.Draw(basedTemplate, basedTemplate.Bounds(), &image.Uniform{image.White}, image.ZP, draw.Src)
	// var radius int = 100
	// drau(basedTemplate)
	// circleRadius(&radius)
	rounding(basedTemplate)

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

	var radius int
	circleRadius(&radius)
	config.LatterSize = radius / 2
	y := config.ImageHeight/2 + (config.LatterSize / 3)
	d.Dot = fixed.Point26_6{
		X: (fixed.I(config.ImageWidth) - d.MeasureString(config.Latters)) / 2,
		Y: fixed.I(y),
	}

	d.DrawString(config.Latters)
	// rounding(template)
	// clipping(template, bgColor)
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

// func drau(figure *image.RGBA) {
// 	figure.Set(0, 100, color.Black)
// 	figure.Set(4, 72, color.Black)
// 	figure.Set(4, 128, color.Black)
// 	figure.Set(20, 40, color.Black)
// 	figure.Set(20, 160, color.Black)
// 	figure.Set(40, 20, color.Black)
// 	figure.Set(40, 180, color.Black)
// 	figure.Set(72, 4, color.Black)
// 	figure.Set(72, 196, color.Black)
// 	figure.Set(100, 0, color.Black)
// 	figure.Set(128, 4, color.Black)
// 	figure.Set(128, 196, color.Black)
// 	figure.Set(160, 20, color.Black)
// 	figure.Set(160, 180, color.Black)
// 	figure.Set(180, 40, color.Black)
// 	figure.Set(180, 160, color.Black)
// 	figure.Set(196, 72, color.Black)
// 	figure.Set(196, 128, color.Black)
//
// 	сT := color.RGBA{255, 0, 0, 255}
// 	cB := color.RGBA{0, 255, 0, 255}
//
// 	figure.Set(0, 50, сT)
// 	figure.Set(100, 50, cB)
// 	figure.Set(2, 36, сT)
// 	figure.Set(2, 64, cB)
// 	figure.Set(10, 20, сT)
// 	figure.Set(10, 80, cB)
// 	figure.Set(36, 2, сT)
// 	figure.Set(36, 98, cB)
// 	figure.Set(50, 0, сT)
// 	figure.Set(64, 2, сT)
// 	figure.Set(64, 98, cB)
// 	figure.Set(80, 10, сT)
// 	figure.Set(80, 90, cB)
// 	figure.Set(90, 20, сT)
// 	figure.Set(90, 80, cB)
// 	figure.Set(98, 36, сT)
// 	figure.Set(98, 64, cB)
//
// }

// func clipping(figure *image.RGBA, col color.Color) {
// 	cr, cg, cb, ca := col.RGBA()
// 	for y := 0; y <= int(config.ImageHeight); y++ {
// 		for x := 0; x <= int(config.ImageWidth/2); x++ {
// 			fr, fg, fb, fa := figure.At(x, y).RGBA()
// 			if cr == fr && cg == fg && cb == fb && ca == fa {
// 				break
// 				} else {
// 					figure.Set(x, y, col)
// 				}
// 			}
// 			for x := int(config.ImageWidth); x >= int(config.ImageWidth/2); x-- {
// 				fr, fg, fb, fa := figure.At(x, y).RGBA()
// 				if cr == fr && cg == fg && cb == fb && ca == fa {
// 					break
// 					} else {
// 						figure.Set(x, y, col)
// 					}
// 				}
// 			}
// 		}

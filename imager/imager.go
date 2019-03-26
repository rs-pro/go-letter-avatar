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

type Config struct {
	Latters string
}

var fBold *truetype.Font
var img *image.RGBA

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Init() {
	log.Println("Start Init")
	fontBytes, err := ioutil.ReadFile("./fonts/kievit-bold.ttf")
	if err != nil {
		panic(err)
	}
	fBold, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	log.Println("Init Complete")
}

func prepareTemplate(w, h int) *image.RGBA {
	template := image.NewRGBA(image.Rect(0, 0, w, h)) //Возвращает новое изображение (прямоугольник)
	bgColor := color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256))}
	draw.Draw(template, template.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)
	size := h / 2
	face := truetype.NewFace(fBold, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	d := &font.Drawer{
		Dst:  template,
		Src:  image.White,
		Face: face,
	}
	str := "TS"
	y := h/2 + (size / 3)
	log.Println(y)
	d.Dot = fixed.Point26_6{
		X: (fixed.I(w) - d.MeasureString(str)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(str)
	return template

}
func ImageMain(c *gin.Context) { //Genral Image function
	log.Println("Start")
	Init()
	log.Println("Start prepareTemplate")
	img = prepareTemplate(200, 200)
	log.Println("End prepareTemplate")
	log.Println("Start MakeImage")
	log.Println("End MakeImage")
	c.Writer.Header().Set("Content-Type", "image/png")
	err := png.Encode(c.Writer, img)
	if err != nil {
		panic(err)
	}
	err = saveToDisk(img, "test.png")
	if err != nil {
		panic(err)
	}
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

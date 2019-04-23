package main

import (
	"bufio"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-letter-avatar/imager"
)

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to the generator of letters")
	})
	r.GET("/test", func(c *gin.Context) {
		imager.Init()
		img := imager.Template()
		c.Writer.Header().Set("Content-Type", "image/png")
		err := png.Encode(c.Writer, img)
		pan(err)
		err = saveToDisk(img, "test.png")
		pan(err)
	})
	log.Println("listening on :3001")
	err := http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatal(err)
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

func pan(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs-pro/go-letter-avatar/imager"
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
		if err != nil {
			panic(err)
		}
	})
	log.Println("listening on :3001")
	err := http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatal(err)
	}
}

func pan(err error) {
	if err != nil {
		panic(err)
	}
}

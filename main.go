package main

import (
	"log"
	"net/http"

	"github.com/go-letter-avatar/imager"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.GET("/test", imager.ImageMain)
	r.GET("/try", imager.Base)
	log.Println("listening on :3001")
	err := http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatal(err)
	}
}

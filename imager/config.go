package imager

import (
	"log"
	"math/rand"
)

type Config struct {
	Latters         string
	LatterSize      int
	LetterCollor    [4]uint8
	BackgroundColor [4]uint8
	ImageWidth      int
	ImageHeight     int
	NoRandValues    bool
	Rounding        int
}

func (c *Config) ValidDatas() {
	if len([]rune(c.Latters)) > 2 {
		runes := []rune(c.Latters)
		log.Println("Many chars")
		c.Latters = string(runes[0]) + string(runes[0])
	}
	if len(c.Latters) == 0 {
		log.Println("empty charset")
		c.randomLettars()
	}
	bgcolor := c.emptyArray("bg")
	lcolor := c.emptyArray("lc")
	if (bgcolor && lcolor) == true {
		c.similarColors()
	}
	c.CheckEndSetImageSize()
	c.letterSize()
}

func (c *Config) similarColors() {
	if c.BackgroundColor[0] == c.LetterCollor[0] || c.BackgroundColor[0] == c.LetterCollor[0]+10 || c.BackgroundColor[0] == c.LetterCollor[0]-10 {
		if c.BackgroundColor[1] == c.LetterCollor[1] || c.BackgroundColor[1] == c.LetterCollor[1]+10 || c.BackgroundColor[0] == c.LetterCollor[0]-10 {
			if c.BackgroundColor[2] == c.LetterCollor[2] || c.BackgroundColor[2] == c.LetterCollor[2]+10 || c.BackgroundColor[0] == c.LetterCollor[0]-10 {
				if c.BackgroundColor[3] == c.LetterCollor[3] || c.BackgroundColor[3] == c.LetterCollor[3]+10 || c.BackgroundColor[0] == c.LetterCollor[0]-10 {
					c.LetterCollor = [4]uint8{uint8(rand.Intn(255)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256))}
				}
			}
		}
	}
}

func (c *Config) emptyArray(flag string) bool {
	fl := false
	if flag == "bg" {
		for _, v := range c.BackgroundColor {
			if v != 0 {
				fl = true
				return false
			}
		}
		if !fl {
			c.randomBG()
			return true
		}
	} else {
		if flag == "lc" {
			for _, v := range c.LetterCollor {
				if v != 0 {
					fl = true
					return false
				}
			}
			if !fl {
				c.randomLC()
				return true
			}
		}
	}
	return false
}

func (c *Config) letterSize() {
	if c.LatterSize == 0 {
		if c.ImageWidth > c.ImageHeight {
			c.LatterSize = rand.Intn(int(float32(c.ImageHeight) * 0.55))
		} else {
			c.LatterSize = rand.Intn(int(float32(c.ImageWidth) * 0.55))
		}
	}
}

func (c *Config) randomLetterSize() {
	c.LatterSize = 0
	c.letterSize()
}

func (c *Config) randomRounding() {
	for c.Rounding < 1 {
		if c.ImageWidth > 0 {
			c.Rounding = rand.Intn(c.ImageWidth)
		} else {
			c.Rounding = rand.Intn(100)
		}
	}
}

func (c *Config) CheckEndSetImageSize() {
	for c.ImageHeight < 1 {
		c.ImageHeight = rand.Intn(7000)
	}
	for c.ImageWidth < c.ImageHeight {
		c.ImageWidth = rand.Intn(7000)
	}
}
func (c *Config) randomImageSize() {
	c.ImageHeight = 400
	c.ImageWidth = 600
	// c.CheckEndSetImageSize()
}

func (c *Config) randomBG() {
	c.BackgroundColor = [4]uint8{uint8(rand.Intn(255)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
}

func (c *Config) randomLC() {
	c.LetterCollor = [4]uint8{uint8(rand.Intn(255)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
}

func (c *Config) randomLettars() {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	c.Latters = string(letter[rand.Intn(len(letter))]) + string(letter[rand.Intn(len(letter))])
	if len([]rune(c.Latters)) < 2 {
		c.randomLettars()
	}
}

package imager

import "math/rand"

func (c *Config) RandomDatas() {
	c.randomLettars()
	c.randomBG()
	c.randomLC()
	c.similarColors()
	c.randomImageSize()
	c.randomLetterSize()
	c.randomRounding()
}

func (c *Config) randomLetterSize() {
	c.LetterSize = 0
	c.letterSize()
}

func (c *Config) randomRounding() {
	c.Rounding = 0
	for c.Rounding < 1 {
		c.Rounding = rand.Intn(min(c.ImageWidth, c.ImageHeight) / 2)
	}
}

func (c *Config) randomImageSize() {
	c.ImageHeight = 0
	c.ImageWidth = 0
	c.CheckEndSetImageSize()
}

func (c *Config) randomBG() {
	c.BackgroundColor = [4]uint8{uint8(rand.Intn(255)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
}

func (c *Config) randomLC() {
	c.LetterColor = [4]uint8{uint8(rand.Intn(255)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
}

func (c *Config) randomLettars() {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	c.Letters = string(letter[rand.Intn(len(letter))]) + string(letter[rand.Intn(len(letter))])
	if len([]rune(c.Letters)) < 2 {
		c.randomLettars()
	}
}

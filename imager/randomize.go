package imager

import "math/rand"

func (c *Config) randomLetterSize() {
	c.LatterSize = 0
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
	c.LetterCollor = [4]uint8{uint8(rand.Intn(255)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(255)}
}

func (c *Config) randomLettars() {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	c.Latters = string(letter[rand.Intn(len(letter))]) + string(letter[rand.Intn(len(letter))])
	if len([]rune(c.Latters)) < 2 {
		c.randomLettars()
	}
}

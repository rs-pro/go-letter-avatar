package imager

import (
	"testing"
)

func TestRandomImageSize(t *testing.T) {
	config := &Config{}
	config.randomImageSize()
	w, h := config.ImageWidth, config.ImageHeight
	config.randomImageSize()
	if config.ImageWidth == w && config.ImageHeight == h {
		t.Fatal("Width and Heigth was unchanged")
	} else {
		t.Log("success")
	}
}

func TestRandomLetterSize(t *testing.T) {
	config := &Config{}
	config.randomLetterSize()
	size := config.LetterSize
	config.randomLetterSize()
	if config.LetterSize == size && config.ImageWidth > 0 && config.ImageHeight > 0 {
		t.Fatal("Size of letter was unchanged")
	} else {
		t.Log("success")
	}
	size = 0
	config.ImageWidth, config.ImageHeight = 200, 200
	config.randomLetterSize()
	size = config.LetterSize
	config.randomLetterSize()
	if config.LetterSize == size {
		t.Fatal("Size of letter was unchanged")
	} else {
		t.Log("success")
	}
}

func TestRandomRounding(t *testing.T) {
	config := &Config{
		ImageWidth:  200,
		ImageHeight: 200,
	}
	config.randomRounding()
	rg := config.Rounding
	config.randomRounding()
	if config.Rounding == rg {
		t.Fatal("Rounding was unchanged")
	} else {
		t.Log("success")
	}
}

func TestRandomBG(t *testing.T) {
	config := &Config{}
	config.randomLC()
	bgColor := config.BackgroundColor
	config.randomBG()
	if config.BackgroundColor == bgColor {
		t.Fatal("Background color was unchanged")
	} else {
		t.Log("success")
	}
}

func TestRandomLC(t *testing.T) {
	config := &Config{}
	config.randomLC()
	lcColor := config.LetterColor
	config.randomLC()
	if config.LetterColor == lcColor {
		t.Fatal("Letter color was unchanged")
	} else {
		t.Log("success")
	}
}

func TestRandomLettars(t *testing.T) {
	config := &Config{}
	config.randomLettars()
	str := config.Letters
	config.randomLettars()
	if config.Letters == str {
		t.Fatal("Letters was unchanged")
	} else {
		t.Log("success")
	}
}

func TestMin(t *testing.T) {
	if min(3, 2) == 3 {
		t.Fatal("Wrong minimum")
	}
}

func TestEmptyArray(t *testing.T) {
	var flag bool
	config := &Config{}
	emptyArray := config.BackgroundColor
	flag = config.emptyArray("bg")
	if !flag {
		t.Fatal("Array isn't empty, background color was unchanged")
	} else {
		t.Log("success")
	}
	flag = false
	emptyArray = [4]uint8{0, 0, 0, 1}
	config.BackgroundColor = emptyArray
	flag = config.emptyArray("bg")
	if flag {
		t.Fatal("Array is empty, background color was changed")
	} else {
		t.Log("success")
	}
	flag = false
	emptyArray = config.LetterColor
	flag = config.emptyArray("lc")
	if !flag {
		t.Fatal("Array isn't empty, letter color color was unchanged")
	} else {
		t.Log("success")
	}
	flag = false
	emptyArray = [4]uint8{0, 0, 0, 1}
	config.LetterColor = emptyArray
	flag = config.emptyArray("lc")
	if flag {
		t.Fatal("Array is empty, letter color color was changed")
	} else {
		t.Log("success")
	}
}

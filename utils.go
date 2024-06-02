/*
@author: sk
@date: 2024/5/30
*/
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func NewImg(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

func SaveImg(img *image.RGBA) {
	file, err := os.Create("res/res.png")
	HandleErr(err)
	defer file.Close()
	err = png.Encode(file, img)
	HandleErr(err)
}

func OpenImage(path string) image.Image {
	file, err := os.Open(path)
	HandleErr(err)
	defer file.Close()
	img, _, err := image.Decode(file)
	HandleErr(err)
	return img
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NewRGB(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}
}

func NewRGBByVec(vec Vec) color.Color {
	return NewRGB(uint8(vec.X/0x100), uint8(vec.Y/0x100), uint8(vec.Z/0x100))
}

func Sign(val float64) float64 {
	if val > 0 {
		return 1
	} else if val < 0 {
		return -1
	} else {
		return 0
	}
}

package main

import (
	"image/color"
	"image/png"
	"os"

	sc ".." // github.com/TR-SLimey/ShapeCreator
)

func main() {

	palette := map[string]color.Color{
		"black":  color.RGBA{0, 0, 0, 255},
		"white":  color.RGBA{255, 255, 255, 255},
		"red":    color.RGBA{255, 0, 0, 255},
		"green":  color.RGBA{0, 255, 0, 255},
		"blue":   color.RGBA{0, 0, 255, 255},
		"orange": color.RGBA{255, 255, 0, 255},
		"yellow": color.RGBA{0, 255, 255, 255},
		"purple": color.RGBA{255, 0, 255, 255},
	}

	// Instantiate canvas
	var mainCanvas = sc.Canvas{}

	mainCanvas.Init(
		[2]int{512, 512},
	)
	go mainCanvas.DrawPixels()

	points, _ := mainCanvas.JoinPoints([][2]int{
		[2]int{0, 0},
		[2]int{512, 256},
		[2]int{0, 512},
		[2]int{512, 0},
		[2]int{0, 256},
	})

	pixels, _ := mainCanvas.PointsToPixels(points, palette["green"])
	mainCanvas.SendPixels(pixels)

	f, err := os.Create("draw.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	png.Encode(f, mainCanvas.GetResult())
}

package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	p ".." // github.com/TR-SLimey/GoPoints
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

	// Create a plane
	plane := p.Plane{}

	plane.Init(
		[2]int{512, 512},
	)

	square := []p.Point{
		p.Point{10, 10},
		p.Point{50, 10},
		p.Point{50, 50},
		p.Point{10, 50},
		p.Point{10, 10},
	}

	points := plane.JoinPoints(square)

	planeDimensions := plane.GetDimensions()

	img := image.NewRGBA(image.Rect(0, 0, planeDimensions[0], planeDimensions[1]))

	for _, point := range points {
		img.Set(point.X, point.Y, palette["red"])
	}

	f, err := os.Create("example.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_ = png.Encode(f, img)
}

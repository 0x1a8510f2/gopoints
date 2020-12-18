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
		"none":   color.RGBA{0, 0, 0, 0},
	}

	// Create a plane
	plane := p.Plane{}

	plane.Init(
		[2]int{512, 512},
	)

	planeDimensions := plane.GetDimensions()

	img := image.NewRGBA(image.Rect(0, 0, planeDimensions[0], planeDimensions[1]))

	// Draw a square
	square := []p.Point{
		p.Point{X: 10, Y: 10},
		p.Point{X: 50, Y: 10},
		p.Point{X: 50, Y: 50},
		p.Point{X: 10, Y: 50},
		p.Point{X: 10, Y: 10},
	}
	points := plane.JoinPoints(square)

	plane.Flip(1) // We're about to draw the points onto an image which has a reversed Y axis, so flip to account for this

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

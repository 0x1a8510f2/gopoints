package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"

	p "github.com/TR-SLimey/gopoints"
)

func main() {

	// Define a colour palette for easier management of colours later
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
	// ...and init it to the right size
	plane.Init(
		[2]int{512, 512},
	)
	// Get the size of the plane
	// Technically unnecessary if we keep track of it above but left here for demo purposes
	planeDimensions := plane.GetDimensions()
	// Create the image. We will add pixels later
	img := image.NewRGBA(image.Rect(0, 0, planeDimensions[0], planeDimensions[1]))

	// Drawing a square

	// The points of a square
	square := []p.Point{
		p.Point{X: 10, Y: 10},
		p.Point{X: 50, Y: 10},
		p.Point{X: 50, Y: 50},
		p.Point{X: 10, Y: 50},
		p.Point{X: 10, Y: 10},
	}
	// Join them
	points := plane.JoinPoints(square)
	// Add them to the plane (non-strict so no need to check error)
	_ = plane.WritePoints(points, false)

	// Drawing a triangle

	// The points of the triangle
	triangle := []p.Point{
		p.Point{X: 150, Y: 190},
		p.Point{X: 80, Y: 70},
		p.Point{X: 220, Y: 70},
		p.Point{X: 150, Y: 190},
	}
	// Join the points
	points = plane.JoinAndFillPoints(triangle)
	// Add them to the plane (non-strict so no need to check error)
	_ = plane.WritePoints(points, false)

	// Drawing a random shape

	randomShape := []p.Point{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		randomShape = append(randomShape, p.Point{
			X: rand.Intn(500-200) + 200,
			Y: rand.Intn(500-200) + 200,
		})
	}
	// Join the points
	points = plane.JoinPoints(randomShape)
	// Add them to the plane (non-strict so no need to check error)
	_ = plane.WritePoints(points, false)

	plane.Flip(1) // We're about to draw the points onto an image which has a reversed Y axis, so flip along the Y axis to account for this

	// Fetch the points from the plane and draw them as pixels on the image
	for _, point := range plane.ReadPoints() {
		img.Set(point.X, point.Y, palette["red"])
	}

	f, err := os.Create("example.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_ = png.Encode(f, img)
}

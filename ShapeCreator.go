package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
)

// Structs
type pixel struct {
	Col color.Color
	Pos [2]int
}

type canvas struct {
	postInit bool
	Cols     map[string]color.Color
	canvas   image.RGBA
	size     [2]int
	buffer   chan pixel
}

func (canvas *canvas) Init(size [2]int) {
	canvas.size = size
	canvas.canvas = *image.NewRGBA(image.Rect(0, 0, canvas.size[0], canvas.size[1]))
	canvas.buffer = make(chan pixel, 0)
	canvas.postInit = true
}

func (canvas *canvas) JoinPoints(points [][2]int) ([][2]int, error) {
	// Go does not provide sets, so to avoid duplicate points due to rounding, use map keys
	// at little extra cost as structs take 0 bytes
	allPoints := make(map[[2]int]struct{})

	pPoint := [2]int{}
	noPPoint := true

	for _, cPoint := range points {
		if noPPoint {
			noPPoint = false
		} else {
			// Work out difference in X and difference in Y between the two points
			posDiffs := []float64{
				float64(cPoint[0] - pPoint[0]),
				float64(cPoint[1] - pPoint[1]),
			}
			// Work out length of line between points (round to nearest pixel) (pythagoras)
			posDiffLength := int(
				math.Round(
					math.Sqrt(math.Pow(posDiffs[0], 2) + math.Pow(posDiffs[1], 2)),
				),
			)
			// Work out the gradient of the line between the two points
			posDiffGradient := posDiffs[1] / posDiffs[0]

			// We now have all the information we need to join the two points

			// The coordinates of the point to be generated/drawn. We need to keep track of this
			// as the location of the next point is calculated from this. We start, of course, with
			// the initial point itself
			cDrawPoint := [2]float64{
				float64(pPoint[0]),
				float64(pPoint[1]),
			}

			// What we should increment the x coordinate by to generate the next pixel. This is simply the
			// ratio of the diagonal length to the x-only length
			xIncrement := posDiffs[0] / float64(posDiffLength)

			// What we should increment the y coordinate by to generate the next pixel. This is the x increment
			// multiplied by the ratio of x to y for proportion
			yIncrement := xIncrement * posDiffGradient

			// For each point needing to be generated
			// TODO: Skip points which round to the same coordinated for efficiency
			for i := 0; i < posDiffLength; i++ {
				// Add the point to the set of points
				point := [2]int{int(math.Round(cDrawPoint[0])), int(math.Round(cDrawPoint[1]))}
				//fmt.Println(point)
				allPoints[point] = struct{}{}
				// Calculate the position of the next point
				cDrawPoint[0] += xIncrement
				cDrawPoint[1] += yIncrement
			}
		}
		pPoint = cPoint
	}

	// Convert our set implementation to a simple array for compatibility
	finalPoints := [][2]int{}
	for element := range allPoints {
		finalPoints = append(finalPoints, element)
	}

	return canvas.FlipPoints(finalPoints), nil
}

func (canvas *canvas) FlipPoints(points [][2]int) [][2]int {
	flipped := [][2]int{}
	for _, point := range points {
		point = [2]int{
			point[0],
			canvas.size[1] - point[1],
		}
		flipped = append(flipped, point)
	}
	return flipped
}

func (canvas *canvas) PointsToPixels(points [][2]int, col string) ([]pixel, error) {
	// TODO: Check if col is valid
	pixels := []pixel{}
	for _, point := range points {
		pixel := pixel{
			Pos: point,
			Col: canvas.Cols[col],
		}
		pixels = append(pixels, pixel)
	}
	return pixels, nil
}

func (canvas *canvas) SendPixels(pixels []pixel) {
	for _, pixel := range pixels {
		canvas.buffer <- pixel
	}
}

func (canvas *canvas) DrawPixels(timeout time.Duration) {
	for pixel := range canvas.buffer {
		canvas.canvas.Set(pixel.Pos[0], pixel.Pos[1], pixel.Col)
	}
}

func (canvas *canvas) GetResult() image.Image {
	//canvasSnapshot := canvas.canvas
	return &canvas.canvas
}

func main() {

	// Instantiate canvas
	var mainCanvas = canvas{
		Cols: map[string]color.Color{
			"black":  color.RGBA{0, 0, 0, 255},
			"white":  color.RGBA{255, 255, 255, 255},
			"red":    color.RGBA{255, 0, 0, 255},
			"green":  color.RGBA{0, 255, 0, 255},
			"blue":   color.RGBA{0, 0, 255, 255},
			"orange": color.RGBA{255, 255, 0, 255},
			"yellow": color.RGBA{0, 255, 255, 255},
			"purple": color.RGBA{255, 0, 255, 255},
		},
	}
	mainCanvas.Init([2]int{512, 512})
	go mainCanvas.DrawPixels(-1)

	points, _ := mainCanvas.JoinPoints([][2]int{
		[2]int{0, 0},
		[2]int{512, 256},
		[2]int{0, 512},
		[2]int{512, 0},
		[2]int{0, 256},
	})

	pixels, _ := mainCanvas.PointsToPixels(points, "red")
	mainCanvas.SendPixels(pixels)

	f, err := os.Create("draw.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	png.Encode(f, mainCanvas.GetResult())
}

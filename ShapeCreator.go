package ShapeCreator

import (
	"image"
	"image/color"
	"math"
)

// Structs
type Pixel struct {
	Colr color.Color
	Pos  [2]int
}

type Canvas struct {
	postInit bool
	canvas   image.RGBA
	size     [2]int
	buffer   chan Pixel
}

func (canvas *Canvas) Init(size [2]int) {
	canvas.size = size
	canvas.canvas = *image.NewRGBA(image.Rect(0, 0, canvas.size[0], canvas.size[1]))
	canvas.buffer = make(chan Pixel, 0)
	canvas.postInit = true
}

func (canvas *Canvas) JoinPoints(points [][2]int) ([][2]int, error) {
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

	// Convert our set implementation to a simple array for compatibility with the other methods
	finalPoints := [][2]int{}
	for element := range allPoints {
		finalPoints = append(finalPoints, element)
	}
	finalPoints = canvas.FlipPoints(finalPoints)

	return finalPoints, nil
}

func (canvas *Canvas) FlipPoints(points [][2]int) [][2]int {
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

func (canvas *Canvas) PointsToPixels(points [][2]int, colr color.Color) ([]Pixel, error) {
	pixels := []Pixel{}
	for _, point := range points {
		pixel := Pixel{
			Pos:  point,
			Colr: colr,
		}
		pixels = append(pixels, pixel)
	}
	return pixels, nil
}

func (canvas *Canvas) SendPixels(pixels []Pixel) {
	for _, pixel := range pixels {
		canvas.buffer <- pixel
	}
}

func (canvas *Canvas) DrawPixels() {
	for pixel := range canvas.buffer {
		canvas.canvas.Set(pixel.Pos[0], pixel.Pos[1], pixel.Colr)
	}
}

func (canvas *Canvas) GetResult() image.Image {
	//canvasSnapshot := canvas.canvas
	return &canvas.canvas
}

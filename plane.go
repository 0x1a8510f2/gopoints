package gopoints

import (
	"fmt"
	"math"
)

// A special type used by Plane.ReadPointsByFilter
type filterFunction func(Point) bool

// Plane is essentially a wrapper around PointSet, but with a specific size for dimension X and Y
// (although these are not strictly enforced by any of the methods) and some methods to interact
// with points both on the plane and passed as an argument.
type Plane struct {
	data       PointSet
	dimensions [2]int
}

// Init initialises the plane with dimensions passed as an array of 2 integers, and creates a PointSet
// to store the points on the plane.
func (pln *Plane) Init(dimensions [2]int) {
	pln.dimensions = dimensions
	pln.data = PointSet{}
}

// WritePoints writes points passed as the first argument to the plane's PointSet. It can optionally
// reject points which are outside of the plane's dimensions if the second argument is `true`.
func (pln *Plane) WritePoints(points []Point, strict bool) error {
	if strict {
		// Check whether all points fit the plane first before adding any
		for _, point := range points {
			if point.X > pln.dimensions[0] || point.X < 0 || point.Y > pln.dimensions[1] || point.Y > 0 {
				return fmt.Errorf("point (%v) is outside of the plane", point)
			}
		}
	}
	// Actually add the points
	for _, point := range points {
		pln.data.Add(point)
	}
	return nil
}

// ErasePoints removes points passed to it in the first argument from the plane's PointSet. It can optionally
// reject points which are not already on the plane if the second argument is `true`.
func (pln *Plane) ErasePoints(points []Point, strict bool) error {
	if strict {
		// Check whether all points are on the plane
		for _, point := range points {
			if !pln.data.CheckFor(point) {
				return fmt.Errorf("point (%v) is outside of the plane", point)
			}
		}
	}
	// Actually remove the points
	for _, point := range points {
		pln.data.Remove(point)
	}
	return nil
}

// ReadPoints reads all points from the plane's PointSet and returns them as an array (technically a slice).
func (pln *Plane) ReadPoints() []Point {
	return pln.data.AsArray()
}

// ReadPointsByFilter accepts a function as an argument, and runs the function with the X and Y coordinates
// of each point on the plane as argument 1 and 2 respectively. For each point that the function returns `true`
// for, that point is returned as part of a slice of points.
func (pln *Plane) ReadPointsByFilter(f filterFunction) []Point {
	allPoints := pln.data.AsArray()
	resultingPoints := []Point{}
	for _, point := range allPoints {
		if f(point) {
			resultingPoints = append(resultingPoints, point)
		}
	}
	return resultingPoints
}

// JoinPoints accepts a slice of points and joins each consecutive point using a line of points. It then returns
// all the points together as a slice.
func (pln *Plane) JoinPoints(points []Point) []Point {
	if len(points) == 0 {
		return []Point{}
	}

	allPoints := PointSet{}

	var pPoint Point
	noPPoint := true

	for _, cPoint := range points {
		if noPPoint {
			noPPoint = false
		} else {
			// Work out difference in X and difference in Y between the two points
			posDiffs := []float64{
				float64(cPoint.X - pPoint.X),
				float64(cPoint.Y - pPoint.Y),
			}
			// Work out length of line between points (round to nearest pixel) (pythagoras)
			posDiffLength := int(
				math.Round(
					math.Sqrt(math.Pow(posDiffs[0], 2) + math.Pow(posDiffs[1], 2)),
				),
			)
			// Work out the gradient of the line between the two points. In case posDiffs[0] == 0
			// and posDiffs [1] != 0 (line is going straight up/down), the gradient is undefined,
			// so we'll just replace it with 0 and the special case will be handled later to draw
			// the correct line
			var posDiffGradient float64
			if posDiffs[0] == 0 && posDiffs[1] != 0 {
				posDiffGradient = 0
			} else {
				posDiffGradient = posDiffs[1] / posDiffs[0]
			}

			// We now have all the information we need to join the two points

			// The coordinates of the point to be generated/drawn. We need to keep track of this
			// as the location of the next point is calculated from this. We start, of course, with
			// the initial point itself
			cDrawPoint := [2]float64{
				float64(pPoint.X),
				float64(pPoint.Y),
			}

			// What we should increment the x coordinate by to generate the next pixel. This is simply the
			// ratio of the diagonal length to the x-only length
			xIncrement := posDiffs[0] / float64(posDiffLength)

			// What we should increment the y coordinate by to generate the next pixel. This is the x increment
			// multiplied by the ratio of x to y for proportion
			yIncrement := xIncrement * posDiffGradient

			// If posDiffs[0] is 0 and posDiffs[1] is not, we are attempting to draw a vertical line which has
			// an undefined gradient (which was replaced with 0 earlier). yIncrement will therefore be 0 even
			// though we do need to increment y to get to the next point. For this reason, depending on whether
			// the next point is above or below, set yIncrement to 1 or -1
			if posDiffs[0] == 0 {
				if posDiffs[1] > 0 {
					yIncrement = 1
				} else if posDiffs[1] < 0 {
					yIncrement = -1
				}
				// No else since if posDiffs[1] also equals 0 we are just drawing a point
			}

			// For each point needing to be generated
			// TODO: Skip points which round to the same coordinates for efficiency
			for i := 0; i < posDiffLength; i++ {
				// Add the point to the set of points
				point := Point{
					X: int(math.Round(cDrawPoint[0])),
					Y: int(math.Round(cDrawPoint[1])),
				}
				//fmt.Println(point)
				allPoints.Add(point)
				// Calculate the position of the next point
				cDrawPoint[0] += xIncrement
				cDrawPoint[1] += yIncrement
			}
		}
		pPoint = cPoint
	}

	return allPoints.AsArray()
}

// JoinAndFillPoints works much like JoinPoints except it also attempts to work out the "inside"
// and "outside" of the shape being created, and then fill it with points. This is a one-size-fits-all
// implementation and works reasonably well for most simple shapes, but it's probably always better to make your own.
func (pln *Plane) JoinAndFillPoints(points []Point) []Point {
	if len(points) == 0 {
		return []Point{}
	}

	// Calculate the points furthest in each direction
	minX, minY, maxX, maxY := pln.dimensions[0], pln.dimensions[1], 0, 0
	for _, point := range points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	// Create a rectangle containing the entire shape
	rect := []Point{
		{
			X: minX - 1,
			Y: minY - 1,
		},
		{
			X: maxX + 1,
			Y: maxY + 1,
		},
	}

	// Work out the points making up the outline of the shape
	points = pln.JoinPoints(points)

	// Function used below to search the array for a particular point
	pointExists := func(searchPoint Point) bool {
		for _, point := range points {
			if point == searchPoint {
				return true
			}
		}
		return false
	}

	// Cast "rays" from one side of the rectangle to the other. If they encounter
	// a point making up the sides of the shape, start drawing points until another
	// point is encountered
	drawnPoints := PointSet{}
	drawnPoints.AddArray(points) // The outline should always be drawn
	for x := rect[0].X; x < rect[1].X; x++ {
		totalRayPoints := []Point{}
		for y := rect[0].Y; y < rect[1].Y; y++ {
			curPoint := Point{X: x, Y: y}
			nextPoint := Point{X: x, Y: y + 1}
			// If the current point exists, count it (unless the next one
			// exists too in which case count them as one on the next loop,
			// as a thick line should still count as the side of the shape).
			if pointExists(curPoint) && !pointExists(nextPoint) {
				totalRayPoints = append(totalRayPoints, curPoint)
			}
		}
		// If the number of points is odd, discard the last one as not to paint the rest of the rectangle
		if len(totalRayPoints)%2 != 0 {
			totalRayPoints = totalRayPoints[:len(totalRayPoints)-1]
		}
		// Pair up the points and draw a line between each pair as this is the inside of the shape
		for i := range totalRayPoints {
			if i%2 != 0 {
				continue
			}
			drawnPoints.AddArray(pln.JoinPoints([]Point{totalRayPoints[i], totalRayPoints[i+1]}))
		}

	}

	return drawnPoints.AsArray()
}

// Flip flips the plane along the X or the Y axis. Sometimes you may need to flip the points on the plane,
// for example when converting to an image where the Y (1) axis it flipped. The axis of the flip is given as
// an integer - 0 for X and 1 for Y.
func (pln *Plane) Flip(dimension int) {
	dimensionMax := pln.dimensions[dimension]
	flipped := PointSet{}
	for _, point := range pln.data.AsArray() {
		pln.data.Remove(point)
		if dimension == 0 {
			point.X = dimensionMax - point.X
		} else if dimension == 1 {
			point.Y = dimensionMax - point.Y
		}
		flipped.Add(point)
	}
	pln.data.AddArray(flipped.AsArray())
}

// FlipPoints is much the same as Flip, but acts on a given set of points as opposed to the whole plane.
// It also returns the flipped points as a slice.
func (pln *Plane) FlipPoints(points []Point, dimension int) []Point {
	dimensionMax := pln.dimensions[dimension]
	for i, point := range points {
		if dimension == 0 {
			point.X = dimensionMax - point.X
		} else if dimension == 1 {
			point.Y = dimensionMax - point.Y
		}
		points[i] = point
	}
	return points
}

// GetDimensions returns the size of the plane as an array of two integers
func (pln *Plane) GetDimensions() [2]int {
	return pln.dimensions
}

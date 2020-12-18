package gopoints

import "math"

type Plane struct {
	data       PointSet
	dimensions [2]int
}

func (pln *Plane) Init(dimensions [2]int) {
	pln.dimensions = dimensions
	pln.data = PointSet{}
}

func (pln *Plane) JoinPoints(points []Point) []Point {
	// Go does not provide sets, so to avoid duplicate points due to rounding, use map keys
	// at little extra cost as structs take 0 bytes
	allPoints := PointSet{}

	var pPoint Point
	noPPoint := true

	for _, cPoint := range points {
		if noPPoint {
			noPPoint = false
		} else {
			// Work out difference in X and difference in Y between the two points
			posDiffs := []float64{
				float64(cPoint.Pos[0] - pPoint.Pos[0]),
				float64(cPoint.Pos[1] - pPoint.Pos[1]),
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
				float64(pPoint.Pos[0]),
				float64(pPoint.Pos[1]),
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
			// TODO: Skip points which round to the same coordinated for efficiency
			for i := 0; i < posDiffLength; i++ {
				// Add the point to the set of points
				point := Point{
					Pos: [2]int{
						int(math.Round(cDrawPoint[0])),
						int(math.Round(cDrawPoint[1])),
					},
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

	return allPoints.ToArray()
}

func (pln *Plane) Flip(dimension int) {
	// Sometimes you may need to flip the points on the plane, for example when converting to
	// an image where the Y (1) axis it flipped
	for _, point := range pln.data.ToArray() {
		point.Pos[dimension] = pln.dimensions[dimension] - point.Pos[dimension]
	}
}

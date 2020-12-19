# GoPoints

A very simple Go library for creating shapes on a 2D plane out of points, using only the Go stdlib. Supports some basic functions like:
- storing points in a `plane` structure
- joining arbitrary points with lines (including diagonal)
- filling shapes with points
- flipping the whole plane or some points along the X or Y axis
- fetching all points from a plane
- fetching points from a plane based on the return value (boolean) of a given function (like a mathematical equation)
- *TODO*: Applying an arbitrary transformation based on a function to the points on the plane

## Use cases

Can be used for any number of things, but my primary use-case is generating images made of 2D, 1D or 0D shapes, such as the [example image](https://github.com/TR-SLimey/gopoints/blob/master/examples/CreatePngImage.png?raw=true) generated with `examples/CreatePngImage.go`:

![example image](https://github.com/TR-SLimey/gopoints/blob/master/examples/CreatePngImage.png?raw=true)

## Usage examples

Some example code can be found in the `examples` directory in the root of this repository


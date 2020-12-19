package gopoints

// A simple set implementation for 2D points as Go does not provide one. Uses map keys as the storage
// as they do not repeat and are unordered. The value used is a blank struct as it causes little
// to no overead due to its 0-byte size.
type PointSet struct {
	// The actual data-storage part
	data map[Point]struct{}
}

// Init initialises the set by creating the map which stores the data. This is called implicitly
// by all other methods of the set so does not need to be called explicitly. It does not overwrite
// the set if called multiple times.
func (set *PointSet) Init() {
	// Only if the set is not already initialised in order to avoid accidentally removing data
	if set.data == nil {
		set.data = make(map[Point]struct{})
	}
}

// Add adds a single point into the set (unless it already exists of course - it's a set after all)
func (set *PointSet) Add(point Point) {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	set.data[point] = struct{}{}
}

// AddArray works much like Add but accepts a slice of points (poor naming) and adds each element of it to the set
func (set *PointSet) AddArray(pointArray []Point) {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	for _, point := range pointArray {
		set.Add(point)
	}
}

// Remove removes a single point from the set
func (set *PointSet) Remove(point Point) {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	delete(set.data, point)
}

// RemoveArray works much like Remove but accepts a slice of points (poor naming) and removes each element of it from the set
func (set *PointSet) RemoveArray(pointArray []Point) {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	for _, point := range pointArray {
		set.Remove(point)
	}
}

// CheckFor checks if a given point is in the set
func (set PointSet) CheckFor(point Point) bool {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	_, ok := set.data[point]
	return ok
}

// CheckForAll checks if all points in the given slice are present in the set and returns a single boolean to indicate the result
func (set *PointSet) CheckForAll(points []Point) bool {
	for _, point := range points {
		if !set.CheckFor(point) {
			return false
		}
	}
	return true
}

// AsArray returns the contents of the set as a slice
func (set PointSet) AsArray() []Point {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	// Convert the set to an array for easy access
	arr := []Point{}
	for point := range set.data {
		arr = append(arr, point)
	}
	return arr
}

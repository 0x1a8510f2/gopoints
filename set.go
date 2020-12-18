package ShapeCreator

// A simple set implementation as Go does not provide one. Uses map keys as the storage medium
// as they do not repeat and are unordered. The value used is a blank struct as it causes little
// to no overead due to its 0-byte size.

type set struct {
	// The actual data-storage part
	data map[interface{}]struct{}
}

func (set *set) Init() {
	// Only if the set is not already initialised in order to avoid accidentally removing data
	if set.data == nil {
		set.data = make(map[interface{}]struct{})
	}
}

func (set *set) Add(element interface{}) {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	set.data[element] = struct{}{}
}

func (set *set) Remove(element interface{}) {
	// Automatically init the set if it's not initialised yet (check handled by Init())
	set.Init()

	delete(set.data, element)
}

func (set *set) ToArray() []interface{} {
	// Convert the set to an array for easy access
	arr := []interface{}{}
	for element, _ := range set.data {
		arr = append(arr, element)
	}
	return arr
}

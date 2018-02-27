package georadis

//Distance unit
const (
	M  = "m"
	KM = "km"
	Mi = "mi"
	Ft = "ft"
)

// Option is the georedius option
type Option = int

const (
	// WithHash returns the value with geohash
	WithHash Option = iota
	// WithDist returns the distance between the location to neighbor
	WithDist
	// WithCoord returns the coordinate of nieghbor
	WithCoord
)

// index for redigo returned data
const (
	lonIdx = iota
	latIdx
)

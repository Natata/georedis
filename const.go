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
	// WithDist returns the distance between the location to neighbor
	WithDist Option = iota
	// WithHash returns the value with geohash
	WithHash
	// WithCoord returns the coordinate of nieghbor
	WithCoord
)

// index for redigo returned data
const (
	lonIdx = iota
	latIdx
)

var optMap = map[Option]string{
	WithDist:  "WITHDIST",
	WithHash:  "WITHHASH",
	WithCoord: "WITHCOORD",
}

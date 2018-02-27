package georadis

// MetaData present the Coordinate and the key of data
type MetaData struct {
	// Coordinate of data
	Coord Coordinate
	// Data key
	DKey string
}

// NewMetaData create a meta dta
func NewMetaData(dKey string, lat, lon float64) *MetaData {
	return &MetaData{
		Coord: Coordinate{
			Lat: lat,
			Lon: lon,
		},
		DKey: dKey,
	}
}

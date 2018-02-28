package georedis

// Member present the Coordinate and the key of data
type Member struct {
	Name  string     `json:"name"`
	Coord Coordinate `json:"coord"`
}

// NewMember create a meta dta
func NewMember(name string, lat, lon float64) *Member {
	return &Member{
		Name: name,
		Coord: Coordinate{
			Lat: lat,
			Lon: lon,
		},
	}
}

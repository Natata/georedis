package georadis

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// NeighborData present the data of a neighbor
type NeighborData struct {
	Name  string
	Dist  float64
	Coord Coordinate
	Hash  int64
}

// NewNeighborData transfers the raw value from GEOREDIUS to NeighborData
func NewNeighborData(raw reflect.Value, opts ...Option) (*NeighborData, error) {

	// no option, only slice of string
	if len(opts) == 0 {
		name, err := toString(unpackValue(raw))
		log.Printf("no options, name: %v", name)
		return &NeighborData{Name: name}, err
	}

	if raw.Kind() != reflect.Slice {
		return nil, fmt.Errorf("new neighbor data fail: %v", raw.Kind())
	}

	nd := &NeighborData{}

	transformers := []func(v reflect.Value) error{
		func(v reflect.Value) error { // name
			fmt.Println("parse name")
			name, err := toString(unpackValue(v))
			if err != nil {
				return err
			}
			fmt.Println("name: ", name)
			nd.Name = name
			return nil
		},
		func(v reflect.Value) error { // distance
			fmt.Println("parse distance")
			dist, err := toFloat64(unpackValue(v))
			if err != nil {
				return err
			}
			fmt.Println("dist: ", dist)
			nd.Dist = dist
			return nil
		},
		func(v reflect.Value) error { // hash (int)
			fmt.Println("parse hash")
			hash := toInt64(unpackValue(v))
			fmt.Println("hash: ", hash)
			nd.Hash = hash
			return nil
		},
		func(v reflect.Value) error { // coordinate
			fmt.Println("parse coordinate")
			fmt.Println(v)
			coord, err := toCoordinate(unpackValue(v))
			if err != nil {
				return err
			}
			fmt.Println("coord lat: ", coord.Lat)
			fmt.Println("coord lon: ", coord.Lon)
			nd.Coord = coord
			return nil
		},
	}

	ckTab := make([]bool, 4)
	ckTab[0] = true
	for _, opt := range opts {
		fmt.Println("has opt: ", optMap[opt])
		ckTab[opt+1] = true
	}
	i := 0
	for fi, ok := range ckTab {
		if !ok {
			continue
		}
		err := transformers[fi](raw.Index(i))
		if err != nil {
			return nil, err
		}
		i++
	}

	fmt.Println("Done")
	return nd, nil
}

func toString(v reflect.Value) (string, error) {
	if v.Kind() != reflect.Slice {
		return "", fmt.Errorf("to string fail: %v", v.Kind())
	}

	b := v.Bytes()
	return string(b), nil
}

func toFloat64(v reflect.Value) (float64, error) {
	s, err := toString(v)
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

func toInt64(v reflect.Value) int64 {
	i := v.Int()
	return i
}

func toCoordinate(v reflect.Value) (Coordinate, error) {
	if v.Kind() != reflect.Slice || v.Len() != 2 {
		return Coordinate{}, fmt.Errorf("invalid data format for coordainate, %v", v)
	}

	var coord Coordinate
	var err error
	coord.Lon, err = toFloat64(unpackValue(v.Index(lonIdx)))
	if err != nil {
		return coord, err
	}

	coord.Lat, err = toFloat64(unpackValue(v.Index(latIdx)))
	if err != nil {
		return coord, err
	}

	return coord, nil
}

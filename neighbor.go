package georedis

import (
	"fmt"
	"log"
	"reflect"
)

// Neighbor is member with the distance and geohash
type Neighbor struct {
	Member
	Dist float64
	Hash int64
}

// NewNeighbor transfers the raw value from GEOREDIUS to Member
func NewNeighbor(raw reflect.Value, opts ...Option) (*Neighbor, error) {

	// no option, only slice of string
	if len(opts) == 0 {
		name, err := toString(unpackValue(raw))
		log.Printf("no options, name: %v", name)
		return &Neighbor{Member: Member{Name: name}}, err
	}

	if raw.Kind() != reflect.Slice {
		return nil, fmt.Errorf("new neighbor data fail: %v", raw.Kind())
	}

	nb := &Neighbor{}

	transformers := []func(v reflect.Value) error{
		func(v reflect.Value) error { // name
			fmt.Println("parse name")
			name, err := toString(unpackValue(v))
			if err != nil {
				return err
			}
			fmt.Println("name: ", name)
			nb.Name = name
			return nil
		},
		func(v reflect.Value) error { // distance
			fmt.Println("parse distance")
			dist, err := toFloat64(unpackValue(v))
			if err != nil {
				return err
			}
			fmt.Println("dist: ", dist)
			nb.Dist = dist
			return nil
		},
		func(v reflect.Value) error { // hash (int)
			fmt.Println("parse hash")
			hash := toInt64(unpackValue(v))
			fmt.Println("hash: ", hash)
			nb.Hash = hash
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
			nb.Coord = coord
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
	return nb, nil
}

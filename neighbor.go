package georedis

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
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
		return &Neighbor{Member: Member{Name: name}}, err
	}

	if raw.Kind() != reflect.Slice {
		log.WithFields(log.Fields{
			"type":     raw.Kind(),
			"raw data": raw,
		}).Error("wrong type, want slice")
		return nil, fmt.Errorf("new neighbor data fail: %v", raw.Kind())
	}

	nb := &Neighbor{}

	transformers := []func(v reflect.Value) error{
		func(v reflect.Value) error { // name
			name, err := toString(unpackValue(v))
			log.Debugf("convert to name, raw: %v, error: %v", v, err)
			if err != nil {
				return err
			}
			nb.Name = name
			return nil
		},
		func(v reflect.Value) error { // distance
			dist, err := toFloat64(unpackValue(v))
			log.Debugf("convert to distance, raw: %v, error: %v", v, err)
			if err != nil {
				return err
			}
			nb.Dist = dist
			return nil
		},
		func(v reflect.Value) error { // hash (int)

			// TODO: handle panic when convert fail

			hash := toInt64(unpackValue(v))
			log.Debugf("convert to distance, raw: %v", v)
			nb.Hash = hash
			return nil
		},
		func(v reflect.Value) error { // coordinate
			coord, err := toCoordinate(unpackValue(v))
			log.Debugf("convert to coordinate, raw: %v, error: %v", v, err)
			if err != nil {
				return err
			}
			nb.Coord = coord
			return nil
		},
	}

	ckTab := make([]bool, 4)
	ckTab[0] = true
	for _, opt := range opts {
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

	return nb, nil
}

package georedis

import (
	"fmt"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

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

func rawToNeighbors(r interface{}, options ...Option) ([]*Neighbor, error) {
	v := reflect.ValueOf(r)

	if v.Kind() != reflect.Slice {
		log.WithFields(log.Fields{
			"type": v.Kind(),
		}).Error("wrong type, want slice")
		return nil, fmt.Errorf("wrong type: %v", v.Kind())
	}

	results := make([]*Neighbor, v.Len())
	var err error
	for i := 0; i < v.Len(); i++ {
		results[i], err = NewNeighbor(unpackValue(v.Index(i)), options...)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"raw data": v.Index(i),
			}).Error("fail to convert raw data to Neighbor type")
			return nil, err
		}
		log.WithFields(log.Fields{
			"name":      results[i].Name,
			"latitude":  results[i].Coord.Lat,
			"longitude": results[i].Coord.Lon,
			"distance":  results[i].Dist,
			"hash":      results[i].Hash,
		}).Info("find neighbor")
	}

	return results, nil
}

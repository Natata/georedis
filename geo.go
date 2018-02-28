package georadis

import (
	"fmt"
	"log"
	"reflect"

	"github.com/garyburd/redigo/redis"
)

// Geo is the core service for geolocation-related operation
type Geo struct {
	pool *redis.Pool
}

// NewGeo creates a Geo service
func NewGeo(pool *redis.Pool) *Geo {
	return &Geo{
		pool: pool,
	}
}

// Set sets key and related meta data to redis
func (s *Geo) Set(key string, data []*MetaData) error {

	// TODO: check lat, lon validate. For now, redis would check the coordinate itself

	conn := s.pool.Get()
	defer conn.Close()

	for _, d := range data {
		r, err := conn.Do("GEOADD", key, d.Coord.Lon, d.Coord.Lat, d.DKey)
		log.Printf("add key %v, data: %v, r: %v, err: %v", key, d, r, err)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get gets the meta data by key
// returned meta data hase the same order of dKeys
// leave nil for the keys have no data
func (s *Geo) Get(key string, dKeys []string) ([]*MetaData, error) {
	conn := s.pool.Get()
	defer conn.Close()

	// get data from redis
	args := []interface{}{key}
	for i := range dKeys {
		args = append(args, dKeys[i])
	}
	r, err := redis.Positions(conn.Do("GEOPOS", args...))
	if err != nil {
		return nil, err
	}
	log.Printf("GEOPOS result: %v", r)

	// create meta data
	data := make([]*MetaData, len(r))
	for i := range r {
		if r[i] == nil {
			log.Printf("no data for %v", dKeys[i])
		} else {
			data[i] = NewMetaData(dKeys[i], r[i][lonIdx], r[i][latIdx])
		}
	}

	return data, nil
}

// Neighbors find the neighbor
func (s *Geo) Neighbors(key string, Coord Coordinate, radius int, unit string, options ...Option) ([]*NeighborData, error) {
	conn := s.pool.Get()
	defer conn.Close()

	// basic command
	args := []interface{}{key, Coord.Lon, Coord.Lat, radius, unit}

	// set options
	for _, opt := range options {
		args = append(args, optMap[opt])
	}
	fmt.Println("args: ", args)

	// execute command
	r, err := conn.Do("GEORADIUS", args...)
	if err != nil {
		return nil, err
	}
	log.Printf("GEORADIUS result: %v", r)

	return rawToNeighbors(r, options...)
}

// Dist cc
func (s *Geo) Dist(key, a, b string, u Unit) (float64, error) {
	conn := s.pool.Get()
	defer conn.Close()

	r, err := conn.Do("GEODIST", key, a, b, u)
	if err != nil {
		return 0, err
	}
	log.Printf("GEODIST result: %v", r)

	v := reflect.ValueOf(r)
	f, err := toFloat64(v)
	if err != nil {
		return 0, err
	}

	return f, nil
}

// GeoPos cc
func (s *Geo) GeoPos(key string, list ...string) ([]Coordinate, error) {
	conn := s.pool.Get()
	defer conn.Close()

	args := []interface{}{key}
	for _, l := range list {
		args = append(args, l)
	}
	r, err := conn.Do("GEOPOS", args...)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(r)
	coords := make([]Coordinate, len(list))
	for i := 0; i < v.Len(); i++ {
		pos := unpackValue(v.Index(i))
		coord, err := toCoordinate(pos)
		if err != nil {
			return nil, err
		}
		coords[i] = coord
	}
	return coords, nil
}

// GeoHash return the geohash of place
func (s *Geo) GeoHash(key string, list ...string) ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	args := []interface{}{key}
	for _, l := range list {
		args = append(args, l)
	}
	r, err := conn.Do("GEOHASH", args...)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(r)
	hashs := make([]string, len(list))
	for i := 0; i < v.Len(); i++ {
		hashv := unpackValue(v.Index(i))
		hash, err := toString(hashv)
		if err != nil {
			return nil, err
		}
		hashs[i] = hash
	}
	return hashs, nil
}

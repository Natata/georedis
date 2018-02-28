package georadis

import (
	"fmt"
	"log"

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

// Distance cc
func (s *Geo) Distance() {}

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

// Add adds key and related meta data to redis
func (s *Geo) Add(key string, data []*Member) error {

	// TODO: check lat, lon validate. For now, redis would check the coordinate itself

	conn := s.pool.Get()
	defer conn.Close()

	for _, d := range data {
		r, err := conn.Do("GEOADD", key, d.Coord.Lon, d.Coord.Lat, d.Name)
		log.Printf("add data: %v, r: %v, err: %v", d, r, err)
		if err != nil {
			return err
		}
	}

	return nil
}

// Pos gets the meta data by key
// returned meta data hase the same order of names
// leave nil for the keys have no data
func (s *Geo) Pos(key string, names []string) ([]*Member, error) {
	conn := s.pool.Get()
	defer conn.Close()

	// get data from redis
	args := []interface{}{key}
	for i := range names {
		args = append(args, names[i])
	}
	r, err := redis.Positions(conn.Do("GEOPOS", args...))
	if err != nil {
		return nil, err
	}
	log.Printf("GEOPOS result: %v", r)

	// create meta data
	data := make([]*Member, len(r))
	for i := range r {
		if r[i] == nil {
			log.Printf("no data for %v", names[i])
		} else {
			data[i] = NewMember(key, names[i], r[i][lonIdx], r[i][latIdx])
		}
	}

	return data, nil
}

// RadiusByName find nearby members of member
func (s *Geo) RadiusByName(key string, name string, radius int, unit string, options ...Option) ([]*Neighbor, error) {
	mems, err := s.Pos(key, []string{name})
	if err != nil {
		return nil, err
	}

	if len(mems) != 1 {
		return nil, fmt.Errorf("%v not exist in the key %v", name, key)
	}

	ns, err := s.Radius(key, mems[0].Coord, radius, unit, options...)
	if err != nil {
		return nil, err
	}

	return ns, nil
}

// Radius find the neighbor with coordinate
func (s *Geo) Radius(key string, Coord Coordinate, radius int, unit string, options ...Option) ([]*Neighbor, error) {
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

// Hash return the geohash of place
func (s *Geo) Hash(key string, list ...string) ([]string, error) {
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

package georedis

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
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
	conn := s.pool.Get()
	defer conn.Close()

	for _, d := range data {
		r, err := conn.Do("GEOADD", key, d.Coord.Lon, d.Coord.Lat, d.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err,
				"key":       key,
				"name":      d.Name,
				"latitude":  d.Coord.Lat,
				"longitude": d.Coord.Lon,
			}).Error("add failed")
			return err
		}
		log.WithFields(log.Fields{
			"result":    r,
			"key":       key,
			"name":      d.Name,
			"latitude":  d.Coord.Lat,
			"longitude": d.Coord.Lon,
		}).Info("add success")
	}

	return nil
}

// Pos gets the meta data by key
// returned meta data hase the same order of names
// leave nil for the keys have no data
func (s *Geo) Pos(key string, names ...string) ([]*Member, error) {
	conn := s.pool.Get()
	defer conn.Close()

	// get data from redis
	args := []interface{}{key}
	for i := range names {
		args = append(args, names[i])
	}
	r, err := redis.Positions(conn.Do("GEOPOS", args...))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
			"names": names,
		}).Error("get positions failed")
		return nil, err
	}

	// create meta data
	data := make([]*Member, len(r))
	for i := range r {
		if r[i] == nil {
			log.WithFields(log.Fields{
				"name": names[i],
			}).Info("no position data")
		} else {
			data[i] = NewMember(names[i], r[i][latIdx], r[i][lonIdx])
			log.WithFields(log.Fields{
				"name":      names[i],
				"latitude":  r[i][latIdx],
				"longitude": r[i][lonIdx],
			}).Info("get position data")
		}
	}

	return data, nil
}

// RadiusByName find nearby members of member
// the result include the name itself
func (s *Geo) RadiusByName(key string, name string, radius int, unit string, options ...Option) ([]*Neighbor, error) {
	mems, err := s.Pos(key, name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
			"name":  name,
		}).Error("no position data of this point")
		return nil, err
	}
	if len(mems) != 1 {
		log.WithFields(log.Fields{
			"key":     key,
			"name":    name,
			"members": mems,
		}).Error("wrong data number")
		return nil, fmt.Errorf("have multiple or zero results, key: %v, name: %v, members: %v", key, name, mems)
	}

	return s.Radius(key, mems[0].Coord, radius, unit, options...)
}

// Radius find the neighbor with coordinate
func (s *Geo) Radius(key string, coord Coordinate, radius int, unit string, options ...Option) ([]*Neighbor, error) {
	conn := s.pool.Get()
	defer conn.Close()

	// basic command
	args := []interface{}{key, coord.Lon, coord.Lat, radius, unit}

	// set options
	for _, opt := range options {
		args = append(args, optMap[opt])
	}

	// execute command
	r, err := conn.Do("GEORADIUS", args...)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"key":       key,
			"latitude":  coord.Lat,
			"longitude": coord.Lon,
		}).Error("get radius fail")
		return nil, err
	}

	return rawToNeighbors(r, options...)
}

// Dist cc
func (s *Geo) Dist(key, a, b string, u Unit) (float64, error) {
	conn := s.pool.Get()
	defer conn.Close()

	r, err := conn.Do("GEODIST", key, a, b, u)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"key":     key,
			"place a": a,
			"place b": b,
			"unit":    u,
		}).Error("get dist fail")
		return 0, err
	}

	v := reflect.ValueOf(r)
	f, err := toFloat64(v)
	if err != nil {
		log.Errorf("convert value to float64 failed, value: %v, error: %v", v, err)
		return 0, err
	}

	log.WithFields(log.Fields{
		"key":      key,
		"place a":  a,
		"place b":  b,
		"distance": f,
		"unit":     u,
	}).Info("get distance")

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
		log.WithFields(log.Fields{
			"key":  key,
			"list": list,
		}).Error("get hash failed")
		return nil, err
	}
	v := reflect.ValueOf(r)
	hashs := make([]string, len(list))
	for i := 0; i < v.Len(); i++ {
		hashv := unpackValue(v.Index(i))
		hash, err := toString(hashv)
		if err != nil {
			log.Errorf("convert to string failed, value: %v, error: %v", v.Index(i), err)
			return nil, err
		}
		hashs[i] = hash
	}
	log.WithFields(log.Fields{
		"key":    key,
		"places": list,
		"hashs":  hashs,
	}).Info("get geohash")
	return hashs, nil
}

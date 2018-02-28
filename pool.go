package georedis

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/garyburd/redigo/redis"
)

// PoolConfig is the struct for config pool
type PoolConfig struct {
	IdleConn     int    `json:"idle_conn"`
	ActiveConn   int    `json:"active_conn"`
	Protocol     string `json:"protocol"`
	Addr         string `json:"addr"`
	DB           int    `json:"db"`
	TobTimeout   string `json:"tob_timeout"`
	IdleTimeout  string `json:"idle_timeout"`
	ConnTimeout  string `json:"conn_timeout"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

// NewPool creates a redis connection pool
func NewPool(configPath string) (*redis.Pool, error) {
	r, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &PoolConfig{}
	err = json.Unmarshal(r, cfg)
	if err != nil {
		return nil, err
	}

	tobT, err := time.ParseDuration(cfg.TobTimeout)
	if err != nil {
		return nil, err
	}
	idleT, err := time.ParseDuration(cfg.IdleTimeout)
	if err != nil {
		return nil, err
	}
	connT, err := time.ParseDuration(cfg.ConnTimeout)
	if err != nil {
		return nil, err
	}
	readT, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		return nil, err
	}
	writeT, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		return nil, err
	}
	return &redis.Pool{
		MaxIdle:   cfg.IdleConn,
		MaxActive: cfg.ActiveConn,
		Dial: func() (redis.Conn, error) {
			connOption := redis.DialConnectTimeout(connT)
			readOption := redis.DialReadTimeout(readT)
			writeOption := redis.DialWriteTimeout(writeT)
			dbOption := redis.DialDatabase(cfg.DB)
			c, err := redis.Dial(
				cfg.Protocol,
				cfg.Addr,
				connOption,
				readOption,
				writeOption,
				dbOption,
			)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < tobT {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
		IdleTimeout: idleT,
		Wait:        true,
	}, nil
}

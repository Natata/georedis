package geochat

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// TODO: use some config package to setup the pool

const (
	idleConn     = 5
	activeConn   = 10
	protocol     = "tcp"
	addr         = ":6379"
	tobTimeout   = time.Minute
	idleTimeout  = 300 * time.Second
	connTimeout  = 6000 * time.Millisecond
	readTimeout  = 6000 * time.Millisecond
	writeTimeout = 6000 * time.Millisecond
	db           = 0
)

// NewPool creates a redis connection pool
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   idleConn,
		MaxActive: activeConn,
		Dial: func() (redis.Conn, error) {
			connOption := redis.DialConnectTimeout(connTimeout)
			readOption := redis.DialReadTimeout(readTimeout)
			writeOption := redis.DialWriteTimeout(writeTimeout)
			dbOption := redis.DialDatabase(db)
			c, err := redis.Dial(protocol, addr, connOption, readOption, writeOption, dbOption)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < tobTimeout {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
		IdleTimeout: idleTimeout,
		Wait:        true,
	}
}

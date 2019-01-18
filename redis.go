package cachita

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"time"

	rds "github.com/gomodule/redigo/redis"
)

var rCache Cache

type redis struct {
	pool   *rds.Pool
	prefix string
	ttl    time.Duration
}

func Redis(addr string) (Cache, error) {
	if rCache == nil {
		var err error
		rCache, err = NewRedisCache(24*time.Hour, 10, "cachita", addr)
		if err != nil {
			return nil, err
		}
	}
	return rCache, nil
}

func NewRedisCache(ttl time.Duration, poolSize int, prefix, addr string) (Cache, error) {
	pool := &rds.Pool{

		IdleTimeout: 240 * time.Second,
		MaxIdle:     poolSize,
		MaxActive:   poolSize,
		Dial: func() (rds.Conn, error) {
			c, err := rds.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c rds.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	c := &redis{
		pool:   pool,
		prefix: prefix,
		ttl:    ttl,
	}

	return c, nil
}

func (c *redis) Get(key string, i interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	s, err := conn.Do("GET", c.k(key))
	if err != nil {
		return err
	}
	if s == nil {
		return ErrNotFound
	}
	data := s.([]byte)
	if data == nil {
		return ErrNotFound
	}
	return msgpack.Unmarshal(data, i)
}

func (c *redis) Put(key string, i interface{}, ttl time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", c.k(key), calculateTtl(ttl, c.ttl).Seconds(), data)
	return err

}

func (c *redis) Incr(key string, ttl time.Duration) error {
	//k := c.k(key)
	//pool, err := c.SingleRedis(c.ctx)
	//if err != nil {
	//	return err
	//}
	//defer pool.Close()
	//sync := rds.SyncCtx{pool}
	//_, err = sync.SendTransaction(c.ctx, []rds.Request{
	//	rds.Req("INCR", k),
	//	rds.Req("EXPIRE", k, calculateTtl(ttl, c.ttl).Seconds()),
	//})

	return nil
}

func (c *redis) Invalidate(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", c.k(key))
	return err
}

func (c *redis) Exists(key string) bool {
	conn := c.pool.Get()
	defer conn.Close()
	b, _ := rds.Bool(conn.Do("EXISTS", c.k(key)))
	return b
}

func (c *redis) k(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

func isInt(i interface{}) bool {
	switch i.(type) {
	case *int:
	case *int8:
	case *int16:
	case *int32:
	case *int64:
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	default:
		return false
	}
	return true
}

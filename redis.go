package cachita

import (
	"fmt"
	"time"

	rds "github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

var rCache Cache

type redis struct {
	c      *rds.Client
	prefix string
	ttl    time.Duration
}

func Redis(addr string) (Cache, error) {
	if rCache == nil {
		c := rds.NewClient(&rds.Options{
			Addr:            addr,
			MaxRetries:      10,
			MaxRetryBackoff: 5 * time.Second,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			PoolTimeout:     10 * time.Second,
			PoolSize:        10,
		})
		_, err := c.Ping().Result()
		if err != nil {
			return nil, err
		}
		rCache = NewRedisCache(24*time.Hour, c, "cachita")
	}
	return rCache, nil
}

func NewRedisCache(ttl time.Duration, c *rds.Client, prefix string) (Cache) {
	rc := &redis{
		c:      c,
		prefix: prefix,
		ttl:    ttl,
	}

	return rc
}

func (c *redis) Get(key string, i interface{}) error {
	isInt := isInt(i)
	var err error
	if isInt {
		i, err = c.c.Get(c.k(key)).Int64()
		if err != nil {
			if err == rds.Nil{
				return ErrNotFound
			}
			return err
		}

		return nil
	}
	data, err := c.c.Get(c.k(key)).Bytes()
	if err != nil {
		if err == rds.Nil{
			return ErrNotFound
		}
		return err
	}
	if data == nil {
		return ErrNotFound
	}
	return msgpack.Unmarshal(data, i)
}

func (c *redis) Put(key string, i interface{}, ttl time.Duration) error {
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}
	return c.c.Set(c.k(key), data, calculateTtl(ttl, c.ttl)).Err()
}

func (c *redis) Incr(key string, ttl time.Duration) error {
	k := c.k(key)
	pipe := c.c.Pipeline()
	pipe.Incr(k)
	pipe.Expire(k, calculateTtl(ttl, c.ttl))
	_, err := pipe.Exec()
	return err
}

func (c *redis) Invalidate(key string) error {
	_, err := c.c.Del(c.k(key)).Result()
	return err
}

func (c *redis) Exists(key string) bool {
	return c.c.Exists(c.k(key)).Val() != 0
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

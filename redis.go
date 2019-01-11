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

func (rc *redis) Get(key string, i interface{}) error {
	data, err := rc.c.Get(rc.k(key)).Bytes()
	if err != nil {
		if err == rds.Nil {
			return ErrNotFound
		}
		return err
	}
	return msgpack.Unmarshal(data, i)
}

func (rc *redis) Put(key string, i interface{}, ttl time.Duration) error {
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}
	return rc.c.Set(rc.k(key), data, calculateTtl(ttl, rc.ttl)).Err()
}

func (rc *redis) Invalidate(key string) error {
	rc.c.Del(rc.k(key))
	return nil
}

func (rc *redis) Exists(key string) bool {
	return rc.c.Exists(rc.k(key)).Val() != 0
}

func (rc *redis) k(key string) string {
	return fmt.Sprintf("%s:%s", rc.prefix, key)
}

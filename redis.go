package cachita

import (
	"fmt"
	"time"

	"github.com/mediocregopher/radix"
	"github.com/vmihailenco/msgpack"
)

var rCache Cache

type redis struct {
	pool   *radix.Pool
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
	pool, err := radix.NewPool("tcp", addr, poolSize)
	if err != nil {
		return nil, err
	}

	rc := &redis{
		pool:   pool,
		prefix: prefix,
		ttl:    ttl,
	}

	return rc, nil
}

func (rc *redis) Get(key string, i interface{}) error {
	var data []byte
	err := rc.pool.Do(radix.FlatCmd(&data, "GET", rc.k(key)))
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNotFound
	}
	return msgpack.Unmarshal(data, i)
}

func (rc *redis) Put(key string, i interface{}, ttl time.Duration) error {
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}
	return rc.pool.Do(radix.FlatCmd(nil, "SETEX", rc.k(key), calculateTtl(ttl, rc.ttl).Seconds(), data))
}

func (rc *redis) Invalidate(key string) error {
	return rc.pool.Do(radix.FlatCmd(nil, "DEL", rc.k(key)))
}

func (rc *redis) Exists(key string) bool {
	var b bool
	rc.pool.Do(radix.FlatCmd(&b, "EXISTS", rc.k(key)))
	return b
}

func (rc *redis) k(key string) string {
	return fmt.Sprintf("%s:%s", rc.prefix, key)
}

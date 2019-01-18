package cachita

import (
	"context"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"time"

	rds "github.com/joomcode/redispipe/redis"
	"github.com/joomcode/redispipe/redisconn"
)

var rCache Cache

type redis struct {
	SingleRedis func(ctx context.Context) (rds.Sender, error)
	pool        rds.SyncCtx
	ctx         context.Context
	prefix      string
	ttl         time.Duration
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
	ctx := context.Background()
	SingleRedis := func(ctx context.Context) (rds.Sender, error) {
		opts := redisconn.Opts{
			DB: 1,
			Logger:   redisconn.NoopLogger{},
		}
		conn, err := redisconn.Connect(ctx, addr, opts)
		return conn, err
	}


	c := &redis{
		SingleRedis: SingleRedis,
		//pool:   sync,
		ctx:    ctx,
		prefix: prefix,
		ttl:    ttl,
	}

	return c, nil
}

func (c *redis) Get(key string, i interface{}) error {
	pool, err := c.SingleRedis(c.ctx)
	if err != nil {
		return err
	}
	defer pool.Close()
	sync := rds.SyncCtx{pool}

	s := sync.Do(c.ctx, "GET", c.k(key))
	err = rds.AsError(s)
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
	pool, err := c.SingleRedis(c.ctx)
	if err != nil {
		return err
	}
	defer pool.Close()
	sync := rds.SyncCtx{pool}
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}

	res := sync.Do(c.ctx, "SETEX", c.k(key), calculateTtl(ttl, c.ttl).Seconds(), data)
	return rds.AsError(res)

}

func (c *redis) Incr(key string, ttl time.Duration) error {
	k := c.k(key)
	pool, err := c.SingleRedis(c.ctx)
	if err != nil {
		return err
	}
	defer pool.Close()
	sync := rds.SyncCtx{pool}
	_, err = sync.SendTransaction(c.ctx, []rds.Request{
		rds.Req("INCR", k),
		rds.Req("EXPIRE", k, calculateTtl(ttl, c.ttl).Seconds()),
	})

	return err
}

func (c *redis) Invalidate(key string) error {
	pool, err := c.SingleRedis(c.ctx)
	if err != nil {
		return err
	}
	defer pool.Close()
	sync := rds.SyncCtx{pool}
	return rds.AsError(sync.Do(c.ctx, "DEL", c.k(key)))
}

func (c *redis) Exists(key string) bool {
	pool, err := c.SingleRedis(c.ctx)
	if err != nil {
		return false
	}
	defer pool.Close()
	sync := rds.SyncCtx{pool}
	res := sync.Do(c.ctx, "EXISTS", c.k(key))
	if rds.AsError(res) != nil {
		return false
	}
	return res != nil && res != int64(0)
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

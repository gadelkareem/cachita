package cachita

import (
    "fmt"
    "time"

    "github.com/mediocregopher/radix/v3"
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

    c := &redis{
        pool:   pool,
        prefix: prefix,
        ttl:    ttl,
    }

    return c, nil
}

func (c *redis) Get(key string, i interface{}) error {
    s := i
    isInt := isInt(i)
    var data []byte
    if !isInt {
        s = &data
    }
    err := c.pool.Do(radix.FlatCmd(s, "GET", c.k(key)))
    if err != nil {
        return err
    }

    if isInt {
        i = s
        return nil
    }

    data = *s.(*[]byte)
    if data == nil {
        return ErrNotFound
    }
    return msgpack.Unmarshal(data, i)
}

func (c *redis) Put(key string, i interface{}, ttl time.Duration) error {
    s := i
    if !isInt(i) {
        data, err := msgpack.Marshal(i)
        if err != nil {
            return err
        }
        s = &data
    }

    return c.pool.Do(radix.FlatCmd(nil, "SETEX", c.k(key), calculateTtl(ttl, c.ttl).Seconds(), s))
}

func (c *redis) Incr(key string, ttl time.Duration) (int64, error) {
    k := c.k(key)
    incr := radix.NewEvalScript(1, `
		local n = redis.call("incr",KEYS[1])
        if tonumber(n) == 1 then
    		redis.call("expire",KEYS[1],ARGV[1])
		end
		return n
`)
    var n int64
    err := c.pool.Do(incr.Cmd(&n, k, fmt.Sprintf("%.0f", calculateTtl(ttl, c.ttl).Seconds())))
    return n, err
}

func (c *redis) Invalidate(key string) error {
    return c.pool.Do(radix.Cmd(nil, "DEL", c.k(key)))
}

func (c *redis) Exists(key string) bool {
    var b bool
    err := c.pool.Do(radix.Cmd(&b, "EXISTS", c.k(key)))
    return err == nil && b
}

func (c *redis) k(key string) string {
    return fmt.Sprintf("%s:keys::%s", c.prefix, key)
}

func (c *redis) t(tag string) string {
    return fmt.Sprintf("%s:tags::%s", c.prefix, tag)
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

func (c *redis) InvalidateMulti(keys ...string) error {
    var rKeys []string
    for _, k := range keys {
        rKeys = append(rKeys, c.k(k))
    }
    return c.pool.Do(radix.Cmd(nil, "DEL", rKeys...))
}

func (c *redis) Tag(key string, tags ...string) (err error) {
    tags = uniqueTags(tags)
    rKey := c.k(key)
    var cmds []radix.CmdAction
    for _, t := range tags {
        cmds = append(cmds, radix.FlatCmd(nil, "SADD", c.t(t), rKey))
    }
    return c.pool.Do(radix.Pipeline(cmds...))
}

func (c *redis) InvalidateTags(tags ...string) error {
    tags = uniqueTags(tags)
    var rKeys, rTags []string
    for _, t := range tags {
        var keys []string
        t = c.t(t)
        err := c.pool.Do(radix.Cmd(&keys, "SMEMBERS", t))
        if err != nil {
            return err
        }
        rKeys = append(rKeys, keys...)
        rTags = append(rTags, t)
    }

    go func() {
        var cmds []radix.CmdAction
        for _, k := range rKeys {
            for _, t := range rTags {
                cmds = append(cmds, radix.FlatCmd(nil, "SREM", t, k))
            }
        }
        _ = c.pool.Do(radix.Pipeline(cmds...))
    }()

    if len(rKeys) == 0 {
        return nil
    }

    return c.pool.Do(radix.Cmd(nil, "DEL", rKeys...))
}

package cachita

import (
	"sync"
	"time"
)

var mCache Cache

type memory struct {
	recordsMu sync.RWMutex
	records   map[string]*record
	tagsMu    sync.Mutex
	tags      map[string][]string
	ttl       time.Duration
}

func Memory() Cache {
	if mCache == nil {
		mCache = NewMemoryCache(1*time.Minute, 1*time.Minute)
	}
	return mCache
}

func NewMemoryCache(ttl, tickerTtl time.Duration) Cache {
	c := &memory{
		records: make(map[string]*record),
		tags:    make(map[string][]string),
		ttl:     ttl,
	}

	runEvery(tickerTtl, func() {
		c.deleteExpired()
	})

	return c
}

func (c *memory) Get(key string, i interface{}) error {
	c.recordsMu.RLock()
	r, exists := c.records[key]
	c.recordsMu.RUnlock()
	if !exists {
		return ErrNotFound
	}
	if r.ExpiredAt.Before(time.Now()) {
		return ErrExpired
	}
	return TypeAssert(r.Data, i)
}

func (c *memory) Put(key string, i interface{}, ttl time.Duration) error {
	r := &record{Data: i, ExpiredAt: expiredAt(ttl, c.ttl)}
	c.recordsMu.Lock()
	defer c.recordsMu.Unlock()
	c.records[key] = r
	return nil
}

func (c *memory) Incr(key string, ttl time.Duration) (int64, error) {
	var n int64
	err := c.Get(key, &n)
	if err != nil && err != ErrNotFound && err != ErrExpired {
		return 0, err
	}
	n++
	err = c.Put(key, n, ttl)
	return n, err
}

func (c *memory) Invalidate(key string) error {
	c.recordsMu.Lock()
	defer c.recordsMu.Unlock()
	delete(c.records, key)
	return nil
}

func (c *memory) Exists(key string) bool {
	c.recordsMu.RLock()
	defer c.recordsMu.RUnlock()
	r, exists := c.records[key]
	return exists && r.ExpiredAt.After(time.Now())
}

func (c *memory) deleteExpired() {
	records := make(map[string]*record)
	c.recordsMu.Lock()
	defer c.recordsMu.Unlock()
	for k, r := range c.records {
		if r.ExpiredAt.After(time.Now()) {
			records[k] = r
		}
	}
	c.records = records
}

func (c *memory) InvalidateMulti(keys ...string) error {
	c.recordsMu.Lock()
	defer c.recordsMu.Unlock()
	for _, key := range keys {
		delete(c.records, key)
	}
	return nil
}

func (c *memory) Tag(key string, tags ...string) error {
	c.tagsMu.Lock()
	defer c.tagsMu.Unlock()
	for _, t := range tags {
		if inArr(c.tags[t], key) {
			continue
		}
		c.tags[t] = append(c.tags[t], key)
	}

	return nil
}

func (c *memory) InvalidateTags(tags ...string) error {
	c.tagsMu.Lock()
	var keys []string
	for _, t := range tags {
		keys = append(keys, c.tags[t]...)
		delete(c.tags, t)
	}
	c.tagsMu.Unlock()

	return c.InvalidateMulti(keys...)
}

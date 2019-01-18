package cachita

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func newCache(c Cache, t *testing.T) {
	var d string
	s := "╭∩╮(Ο_Ο)╭∩╮"
	test(c, "٩(̾●̮̮̃̾•̃̾)۶", &s, &d, t)
}

func cacheExpires(c Cache, t *testing.T, ttl, tts time.Duration) {
	var d string
	s := "><(((('>"
	k := "><>"
	err := c.Put(k, s, ttl)
	isError(err, t)
	assert.True(t, c.Exists(k))
	time.Sleep(tts)
	assert.False(t, c.Exists(k))

	err = c.Get(k, &d)
	assert.Equal(t, ErrNotFound, err)
	if err != ErrNotFound {
		isError(err, t)
	}
}

func test(c Cache, k string, s, d interface{}, t assert.TestingT, f ... func(t assert.TestingT, s, d interface{})) {
	k = fmt.Sprintf("%s%d", k, rand.Int())
	disableAssert := isBenchmark(t)

	err := c.Put(k, s, 0)
	isError(err, t)
	if !disableAssert {
		assert.True(t, c.Exists(k))
	}

	err = c.Get(k, d)
	isError(err, t)

	if !disableAssert {
		if len(f) > 0 {
			f[0](t, s, d)
		} else {
			assert.Equal(t, s, d)
		}
	}

	err = c.Invalidate(k)
	isError(err, t)
	if !disableAssert {
		assert.False(t, c.Exists(k))
	}

}
func isBenchmark(t assert.TestingT) bool {
	_, ok := t.(*testing.B)
	return ok
}

func testIncr(c Cache, k string, t assert.TestingT) {
	k = fmt.Sprintf("%s%d", k, rand.Int())
	disableAssert := isBenchmark(t)

	n, err := c.Incr(k, 0)
	isError(err, t)
	if !disableAssert {
		assert.True(t, c.Exists(k))
		assert.Equal(t, int64(1), n)
	}

	n, err = c.Incr(k, 0)
	isError(err, t)

	if !disableAssert {
		assert.Equal(t, int64(2), n)
	}

	err = c.Invalidate(k)
	isError(err, t)
	if !disableAssert {
		assert.False(t, c.Exists(k))
	}
}

func benchmarkCacheWithInt(c Cache, b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cacheWithInt(c, b)
		}
	})
}

func cacheWithInt(c Cache, t assert.TestingT) {
	s := 10000
	var d int
	test(c, "k", &s, &d, t)
}

func benchmarkCacheWithString(c Cache, b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cacheWithString(c, b)
		}
	})
}

func cacheWithString(c Cache, t assert.TestingT) {
	s := "test"
	var d string
	test(c, "k", &s, &d, t)
}

func benchmarkCacheWithMapInterface(c Cache, b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cacheWithMapInterface(c, b)
		}
	})
}

func cacheWithMapInterface(c Cache, t assert.TestingT) {
	s := map[string]interface{}{
		"Ƹ̵̡Ӝ̵̨̄Ʒ": 1,
		"ô¿ô":      "┌∩┐(◣_◢)┌∩┐",
		"●̮̮̃̾•": []interface{}{
			"ooo", "°º¤ø,¸¸,ø¤º°`°º¤ø,¸,ø¤°º¤ø,¸¸,ø¤º°`°º¤ø,¸",
			"♫♪.ılılıll|̲̅̅●̲̅̅|̲̅̅=̲̅̅|̲̅̅●̲̅̅|llılılı.♫♪",
		},
	}
	var d map[string]interface{}
	test(c, "k", &s, &d, t, compareMap)
}

func benchmarkCacheWithStruct(c Cache, b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cacheWithStruct(c, b)
		}
	})
}

func cacheWithStruct(c Cache, t assert.TestingT) {
	type str struct {
		A int
		B string
		C []string
		D map[string]interface{}
	}
	s := str{
		A: 1,
		B: "┌∩┐(◣_◢)┌∩┐",
		C: []string{
			"(♥_♥)", "ε(´סּ︵סּ`)з",
		},
		D: map[string]interface{}{
			"Ƹ̵̡Ӝ̵̨̄Ʒ": 1,
			"ô¿ô":      "┌∩┐(◣_◢)┌∩┐",
			"(̾●̮̮̃̾•̃̾)۶": []interface{}{
				"ooo", "°º¤ø,¸¸,ø¤º°`°º¤ø,¸,ø¤°º¤ø,¸¸,ø¤º°`°º¤ø,¸",
				"♫♪.ılılıll|̲̅̅●̲̅̅|̲̅̅=̲̅̅|̲̅̅●̲̅̅|llılılı.♫♪",
			},
		},
	}
	var d str
	test(c, "k", &s, &d, t, func(t assert.TestingT, s1, d1 interface{}) {
		s := s1.(*str)
		d := d1.(*str)
		assert.Equal(t, s.A, d.A)
		assert.Equal(t, s.B, d.B)
		assert.Equal(t, s.C, d.C)
		compareMap(t, &s.D, &d.D)
	})
}

func cacheIncr(c Cache, t assert.TestingT) {
	testIncr(c, "x", t)
}

func benchmarkCacheIncr(c Cache, b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cacheIncr(c, b)
		}
	})
}

func compareMap(t assert.TestingT, s1, d1 interface{}) {
	s := *s1.(*map[string]interface{})
	d := *d1.(*map[string]interface{})
	for k := range s {
		assert.EqualValues(t, s[k], d[k])
	}
}

func isError(err error, t assert.TestingT) {
	if err != nil {
		t.Errorf("%s", err)
	}
}

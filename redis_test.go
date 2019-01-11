package cachita

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRedisCache(t *testing.T) {
	t.Parallel()
	newCache(rc(t), t)
}

func TestRedisCacheExpires(t *testing.T) {
	t.Parallel()
	cacheExpires(rc(t), t, time.Second, 1200*time.Millisecond)
}

func TestRedisCacheWithInt(t *testing.T) {
	t.Parallel()
	cacheWithInt(rc(t), "x", t)
}
func BenchmarkRedisCacheWithInt(b *testing.B) {
	benchmarkParallel(rc(b), b, cacheWithInt)
}

func TestRedisCacheWithString(t *testing.T) {
	t.Parallel()
	cacheWithString(rc(t), "x1", t)
}

func BenchmarkRedisCacheWithString(b *testing.B) {
	benchmarkParallel(rc(b), b, cacheWithString)
}

func TestRedisCacheWithMapInterface(t *testing.T) {
	t.Parallel()
	cacheWithMapInterface(rc(t), "x2", t)
}

func BenchmarkRedisCacheWithMapInterface(b *testing.B) {
	benchmarkParallel(rc(b), b, cacheWithMapInterface)
}

func TestRedisCacheWithStruct(t *testing.T) {
	t.Parallel()
	cacheWithStruct(rc(t), "x3", t)
}

func BenchmarkRedisCacheWithStruct(b *testing.B) {
	benchmarkParallel(rc(b), b, cacheWithStruct)
}

func rc(t assert.TestingT) (c Cache) {
	c, err := Redis("127.0.0.1:6379")
	isError(err, t)
	return
}

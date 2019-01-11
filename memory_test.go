package cachita

import (
	"testing"
	"time"
)

func TestNewMemoryCache(t *testing.T) {
	t.Parallel()
	newCache(Memory(), t)
}

func TestMemoryCacheExpires(t *testing.T) {
	t.Parallel()
	cacheExpires(NewMemoryCache(2*time.Minute, 5*time.Millisecond), t, 50*time.Millisecond, 150*time.Millisecond)
}

func TestMemoryCacheWithInt(t *testing.T) {
	t.Parallel()
	cacheWithInt(Memory(), "x", t)
}

func BenchmarkMemoryCacheWithInt(b *testing.B) {
	benchmarkCacheWithInt(Memory(), b)
}

func TestMemoryCacheWithString(t *testing.T) {
	t.Parallel()
	cacheWithString(Memory(), "x1", t)
}

func BenchmarkMemoryCacheWithString(b *testing.B) {
	benchmarkCacheWithString(Memory(), b)
}

func TestMemoryCacheWithMapInterface(t *testing.T) {
	t.Parallel()
	cacheWithMapInterface(Memory(), "x2", t)
}

func BenchmarkMemoryCacheWithMapInterface(b *testing.B) {
	benchmarkCacheWithMapInterface(Memory(), b)
}

func TestMemoryCacheWithStruct(t *testing.T) {
	t.Parallel()
	cacheWithStruct(Memory(), "x3", t)
}

func BenchmarkMemoryCacheWithStruct(b *testing.B) {
	benchmarkCacheWithStruct(Memory(), b)
}

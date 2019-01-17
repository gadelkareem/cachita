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
	cacheWithInt(Memory(), t)
}

func BenchmarkMemoryCacheWithInt(b *testing.B) {
	benchmarkCacheWithInt(Memory(), b)
}

func TestMemoryCacheWithString(t *testing.T) {
	t.Parallel()
	cacheWithString(Memory(), t)
}

func BenchmarkMemoryCacheWithString(b *testing.B) {
	benchmarkCacheWithString(Memory(), b)
}

func TestMemoryCacheWithMapInterface(t *testing.T) {
	t.Parallel()
	cacheWithMapInterface(Memory(), t)
}

func BenchmarkMemoryCacheWithMapInterface(b *testing.B) {
	benchmarkCacheWithMapInterface(Memory(), b)
}

func TestMemoryCacheWithStruct(t *testing.T) {
	t.Parallel()
	cacheWithStruct(Memory(), t)
}

func BenchmarkMemoryCacheWithStruct(b *testing.B) {
	benchmarkCacheWithStruct(Memory(), b)
}

func TestMemory_Incr(t *testing.T) {
	t.Parallel()
	cacheIncr(Memory(), t)
}

func BenchmarkMemory_Incr(b *testing.B) {
	benchmarkCacheIncr(Memory(), b)
}

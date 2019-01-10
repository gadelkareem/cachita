package cachita

import (
	"testing"
	"time"
)

func TestNewMemoryCache(t *testing.T) {
	newCache(Memory(), t)
}

func TestMemoryCacheExpires(t *testing.T) {
	cacheExpires(NewMemoryCache(2*time.Minute, 5*time.Millisecond), t)
}

func TestMemoryCacheWithInt(t *testing.T) {
	cacheWithInt(Memory(), "x", t)
}

func BenchmarkMemoryCacheWithInt(b *testing.B) {
	benchmarkCacheWithInt(Memory(), b)
}

func TestMemoryCacheWithString(t *testing.T) {
	cacheWithString(Memory(), "x1", t)
}

func BenchmarkMemoryCacheWithString(b *testing.B) {
	benchmarkCacheWithString(Memory(), b)
}

func TestMemoryCacheWithMapInterface(t *testing.T) {
	cacheWithMapInterface(Memory(), "x2", t)
}

func BenchmarkMemoryCacheWithMapInterface(b *testing.B) {
	benchmarkCacheWithMapInterface(Memory(), b)
}

func TestMemoryCacheWithStruct(t *testing.T) {
	cacheWithStruct(Memory(), "x3", t)
}

func BenchmarkMemoryCacheWithStruct(b *testing.B) {
	benchmarkCacheWithStruct(Memory(), b)
}

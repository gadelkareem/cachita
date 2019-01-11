package cachita

import (
	"testing"
	"time"
)

func TestNewMemoryCache(t *testing.T) {
	newCache(Memory(), t)
}

func TestMemoryCacheExpires(t *testing.T) {
	cacheExpires(NewMemoryCache(2*time.Minute, 5*time.Millisecond), t, 50*time.Millisecond, 150*time.Millisecond)
}

func TestMemoryCacheWithInt(t *testing.T) {
	cacheWithInt(Memory(), "x", t)
}

func BenchmarkMemoryCacheWithInt(b *testing.B) {
	benchmark(rc(b), b, cacheWithInt)
}

func TestMemoryCacheWithString(t *testing.T) {
	cacheWithString(Memory(), "x1", t)
}

func BenchmarkMemoryCacheWithString(b *testing.B) {
	benchmark(rc(b), b, cacheWithString)
}

func TestMemoryCacheWithMapInterface(t *testing.T) {
	cacheWithMapInterface(Memory(), "x2", t)
}

func BenchmarkMemoryCacheWithMapInterface(b *testing.B) {
	benchmark(rc(b), b, cacheWithMapInterface)
}

func TestMemoryCacheWithStruct(t *testing.T) {
	cacheWithStruct(Memory(), "x3", t)
}

func BenchmarkMemoryCacheWithStruct(b *testing.B) {
	benchmark(rc(b), b, cacheWithStruct)
}

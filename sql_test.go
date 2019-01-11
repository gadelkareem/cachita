package cachita

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSqlCache(t *testing.T) {
	t.Parallel()
	newCache(sc(t), t)
}

func TestSqlCacheExpires(t *testing.T) {
	t.Parallel()
	sqlDriver, err := sql.Open("postgres", "postgres://postgres@localhost/test?sslmode=disable")
	isError(err, t)
	c, err := NewSqlCache(2*time.Minute, time.Second, sqlDriver, "cachita_cache", true)
	isError(err, t)
	cacheExpires(c, t, time.Second, 1200*time.Millisecond)
}

func TestSqlCacheWithInt(t *testing.T) {
	t.Parallel()
	cacheWithInt(sc(t), "x", t)
}
func BenchmarkSqlCacheWithInt(b *testing.B) {
	benchmarkCacheWithInt(sc(b), b)
}

func TestSqlCacheWithString(t *testing.T) {
	t.Parallel()
	cacheWithString(sc(t), "x1", t)
}

func BenchmarkSqlCacheWithString(b *testing.B) {
	benchmarkCacheWithString(sc(b), b)
}

func TestSqlCacheWithMapInterface(t *testing.T) {
	t.Parallel()
	cacheWithMapInterface(sc(t), "x2", t)
}

func BenchmarkSqlCacheWithMapInterface(b *testing.B) {
	benchmarkCacheWithMapInterface(sc(b), b)
}

func TestSqlCacheWithStruct(t *testing.T) {
	t.Parallel()
	cacheWithStruct(sc(t), "x3", t)
}

func BenchmarkSqlCacheWithStruct(b *testing.B) {
	benchmarkCacheWithStruct(sc(b), b)
}

func sc(t assert.TestingT) (c Cache) {
	c, err := Sql("postgres", "postgres://postgres@localhost/test?sslmode=disable")
	isError(err, t)
	return
}
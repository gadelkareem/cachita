package cachita

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFileCache(t *testing.T) {
	t.Parallel()
	c, err := File()
	isError(err, t)
	newCache(c, t)
}

func TestFileCacheExpires(t *testing.T) {
	t.Parallel()
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	isError(err, t)
	path = filepath.Join(path, "tmp2/file-cache")
	c, err := NewFileCache(path, 2*time.Minute, 100*time.Millisecond)
	isError(err, t)
	cacheExpires(c, t)
}

func TestFileCacheWithInt(t *testing.T) {
	t.Parallel()
	c, err := File()
	isError(err, t)
	cacheWithInt(c, "x", t)
}
func BenchmarkFileCacheWithInt(b *testing.B) {
	c, err := File()
	isError(err, b)
	benchmarkCacheWithInt(c, b)
}

func TestFileCacheWithString(t *testing.T) {
	t.Parallel()
	c, err := File()
	isError(err, t)
	cacheWithString(c, "x1", t)
}

func BenchmarkFileCacheWithString(b *testing.B) {
	c, err := File()
	isError(err, b)
	benchmarkCacheWithString(c, b)
}

func TestFileCacheWithMapInterface(t *testing.T) {
	t.Parallel()
	c, err := File()
	isError(err, t)
	cacheWithMapInterface(c, "x2", t)
}

func BenchmarkFileCacheWithMapInterface(b *testing.B) {
	c, err := File()
	isError(err, b)
	benchmarkCacheWithMapInterface(c, b)
}

func TestFileCacheWithStruct(t *testing.T) {
	t.Parallel()
	c, err := File()
	isError(err, t)
	cacheWithStruct(c, "x3", t)
}

func BenchmarkFileCacheWithStruct(b *testing.B) {
	c, err := File()
	isError(err, b)
	benchmarkCacheWithStruct(c, b)
}

func TestIndexFileCreated(t *testing.T) {
	t.Parallel()
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	isError(err, t)
	path = filepath.Join(path, "tmp5/file-cache")
	c, err := NewFileCache(path, 1*time.Millisecond, 100*time.Millisecond)
	isError(err, t)
	assert.NotNil(t, c)
	k := "x4"
	s := "⺌∅‿∅⺌"
	ttl := 1 * time.Hour
	isError(c.Put(k, s, ttl), t)
	time.Sleep(150 * time.Millisecond)
	indexPath := filepath.Join(path, Id(FileIndex))
	assert.FileExists(t, indexPath)
	var i fileIndex
	isError(readData(indexPath, &i.records), t)
	e, exists := i.records[Id(k)]
	assert.True(t, exists, "Index file should have the record id")
	assert.True(t, e.After(time.Now().Add(ttl-2*time.Second)), "Index expiry should equal first set expiry")
	c, err = NewFileCache(path, 1*time.Millisecond, 1*time.Hour)
	isError(err, t)
	var d string
	isError(c.Get(k, &d), t)
	assert.Equal(t, &s, &d)
}

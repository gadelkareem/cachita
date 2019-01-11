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
	newCache(fc(t), t)
}

func TestFileCacheExpires(t *testing.T) {
	t.Parallel()
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	isError(err, t)
	path = filepath.Join(path, "tmp2/file-cache")
	c, err := NewFileCache(path, 2*time.Minute, 100*time.Millisecond)
	isError(err, t)
	cacheExpires(c, t, 50*time.Millisecond, 150*time.Millisecond)
}

func TestFileCacheWithInt(t *testing.T) {
	t.Parallel()
	cacheWithInt(fc(t), "x", t)
}
func BenchmarkFileCacheWithInt(b *testing.B) {
	benchmark(rc(b), b, cacheWithInt)
}

func TestFileCacheWithString(t *testing.T) {
	t.Parallel()
	cacheWithString(fc(t), "x1", t)
}

func BenchmarkFileCacheWithString(b *testing.B) {
	benchmark(rc(b), b, cacheWithString)
}

func TestFileCacheWithMapInterface(t *testing.T) {
	t.Parallel()
	cacheWithMapInterface(fc(t), "x2", t)
}

func BenchmarkFileCacheWithMapInterface(b *testing.B) {
	benchmark(rc(b), b, cacheWithMapInterface)
}

func TestFileCacheWithStruct(t *testing.T) {
	t.Parallel()
	cacheWithStruct(fc(t), "x3", t)
}

func BenchmarkFileCacheWithStruct(b *testing.B) {
	benchmark(rc(b), b, cacheWithStruct)
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

func fc(t assert.TestingT) (c Cache) {
	c, err := File()
	isError(err, t)
	return
}

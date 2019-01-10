package cachita

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

//----- Memory

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

//---- File

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
	cacheWithInt(Memory(), "x", t)
}
func BenchmarkFileCacheWithInt(b *testing.B) {
	benchmarkCacheWithInt(Memory(), b)
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

//-----------------------------------------------------

func newCache(c Cache, t *testing.T) {
	var d string
	s := "╭∩╮(Ο_Ο)╭∩╮"
	k := "٩(̾●̮̮̃̾•̃̾)۶"
	test(c, k, &s, &d, t)
}

func cacheExpires(c Cache, t *testing.T) {
	var d string
	s := "><(((('>"
	k := "><>"
	err := c.Put(k, s, 50*time.Millisecond)
	isError(err, t)
	assert.True(t, c.Exists(k))
	time.Sleep(150 * time.Millisecond)
	assert.False(t, c.Exists(k))

	err = c.Get(k, &d)
	assert.Equal(t, err, ErrNotFound)
}

func test(c Cache, k string, s, d interface{}, t assert.TestingT, f ... func(t assert.TestingT, s, d interface{})) {
	err := c.Put(k, s, 0)
	isError(err, t)

	assert.True(t, c.Exists(k))

	err = c.Get(k, d)
	isError(err, t)

	if len(f) > 0 {
		f[0](t, s, d)
	} else {
		assert.Equal(t, s, d)
	}

	err = c.Invalidate(k)
	isError(err, t)
	assert.False(t, c.Exists(k))
}

func benchmarkCacheWithInt(c Cache, b *testing.B) {
	for n := 0; n < b.N; n++ {
		cacheWithInt(c, string(b.N), b)
	}
}

func cacheWithInt(c Cache, k string, t assert.TestingT) {
	s := 10000
	var d int
	test(c, k, &s, &d, t)
}

func benchmarkCacheWithString(c Cache, b *testing.B) {
	for n := 0; n < b.N; n++ {
		cacheWithString(c, string(b.N), b)
	}
}

func cacheWithString(c Cache, k string, t assert.TestingT) {
	s := "test"
	var d string
	test(c, k, &s, &d, t)
}

func benchmarkCacheWithMapInterface(c Cache, b *testing.B) {
	for n := 0; n < b.N; n++ {
		cacheWithMapInterface(c, string(b.N), b)
	}
}

func cacheWithMapInterface(c Cache, k string, t assert.TestingT) {
	s := map[string]interface{}{
		"Ƹ̵̡Ӝ̵̨̄Ʒ": 1,
		"ô¿ô":      "┌∩┐(◣_◢)┌∩┐",
		"●̮̮̃̾•": []interface{}{
			"ooo", "°º¤ø,¸¸,ø¤º°`°º¤ø,¸,ø¤°º¤ø,¸¸,ø¤º°`°º¤ø,¸",
			"♫♪.ılılıll|̲̅̅●̲̅̅|̲̅̅=̲̅̅|̲̅̅●̲̅̅|llılılı.♫♪",
		},
	}
	var d map[string]interface{}
	test(c, k, &s, &d, t, compareMap)
}

func benchmarkCacheWithStruct(c Cache, b *testing.B) {
	for n := 0; n < b.N; n++ {
		cacheWithStruct(c, string(b.N), b)
	}
}

func cacheWithStruct(c Cache, k string, t assert.TestingT) {
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
	test(c, k, &s, &d, t, func(t assert.TestingT, s1, d1 interface{}) {
		s := s1.(*str)
		d := d1.(*str)
		assert.Equal(t, s.A, d.A)
		assert.Equal(t, s.B, d.B)
		assert.Equal(t, s.C, d.C)
		compareMap(t, &s.D, &d.D)
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

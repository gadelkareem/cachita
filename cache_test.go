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
	c, err := File()
	isError(err, t)
	newCache(c, t)
}

func TestFileCacheExpires(t *testing.T) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	isError(err, t)
	path = filepath.Join(path, "tmp/file-cache")
	c, err := NewFileCache(path, 2*time.Minute, 5*time.Millisecond)
	isError(err, t)
	cacheExpires(c, t)
}

func TestFileCacheWithInt(t *testing.T) {
	cacheWithInt(Memory(), "x", t)
}
func BenchmarkFileCacheWithInt(b *testing.B) {
	benchmarkCacheWithInt(Memory(), b)
}

func TestFileCacheWithString(t *testing.T) {
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
	t.Skip("msgpack problem")
	c, err := File()
	isError(err, t)
	cacheWithMapInterface(c, "x2", t)
}

func BenchmarkFileCacheWithMapInterface(b *testing.B) {
	b.Skip("msgpack problem")
	c, err := File()
	isError(err, b)
	benchmarkCacheWithMapInterface(c, b)
}

func TestFileCacheWithStruct(t *testing.T) {
	t.Skip("msgpack problem")
	c, err := File()
	isError(err, t)
	cacheWithStruct(c, "x3", t)
}

func BenchmarkFileCacheWithStruct(b *testing.B) {
	b.Skip("msgpack problem")
	c, err := File()
	isError(err, b)
	benchmarkCacheWithStruct(c, b)
}

//-----------------------------------------------------

func newCache(c Cache, t *testing.T) {
	t.Parallel()
	var d string
	s := "╭∩╮(Ο_Ο)╭∩╮"
	k := "٩(̾●̮̮̃̾•̃̾)۶"
	test(c, k, &s, &d, t)
}

func cacheExpires(c Cache, t *testing.T) {
	t.Parallel()
	var d string
	s := "><(((('>"
	k := "><>"
	err := c.Put(k, s, 1*time.Millisecond)
	isError(err, t)
	time.Sleep(10 * time.Millisecond)
	assert.False(t, c.Exists(k))

	err = c.Get(k, &d)
	assert.Equal(t, err, ErrNotFound)
}

func test(c Cache, k string, s, d interface{}, t assert.TestingT) {
	err := c.Put(k, s, 0)
	isError(err, t)

	assert.True(t, c.Exists(k))

	err = c.Get(k, d)
	isError(err, t)

	assert.Equal(t, s, d)

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
		"●̮̮̃̾•": []string{
			"ooo", "°º¤ø,¸¸,ø¤º°`°º¤ø,¸,ø¤°º¤ø,¸¸,ø¤º°`°º¤ø,¸",
			"♫♪.ılılıll|̲̅̅●̲̅̅|̲̅̅=̲̅̅|̲̅̅●̲̅̅|llılılı.♫♪",
		},
	}
	var d map[string]interface{}
	test(c, k, &s, &d, t)
}

func benchmarkCacheWithStruct(c Cache, b *testing.B) {
	for n := 0; n < b.N; n++ {
		cacheWithStruct(c, string(b.N), b)
	}
}

func cacheWithStruct(c Cache, k string, t assert.TestingT) {
	type str struct {
		a int
		b string
		c []string
		d map[string]interface{}
	}
	s := str{
		a: 1,
		b: "┌∩┐(◣_◢)┌∩┐",
		c: []string{
			"(♥_♥)", "ε(´סּ︵סּ`)з",
		},
		d: map[string]interface{}{
			"Ƹ̵̡Ӝ̵̨̄Ʒ": 1,
			"ô¿ô":      "┌∩┐(◣_◢)┌∩┐",
			"(̾●̮̮̃̾•̃̾)۶": []string{
				"ooo", "°º¤ø,¸¸,ø¤º°`°º¤ø,¸,ø¤°º¤ø,¸¸,ø¤º°`°º¤ø,¸",
				"♫♪.ılılıll|̲̅̅●̲̅̅|̲̅̅=̲̅̅|̲̅̅●̲̅̅|llılılı.♫♪",
			},
		},
	}
	var d str
	test(c, k, &s, &d, t)
}

func isError(err error, t assert.TestingT) {
	if err != nil {
		t.Errorf("%s", err)
	}
}

package cachita

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack"
)

const FileIndex = "github.com/gadelkareem/cachita/file-index"

var fCache Cache

type file struct {
	dir string
	ttl time.Duration
	i   *fileIndex
}

type fileIndex struct {
	sync.RWMutex
	records map[string]time.Time
	path    string
}

func File() (Cache, error) {
	if fCache == nil {
		path, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return nil, err
		}
		path = filepath.Join(path, "tmp/file-cache")
		fCache, err = NewFileCache(path, 24*time.Hour, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return fCache, nil
}

func NewFileCache(dir string, ttl, tickerTtl time.Duration) (Cache, error) {
	var (
		err error
	)
	i, err := newIndex(dir, ttl)
	if err != nil {
		return nil, err
	}

	c := &file{
		dir: dir,
		ttl: ttl,
		i:   i,
	}
	if tickerTtl != 0 {
		runEvery(tickerTtl, func() {
			c.deleteExpired()
		})
	}
	return c, nil
}

func (c *file) Exists(key string) bool {
	err := c.i.check(Id(key))
	return err == nil
}

func (c *file) Get(key string, i interface{}) error {
	id := Id(key)
	if err := c.i.check(id); err != nil {
		return err
	}
	return readData(c.path(id), i)
}

func (c *file) Put(key string, i interface{}, ttl time.Duration) error {
	id := Id(key)
	c.i.add(id, expiredAt(ttl, c.ttl))
	return writeData(c.path(id), i)
}

func (c *file) Incr(key string, ttl time.Duration) (int64, error) {
	var n int64
	err := c.Get(key, &n)
	if err != nil && err != ErrNotFound && err != ErrExpired {
		return 0, err
	}
	n++
	err = c.Put(key, n, ttl)
	return n, err
}

func (c *file) Invalidate(key string) error {
	id := Id(key)
	c.i.remove(id)
	err := os.Remove(c.path(id))
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}

func (c *file) path(id string) string {
	return filepath.Join(c.dir, string(id[0]), string(id[1]), id)
}

func (c *file) deleteExpired() {
	expired := c.i.expiredRecords()
	for _, id := range expired {
		os.Remove(c.path(id))
	}
}

//----------------------- fileIndex

func newIndex(dir string, ttl time.Duration) (i *fileIndex, err error) {
	i = &fileIndex{path: filepath.Join(dir, Id(FileIndex))}
	i.records = make(map[string]time.Time)

	err = readData(i.path, &i.records)
	if err != nil && err != ErrNotFound {
		return
	}

	var (
		currentDir string
		files      []os.FileInfo
	)
	i.Lock()
	defer i.Unlock()
	characters := "0123456789abcdef"
	for _, char1 := range characters {
		for _, char2 := range characters {
			currentDir = filepath.Join(dir, string(char1), string(char2))
			if ok, _ := exists(currentDir); !ok {
				err = os.MkdirAll(currentDir, os.ModePerm)
				if err != nil {
					return
				}
			}
			files, err = ioutil.ReadDir(currentDir)
			if err != nil {
				return
			}

			for _, f := range files {
				if f.IsDir() {
					continue
				}
				if _, exists := i.records[f.Name()]; exists {
					continue
				}
				expiredAt := f.ModTime().Add(ttl)
				if expiredAt.After(time.Now()) {
					i.records[f.Name()] = expiredAt
				}
			}
		}
	}
	err = writeData(i.path, &i.records)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *fileIndex) check(id string) error {
	i.RLock()
	defer i.RUnlock()
	expiredAt, exists := i.records[id]
	if !exists {
		return ErrNotFound
	}
	if expiredAt.Before(time.Now()) {
		return ErrExpired
	}
	return nil
}

func (i *fileIndex) expiredRecords() []string {
	i.Lock()
	defer i.Unlock()
	var (
		expired []string
		records = make(map[string]time.Time)
	)
	for id, expiredAt := range i.records {
		if expiredAt.Before(time.Now()) {
			expired = append(expired, id)
			continue
		}
		records[id] = expiredAt
	}
	i.records = records
	err := writeData(i.path, &i.records)
	if err != nil {
		fmt.Printf("cachita: error writing index file: %v", err)
	}
	return expired
}

func (i *fileIndex) add(id string, expiredAt time.Time) {
	i.Lock()
	defer i.Unlock()
	i.records[id] = expiredAt
	return
}

func (i *fileIndex) remove(id string) {
	i.Lock()
	defer i.Unlock()
	delete(i.records, id)
}

//--------------------

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readData(path string, i interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) || err == io.EOF {
			return ErrNotFound
		}
		return err
	}
	err = msgpack.Unmarshal(data, i)
	if err != nil {
		if err == io.EOF {
			return ErrNotFound
		}
		return err
	}
	return nil
}
func writeData(path string, i interface{}) error {
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}

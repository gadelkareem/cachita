package cache

import (
	"fmt"
	"github.com/gadelkareem/go-helpers"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
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
		fCache, err = NewFileCache(path, 24*time.Hour, 1*time.Hour)
		if err != nil {
			return nil, err
		}
	}
	return fCache, nil
}

func NewFileCache(dir string, ttl, tickerTtl time.Duration) (Cache, error) {
	var (
		currentDir string
		err        error
	)
	if ok, _ := exists(dir); !ok {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	characters := "0123456789abcdefghijklmnopqrstuvwxyz"
	for _, char1 := range characters {
		for _, char2 := range characters {
			currentDir = filepath.Join(dir, string(char1), string(char2))
			if ok, _ := exists(currentDir); !ok {
				err = os.MkdirAll(currentDir, os.ModePerm)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	i, err := newIndex(dir, ttl)
	if err != nil {
		return nil, err
	}

	c := &file{
		dir: dir,
		ttl: ttl,
		i:   i,
	}

	helpers.RunEvery(tickerTtl, func() {
		err := c.Put(FileIndex, &c.i, -1)
		if err != nil {
			fmt.Printf("cachita: error writing index file: %v", err)
		}
		c.deleteExpired()
	})

	return c, nil
}

func (c *file) Exists(key string) bool {
	err := c.i.check(id(key))
	return err == nil
}

func (c *file) Get(key string, i interface{}) error {
	id := id(key)
	if err := c.i.check(id); err != nil {
		return err
	}
	return readData(c.path(id), i)
}
func (c *file) Put(key string, i interface{}, ttl time.Duration) error {
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}
	id := id(key)
	go c.i.add(id, ExpiredAt(ttl, c.ttl))
	return ioutil.WriteFile(c.path(id), data, 0666)
}

func (c *file) Invalidate(key string) error {
	id := id(key)
	c.i.remove(id)
	return os.Remove(c.path(id))
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
	i = &fileIndex{path: filepath.Join(dir, id(FileIndex))}
	i.records = make(map[string]time.Time)
	err = readData(i.path, &i)
	if err != nil && err != ErrNotFound {
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	i.Lock()
	defer i.Unlock()
	for _, f := range files {
		expiredAt := f.ModTime().Add(ttl)
		if expiredAt.Before(time.Now()) {
			i.records[f.Name()] = expiredAt
		}
	}

	return
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

func id(key string) string {
	return helpers.Md5(key)
}

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
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	err = msgpack.Unmarshal(data, i)
	if err != nil {
		return err
	}
	return nil
}

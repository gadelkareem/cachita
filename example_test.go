package cachita_test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gadelkareem/cachita"
)

func ExampleCache() {
	cache := cachita.Memory()
	err := cache.Put("cache_key", "some data", 1*time.Minute)
	if err != nil {
		panic(err)
	}

	if cache.Exists("cache_key") {
		//do something
	}

	var holder string
	err = cache.Get("cache_key", &holder)
	if err != nil && err != cachita.ErrNotFound {
		panic(err)
	}

	fmt.Printf("%s", holder) //prints "some data"

	err = cache.Invalidate("cache_key")
	if err != nil {
		panic(err)
	}

	//Output: some data

}

func ExampleMemory() {
	var u url.URL
	cacheId := cachita.Id(u.Scheme, u.Host, u.RequestURI())
	obj := make(map[string]interface{})
	obj["test"] = "data"
	err := cachita.Memory().Put(cacheId, obj, 0)
	if err != nil {
		panic(err)
	}

	var cacheObj map[string]interface{}
	err = cachita.Memory().Get(cacheId, &cacheObj)
	if err != nil && err != cachita.ErrNotFound && err != cachita.ErrExpired {
		panic(err)
	}
	fmt.Printf("%+v", cacheObj)

	//Output: map[test:data]

}

func ExampleFile() {
	cache, err := cachita.File()
	if err != nil {
		panic(err)
	}

	err = cache.Put("cache_key", "some data", 1*time.Minute)
	if err != nil {
		panic(err)
	}

	var holder string
	err = cache.Get("cache_key", &holder)
	if err != nil && err != cachita.ErrNotFound {
		panic(err)
	}

	fmt.Printf("%s", holder) //prints "some data"

	//Output: some data

}

func ExampleNewMemoryCache() {
	cache := cachita.NewMemoryCache(1*time.Millisecond, 1*time.Minute) //default ttl 1 millisecond

	err := cache.Put("cache_key", "some data", 0) //ttl = 0 means use default
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Millisecond)
	fmt.Printf("%t", cache.Exists("cache_key"))

	//Output: false

}

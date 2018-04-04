package cachita_test

import (
	"fmt"
	"github.com/gadelkareem/cachita"
	"net/url"
	"time"
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

	err = cache.Invalidate("cache_key")
	if err != nil {
		panic(err)
	}

}


func ExampleCachingPageContents() {
	var u url.URL
	cacheId := cachita.Id(u.Scheme, u.Host, u.RequestURI())
	var cacheObj map[string]interface{}
	err := cachita.Memory().Get(cacheId, cacheObj)
	if err != nil && err != cachita.ErrNotFound && err != cachita.ErrExpired {
		panic(err)
	}

	fmt.Printf("%+v", cacheObj)

}

func ExampleFileCache() {
	cache, err := cachita.File()
	if err != nil {
		panic(err)
	}

	err = cache.Put("cache_key", "some data", 1*time.Minute)
	if err != nil {
		panic(err)
	}
	//...

}

func ExampleCustomMemoryCache() {
	cache := cachita.NewMemoryCache(1*time.Minute, 1*time.Minute)

	err := cache.Put("cache_key", "some data", 1*time.Minute)
	if err != nil {
		panic(err)
	}
	//...

}

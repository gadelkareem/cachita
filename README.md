# Cachita
Cachita is a golang file, memory, SQL, Redis cache library 

[![Build Status](https://travis-ci.org/gadelkareem/cachita.svg)](https://travis-ci.org/gadelkareem/cachita)
[![GoDoc](https://godoc.org/github.com/gadelkareem/cachita?status.svg)](https://godoc.org/github.com/gadelkareem/cachita)

- Simple caching with auto type assertion included.
- In memory file cache index to avoid unneeded I/O.
- [Msgpack](https://msgpack.org/index.html) based binary serialization using [msgpack](https://github.com/vmihailenco/msgpack) library for file caching.


API docs: https://godoc.org/github.com/gadelkareem/cachita.

Examples: https://godoc.org/github.com/gadelkareem/cachita#pkg-examples.

## Installation

Install:

```shell
go get -u github.com/gadelkareem/cachita
```

## Quickstart

```go

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

```

## Benchmark

```
> go test -v -bench=. -benchmem
BenchmarkFileCacheWithInt-8              	    5000	    294487 ns/op	    3681 B/op	      46 allocs/op
BenchmarkFileCacheWithString-8           	    5000	    289703 ns/op	    3710 B/op	      47 allocs/op
BenchmarkFileCacheWithMapInterface-8     	    5000	    313341 ns/op	    6038 B/op	      73 allocs/op
BenchmarkFileCacheWithStruct-8           	    5000	    327784 ns/op	    7597 B/op	      94 allocs/op
BenchmarkMemoryCacheWithInt-8            	  300000	      5145 ns/op	     728 B/op	      14 allocs/op
BenchmarkMemoryCacheWithString-8         	  300000	      5344 ns/op	     756 B/op	      14 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	  200000	      7226 ns/op	    1444 B/op	      20 allocs/op
BenchmarkMemoryCacheWithStruct-8         	  200000	     11230 ns/op	    2184 B/op	      33 allocs/op
BenchmarkSqlCacheWithInt-8               	     100	  10351146 ns/op	    8422 B/op	     220 allocs/op
BenchmarkSqlCacheWithString-8            	     200	   9000744 ns/op	    8265 B/op	     212 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	     200	  16634691 ns/op	   14931 B/op	     456 allocs/op
BenchmarkSqlCacheWithStruct-8            	     200	   8496496 ns/op	   18325 B/op	     546 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


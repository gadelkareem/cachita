# Cachita
Cachita is a golang file, memory, SQL, Redis cache library 

[![Build Status](https://travis-ci.org/gadelkareem/cachita.svg)](https://travis-ci.org/gadelkareem/cachita)
[![GoDoc](https://godoc.org/github.com/gadelkareem/cachita?status.svg)](https://godoc.org/github.com/gadelkareem/cachita)

- Simple caching with auto type assertion included.
- In memory file cache index to avoid unneeded I/O.
- [Msgpack](https://msgpack.org/index.html) based binary serialization using [msgpack](https://github.com/vmihailenco/msgpack) library for file caching.
- [go-redis](https://github.com/go-redis/redis) Redis client.


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
BenchmarkMemoryCacheWithInt-8            	  300000	      4789 ns/op	     728 B/op	      14 allocs/op
BenchmarkMemoryCacheWithString-8         	  300000	      4886 ns/op	     756 B/op	      14 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	  200000	      7276 ns/op	    1444 B/op	      20 allocs/op
BenchmarkMemoryCacheWithStruct-8         	  200000	     10955 ns/op	    2184 B/op	      33 allocs/op
BenchmarkFileCacheWithInt-8              	    5000	    284634 ns/op	    3681 B/op	      46 allocs/op
BenchmarkFileCacheWithString-8           	    5000	    286430 ns/op	    3710 B/op	      47 allocs/op
BenchmarkFileCacheWithMapInterface-8     	    5000	    306349 ns/op	    6039 B/op	      73 allocs/op
BenchmarkFileCacheWithStruct-8           	    5000	    320347 ns/op	    7600 B/op	      94 allocs/op
BenchmarkRedisCacheWithInt-8             	     500	   2764354 ns/op	    2525 B/op	      67 allocs/op
BenchmarkRedisCacheWithString-8          	     500	   2927454 ns/op	    2548 B/op	      68 allocs/op
BenchmarkRedisCacheWithMapInterface-8    	     500	   2954466 ns/op	    5167 B/op	      95 allocs/op
BenchmarkRedisCacheWithStruct-8          	     500	   2884623 ns/op	    6802 B/op	     117 allocs/op
BenchmarkSqlCacheWithInt-8               	     100	  10689678 ns/op	    8437 B/op	     221 allocs/op
BenchmarkSqlCacheWithString-8            	     200	   9662987 ns/op	    8282 B/op	     212 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	     200	   8731074 ns/op	   14735 B/op	     453 allocs/op
BenchmarkSqlCacheWithStruct-8            	     200	   8001558 ns/op	   18316 B/op	     546 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


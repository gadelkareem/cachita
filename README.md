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
BenchmarkMemoryCacheWithInt-8            	 1000000	      1335 ns/op	     128 B/op	       7 allocs/op
BenchmarkMemoryCacheWithString-8         	 1000000	      1386 ns/op	     144 B/op	       7 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	 1000000	      1516 ns/op	     544 B/op	      11 allocs/op
BenchmarkMemoryCacheWithStruct-8         	 1000000	      1733 ns/op	     688 B/op	      12 allocs/op
BenchmarkFileCacheWithInt-8              	   10000	    109377 ns/op	    2951 B/op	      35 allocs/op
BenchmarkFileCacheWithString-8           	   10000	    113310 ns/op	    2967 B/op	      36 allocs/op
BenchmarkFileCacheWithMapInterface-8     	   10000	    122770 ns/op	    4998 B/op	      59 allocs/op
BenchmarkFileCacheWithStruct-8           	   10000	    120276 ns/op	    5903 B/op	      64 allocs/op
BenchmarkRedisCacheWithInt-8             	    3000	    803333 ns/op	    1471 B/op	      43 allocs/op
BenchmarkRedisCacheWithString-8          	    3000	    776620 ns/op	    1479 B/op	      44 allocs/op
BenchmarkRedisCacheWithMapInterface-8    	    3000	    591221 ns/op	    3797 B/op	      68 allocs/op
BenchmarkRedisCacheWithStruct-8          	    3000	    513130 ns/op	    4782 B/op	      74 allocs/op
BenchmarkSqlCacheWithInt-8               	    1000	   1856955 ns/op	    5136 B/op	     143 allocs/op
BenchmarkSqlCacheWithString-8            	    1000	   1937829 ns/op	    5022 B/op	     135 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	    1000	   1899502 ns/op	   10981 B/op	     373 allocs/op
BenchmarkSqlCacheWithStruct-8            	    1000	   1965525 ns/op	   13874 B/op	     451 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


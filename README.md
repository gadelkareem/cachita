# Cachita
Cachita is a golang file, memory, SQL, Redis cache library 

[![Build Status](https://travis-ci.org/gadelkareem/cachita.svg)](https://travis-ci.org/gadelkareem/cachita)
[![GoDoc](https://godoc.org/github.com/gadelkareem/cachita?status.svg)](https://godoc.org/github.com/gadelkareem/cachita)

- Simple caching with auto type assertion included.
- In memory file cache index to avoid unneeded I/O.
- [Msgpack](https://msgpack.org/index.html) based binary serialization using [msgpack](https://github.com/vmihailenco/msgpack) library for file caching.
- [radix](https://github.com/mediocregopher/radix) Redis client.


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
BenchmarkMemoryCacheWithInt-8            	 1000000	      1792 ns/op	     120 B/op	       6 allocs/op
BenchmarkMemoryCacheWithString-8         	 1000000	      1525 ns/op	     136 B/op	       6 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	 1000000	      1711 ns/op	     536 B/op	      10 allocs/op
BenchmarkMemoryCacheWithStruct-8         	 1000000	      1832 ns/op	     680 B/op	      11 allocs/op
BenchmarkMemory_Incr-8                   	 1000000	      2768 ns/op	     192 B/op	      10 allocs/op
BenchmarkFileCacheWithInt-8              	   10000	    128861 ns/op	    2946 B/op	      34 allocs/op
BenchmarkFileCacheWithString-8           	   10000	    127513 ns/op	    2968 B/op	      35 allocs/op
BenchmarkFileCacheWithMapInterface-8     	   10000	    126324 ns/op	    4998 B/op	      58 allocs/op
BenchmarkFileCacheWithStruct-8           	   10000	    144674 ns/op	    5905 B/op	      63 allocs/op
BenchmarkFile_Incr-8                     	    5000	    223372 ns/op	    7098 B/op	      74 allocs/op
BenchmarkRedisCacheWithInt-8             	    5000	    357409 ns/op	     705 B/op	      25 allocs/op
BenchmarkRedisCacheWithString-8          	    5000	    327421 ns/op	    1204 B/op	      35 allocs/op
BenchmarkRedisCacheWithMapInterface-8    	    5000	    341155 ns/op	    3284 B/op	      59 allocs/op
BenchmarkRedisCacheWithStruct-8          	    5000	    369480 ns/op	    4177 B/op	      64 allocs/op
BenchmarkRedis_Incr-8                    	    2000	    798698 ns/op	    1323 B/op	      45 allocs/op
BenchmarkSqlCacheWithInt-8               	    1000	   1998693 ns/op	    5179 B/op	     143 allocs/op
BenchmarkSqlCacheWithString-8            	     500	   2163233 ns/op	    5129 B/op	     135 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	     500	   2012257 ns/op	   11285 B/op	     374 allocs/op
BenchmarkSqlCacheWithStruct-8            	    1000	   2055111 ns/op	   14014 B/op	     452 allocs/op
BenchmarkSql_Incr-8                      	     500	   3192285 ns/op	    9687 B/op	     268 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


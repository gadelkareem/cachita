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
BenchmarkMemoryCacheWithInt-8            	  300000	      4885 ns/op	     728 B/op	      14 allocs/op
BenchmarkMemoryCacheWithString-8         	  300000	      4924 ns/op	     756 B/op	      14 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	  200000	      7326 ns/op	    1444 B/op	      20 allocs/op
BenchmarkMemoryCacheWithStruct-8         	  200000	     11038 ns/op	    2184 B/op	      33 allocs/op
BenchmarkFileCacheWithInt-8              	    5000	    288879 ns/op	    3685 B/op	      46 allocs/op
BenchmarkFileCacheWithString-8           	    5000	    290468 ns/op	    3713 B/op	      47 allocs/op
BenchmarkFileCacheWithMapInterface-8     	    5000	    311854 ns/op	    6042 B/op	      73 allocs/op
BenchmarkFileCacheWithStruct-8           	    5000	    330004 ns/op	    7606 B/op	      94 allocs/op
BenchmarkRedisCacheWithInt-8             	     500	   4158970 ns/op	    2282 B/op	      62 allocs/op
BenchmarkRedisCacheWithString-8          	     300	   4136612 ns/op	    2316 B/op	      63 allocs/op
BenchmarkRedisCacheWithMapInterface-8    	     300	   4150110 ns/op	    4703 B/op	      89 allocs/op
BenchmarkRedisCacheWithStruct-8          	     300	   4274845 ns/op	    6241 B/op	     111 allocs/op
BenchmarkSqlCacheWithInt-8               	     200	   8440496 ns/op	    8406 B/op	     221 allocs/op
BenchmarkSqlCacheWithString-8            	     200	   8506969 ns/op	    8291 B/op	     212 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	     100	  14260108 ns/op	   14950 B/op	     455 allocs/op
BenchmarkSqlCacheWithStruct-8            	     200	   8682163 ns/op	   18398 B/op	     547 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


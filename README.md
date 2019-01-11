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
BenchmarkMemoryCacheWithInt-8            	 1000000	      1298 ns/op	     128 B/op	       7 allocs/op
BenchmarkMemoryCacheWithString-8         	 1000000	      1432 ns/op	     144 B/op	       7 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	 1000000	      1832 ns/op	     544 B/op	      11 allocs/op
BenchmarkMemoryCacheWithStruct-8         	 1000000	      1833 ns/op	     688 B/op	      12 allocs/op
BenchmarkFileCacheWithInt-8              	   10000	    108644 ns/op	    2954 B/op	      35 allocs/op
BenchmarkFileCacheWithString-8           	   10000	    109820 ns/op	    2968 B/op	      36 allocs/op
BenchmarkFileCacheWithMapInterface-8     	   10000	    130769 ns/op	    5000 B/op	      59 allocs/op
BenchmarkFileCacheWithStruct-8           	   10000	    125431 ns/op	    5910 B/op	      64 allocs/op
BenchmarkRedisCacheWithInt-8             	    5000	    342156 ns/op	    1122 B/op	      33 allocs/op
BenchmarkRedisCacheWithString-8          	    5000	    321191 ns/op	    1132 B/op	      34 allocs/op
BenchmarkRedisCacheWithMapInterface-8    	    5000	    334021 ns/op	    3212 B/op	      58 allocs/op
BenchmarkRedisCacheWithStruct-8          	    3000	    334561 ns/op	    4101 B/op	      63 allocs/op
BenchmarkSqlCacheWithInt-8               	    1000	   1849025 ns/op	    5195 B/op	     144 allocs/op
BenchmarkSqlCacheWithString-8            	    1000	   1884996 ns/op	    5106 B/op	     135 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	    1000	   1872619 ns/op	   11082 B/op	     373 allocs/op
BenchmarkSqlCacheWithStruct-8            	    1000	   2117309 ns/op	   13910 B/op	     451 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


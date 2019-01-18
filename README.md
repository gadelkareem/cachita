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
BenchmarkMemoryCacheWithInt-8            	 1000000	      1218 ns/op	     120 B/op	       6 allocs/op
BenchmarkMemoryCacheWithString-8         	 1000000	      1234 ns/op	     136 B/op	       6 allocs/op
BenchmarkMemoryCacheWithMapInterface-8   	 1000000	      1445 ns/op	     536 B/op	      10 allocs/op
BenchmarkMemoryCacheWithStruct-8         	 1000000	      1588 ns/op	     680 B/op	      11 allocs/op
BenchmarkMemory_Incr-8                   	  500000	      2389 ns/op	     192 B/op	      10 allocs/op
BenchmarkFileCacheWithInt-8              	   10000	    110629 ns/op	    2946 B/op	      34 allocs/op
BenchmarkFileCacheWithString-8           	   10000	    117502 ns/op	    2968 B/op	      35 allocs/op
BenchmarkFileCacheWithMapInterface-8     	   10000	    121150 ns/op	    4998 B/op	      58 allocs/op
BenchmarkFileCacheWithStruct-8           	   10000	    120383 ns/op	    5909 B/op	      63 allocs/op
BenchmarkFile_Incr-8                     	   10000	    188167 ns/op	    7095 B/op	      74 allocs/op
BenchmarkRedisCacheWithInt-8             	    5000	    331572 ns/op	     703 B/op	      25 allocs/op
BenchmarkRedisCacheWithString-8          	    5000	    351982 ns/op	    1202 B/op	      35 allocs/op
BenchmarkRedisCacheWithMapInterface-8    	    5000	    331931 ns/op	    3284 B/op	      59 allocs/op
BenchmarkRedisCacheWithStruct-8          	    5000	    336453 ns/op	    4184 B/op	      64 allocs/op
BenchmarkRedis_Incr-8                    	    2000	    774163 ns/op	    1598 B/op	      45 allocs/op
BenchmarkSqlCacheWithInt-8               	    1000	   2468703 ns/op	    5168 B/op	     143 allocs/op
BenchmarkSqlCacheWithString-8            	    1000	   2121222 ns/op	    5121 B/op	     135 allocs/op
BenchmarkSqlCacheWithMapInterface-8      	    1000	   2838557 ns/op	   11137 B/op	     373 allocs/op
BenchmarkSqlCacheWithStruct-8            	    1000	   1903278 ns/op	   13880 B/op	     450 allocs/op
BenchmarkSql_Incr-8                      	     500	   3175832 ns/op	    9693 B/op	     268 allocs/op
```

## Howto

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


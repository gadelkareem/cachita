# Cachita
Cachita is a golang file, memory, SQL, Redis cache library 

[![Build Status](https://github.com/gadelkareem/cachita/actions/workflows/go.yml/badge.svg)](https://github.com/gadelkareem/cachita/actions)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/github.com/gadelkareem/cachita)

- Simple caching with auto type assertion included.
- In memory file cache index to avoid unneeded I/O.
- [Msgpack](https://msgpack.org/index.html) based binary serialization using [msgpack](https://github.com/vmihailenco/msgpack) library for file caching.
- [radix](https://github.com/mediocregopher/radix) Redis client.
- Tag cache and invalidate cache keys based on tags, check in the [examples](./example_test.go).


API docs: https://pkg.go.dev/github.com/gadelkareem/cachita.

Examples: [examples](./example_test.go).

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
		// do something
	}

	var holder string
	err = cache.Get("cache_key", &holder)
	if err != nil && err != cachita.ErrNotFound {
		panic(err)
	}

	fmt.Printf("%s", holder) // prints "some data"

	err = cache.Invalidate("cache_key")
	if err != nil {
		panic(err)
	}

	// Output: some data

}

```

## Benchmark

```
> go test -v -bench=. -benchmem
BenchmarkFileCacheWithInt
BenchmarkFileCacheWithInt-8              	   12104	    102738 ns/op	    2446 B/op	      31 allocs/op
BenchmarkFileCacheWithString
BenchmarkFileCacheWithString-8           	   10000	    104694 ns/op	    2469 B/op	      32 allocs/op
BenchmarkFileCacheWithMapInterface
BenchmarkFileCacheWithMapInterface-8     	    9319	    112302 ns/op	    4499 B/op	      55 allocs/op
BenchmarkFileCacheWithStruct
BenchmarkFileCacheWithStruct-8           	   10000	    118305 ns/op	    5407 B/op	      60 allocs/op
BenchmarkFile_Incr
BenchmarkFile_Incr-8                     	    7056	    225759 ns/op	    3052 B/op	      41 allocs/op
BenchmarkFile_Tag
BenchmarkFile_Tag-8                      	    4236	    316565 ns/op	    2802 B/op	      51 allocs/op
BenchmarkMemoryCacheWithInt
BenchmarkMemoryCacheWithInt-8            	  873777	      1145 ns/op	     120 B/op	       6 allocs/op
BenchmarkMemoryCacheWithString
BenchmarkMemoryCacheWithString-8         	  866904	      1185 ns/op	     136 B/op	       6 allocs/op
BenchmarkMemoryCacheWithMapInterface
BenchmarkMemoryCacheWithMapInterface-8   	  838364	      1505 ns/op	     536 B/op	      10 allocs/op
BenchmarkMemoryCacheWithStruct
BenchmarkMemoryCacheWithStruct-8         	  790718	      1445 ns/op	     680 B/op	      11 allocs/op
BenchmarkMemory_Incr
BenchmarkMemory_Incr-8                   	  731803	      1582 ns/op	     128 B/op	       7 allocs/op
BenchmarkMemory_Tag
BenchmarkMemory_Tag-8                    	  349209	      3465 ns/op	     513 B/op	      19 allocs/op
BenchmarkRedisCacheWithInt
BenchmarkRedisCacheWithInt-8             	    1461	    838600 ns/op	     496 B/op	      21 allocs/op
BenchmarkRedisCacheWithString
BenchmarkRedisCacheWithString-8          	    1593	    765341 ns/op	     995 B/op	      32 allocs/op
BenchmarkRedisCacheWithMapInterface
BenchmarkRedisCacheWithMapInterface-8    	    1435	    755585 ns/op	    3071 B/op	      55 allocs/op
BenchmarkRedisCacheWithStruct
BenchmarkRedisCacheWithStruct-8          	    1506	    821237 ns/op	    3963 B/op	      61 allocs/op
BenchmarkRedis_Incr
BenchmarkRedis_Incr-8                    	    1051	   1042468 ns/op	    1237 B/op	      32 allocs/op
BenchmarkRedis_Tag
BenchmarkRedis_Tag-8                     	     452	   2752817 ns/op	    3509 B/op	     117 allocs/op
BenchmarkSqlCacheWithInt
BenchmarkSqlCacheWithInt-8               	     253	   7927815 ns/op	    4984 B/op	     118 allocs/op
BenchmarkSqlCacheWithString
BenchmarkSqlCacheWithString-8            	       1	1026688775 ns/op	   38224 B/op	     620 allocs/op
BenchmarkSqlCacheWithMapInterface
BenchmarkSqlCacheWithMapInterface-8      	     331	   4040482 ns/op	   10741 B/op	     345 allocs/op
BenchmarkSqlCacheWithStruct
BenchmarkSqlCacheWithStruct-8            	     250	   4357675 ns/op	   14142 B/op	     429 allocs/op
BenchmarkSql_Incr
BenchmarkSql_Incr-8                      	     234	   5412791 ns/op	    5812 B/op	     150 allocs/op
BenchmarkSql_Tag
BenchmarkSql_Tag-8                       	      82	  14308321 ns/op	   15407 B/op	     366 allocs/op
PASS
ok  	github.com/gadelkareem/cachita	41.180s
```

## How to

Please go through [examples](./example_test.go) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


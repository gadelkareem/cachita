# Cachita
Cachita is a golang file, memory, SQL, Redis cache library 

[![Build Status](https://travis-ci.org/gadelkareem/cachita.svg)](https://travis-ci.org/gadelkareem/cachita)
[![GoDoc](https://godoc.org/github.com/gadelkareem/cachita?status.svg)](https://godoc.org/github.com/gadelkareem/cachita)

- Simple caching with auto type assertion included.
- In memory file cache index to avoid unneeded I/O.
- [Msgpack](https://msgpack.org/index.html) based binary serialization using [msgpack](https://github.com/vmihailenco/msgpack) library for file caching.
- [radix](https://github.com/mediocregopher/radix) Redis client.
- Tag cache and invalidate cache keys based on tags, check in the [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples).


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
BenchmarkFileCacheWithInt-8              	   10000	    116118 ns/op	    2447 B/op	      31 allocs/op
BenchmarkFileCacheWithString
BenchmarkFileCacheWithString-8           	   10909	    123491 ns/op	    2470 B/op	      32 allocs/op
BenchmarkFileCacheWithMapInterface
BenchmarkFileCacheWithMapInterface-8     	    9862	    124641 ns/op	    4499 B/op	      55 allocs/op
BenchmarkFileCacheWithStruct
BenchmarkFileCacheWithStruct-8           	    9356	    130355 ns/op	    5404 B/op	      60 allocs/op
BenchmarkFile_Incr
BenchmarkFile_Incr-8                     	    6331	    192199 ns/op	    3113 B/op	      44 allocs/op
BenchmarkFile_Tag
BenchmarkFile_Tag-8                      	    3885	    286273 ns/op	    2720 B/op	      47 allocs/op
BenchmarkMemoryCacheWithInt
BenchmarkMemoryCacheWithInt-8            	  870573	      1288 ns/op	     120 B/op	       6 allocs/op
BenchmarkMemoryCacheWithString
BenchmarkMemoryCacheWithString-8         	  938899	      1161 ns/op	     136 B/op	       6 allocs/op
BenchmarkMemoryCacheWithMapInterface
BenchmarkMemoryCacheWithMapInterface-8   	  835402	      1618 ns/op	     536 B/op	      10 allocs/op
BenchmarkMemoryCacheWithStruct
BenchmarkMemoryCacheWithStruct-8         	  771076	      1591 ns/op	     680 B/op	      11 allocs/op
BenchmarkMemory_Incr
BenchmarkMemory_Incr-8                   	  649772	      1784 ns/op	     184 B/op	       9 allocs/op
BenchmarkMemory_Tag
BenchmarkMemory_Tag-8                    	  361974	      3458 ns/op	     439 B/op	      14 allocs/op
BenchmarkRedisCacheWithInt
BenchmarkRedisCacheWithInt-8             	    1404	    787836 ns/op	     492 B/op	      21 allocs/op
BenchmarkRedisCacheWithString
BenchmarkRedisCacheWithString-8          	    1573	    775092 ns/op	     995 B/op	      32 allocs/op
BenchmarkRedisCacheWithMapInterface
BenchmarkRedisCacheWithMapInterface-8    	    1506	    709349 ns/op	    3074 B/op	      55 allocs/op
BenchmarkRedisCacheWithStruct
BenchmarkRedisCacheWithStruct-8          	    1714	    872728 ns/op	    3969 B/op	      61 allocs/op
BenchmarkRedis_Incr
BenchmarkRedis_Incr-8                    	    1153	   1096139 ns/op	    1235 B/op	      32 allocs/op
BenchmarkRedis_Tag
BenchmarkRedis_Tag-8                     	     379	   3356175 ns/op	    8325 B/op	     201 allocs/op
BenchmarkSqlCacheWithInt
BenchmarkSqlCacheWithInt-8               	     277	   3960950 ns/op	    4741 B/op	     115 allocs/op
BenchmarkSqlCacheWithString
BenchmarkSqlCacheWithString-8            	     280	   3979248 ns/op	    4679 B/op	     106 allocs/op
BenchmarkSqlCacheWithMapInterface
BenchmarkSqlCacheWithMapInterface-8      	     282	   4816726 ns/op	   11444 B/op	     352 allocs/op
BenchmarkSqlCacheWithStruct
BenchmarkSqlCacheWithStruct-8            	     230	   4375050 ns/op	   13730 B/op	     425 allocs/op
BenchmarkSql_Incr
BenchmarkSql_Incr-8                      	     199	   6042507 ns/op	    6220 B/op	     154 allocs/op
BenchmarkSql_Tag
BenchmarkSql_Tag-8                       	      57	  35618536 ns/op	  836763 B/op	     967 allocs/op
PASS
ok  	github.com/gadelkareem/cachita	40.188s
```

## How to

Please go through [examples](https://godoc.org/github.com/gadelkareem/cachita#pkg-examples) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


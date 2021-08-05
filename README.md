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
BenchmarkFileCacheWithInt-8              	    8966	    137991 ns/op	    2400 B/op	      37 allocs/op
BenchmarkFileCacheWithString
BenchmarkFileCacheWithString-8           	    8258	    153027 ns/op	    2417 B/op	      38 allocs/op
BenchmarkFileCacheWithMapInterface
BenchmarkFileCacheWithMapInterface-8     	    6908	    157232 ns/op	    4238 B/op	      60 allocs/op
BenchmarkFileCacheWithStruct
BenchmarkFileCacheWithStruct-8           	    7179	    165023 ns/op	    4992 B/op	      65 allocs/op
BenchmarkFile_Incr
BenchmarkFile_Incr-8                     	    4814	    279499 ns/op	    3004 B/op	      47 allocs/op
BenchmarkFile_Tag
BenchmarkFile_Tag-8                      	    3506	    324306 ns/op	    3028 B/op	      71 allocs/op
BenchmarkMemoryCacheWithInt
BenchmarkMemoryCacheWithInt-8            	  954620	      1340 ns/op	     112 B/op	       6 allocs/op
BenchmarkMemoryCacheWithString
BenchmarkMemoryCacheWithString-8         	  774190	      1333 ns/op	     128 B/op	       6 allocs/op
BenchmarkMemoryCacheWithMapInterface
BenchmarkMemoryCacheWithMapInterface-8   	  614234	      1850 ns/op	     520 B/op	      10 allocs/op
BenchmarkMemoryCacheWithStruct
BenchmarkMemoryCacheWithStruct-8         	  629415	      2244 ns/op	     664 B/op	      11 allocs/op
BenchmarkMemory_Incr
BenchmarkMemory_Incr-8                   	  527721	      2201 ns/op	      96 B/op	       4 allocs/op
BenchmarkMemory_Tag
BenchmarkMemory_Tag-8                    	  197188	      5393 ns/op	     598 B/op	      24 allocs/op
BenchmarkRedisCacheWithInt
BenchmarkRedisCacheWithInt-8             	    1591	    644581 ns/op	     460 B/op	      21 allocs/op
BenchmarkRedisCacheWithString
BenchmarkRedisCacheWithString-8          	    1770	    805753 ns/op	     945 B/op	      32 allocs/op
BenchmarkRedisCacheWithMapInterface
BenchmarkRedisCacheWithMapInterface-8    	    2138	    472726 ns/op	    2988 B/op	      54 allocs/op
BenchmarkRedisCacheWithStruct
BenchmarkRedisCacheWithStruct-8          	    2751	    475874 ns/op	    3876 B/op	      59 allocs/op
BenchmarkRedis_Incr
BenchmarkRedis_Incr-8                    	    1486	    826275 ns/op	    1201 B/op	      32 allocs/op
BenchmarkRedis_Tag
BenchmarkRedis_Tag-8                     	     660	   1822219 ns/op	    3309 B/op	     122 allocs/op
BenchmarkSqlCacheWithInt
BenchmarkSqlCacheWithInt-8               	     288	   4125553 ns/op	    4530 B/op	     111 allocs/op
BenchmarkSqlCacheWithString
BenchmarkSqlCacheWithString-8            	     302	   3839348 ns/op	    4373 B/op	     101 allocs/op
BenchmarkSqlCacheWithMapInterface
BenchmarkSqlCacheWithMapInterface-8      	     271	   4057435 ns/op	   10359 B/op	     339 allocs/op
BenchmarkSqlCacheWithStruct
BenchmarkSqlCacheWithStruct-8            	     286	   4106997 ns/op	   13065 B/op	     417 allocs/op
BenchmarkSql_Incr
BenchmarkSql_Incr-8                      	     219	   5618603 ns/op	    5771 B/op	     145 allocs/op
BenchmarkSql_Tag
BenchmarkSql_Tag-8                       	      91	  13543030 ns/op	   14772 B/op	     354 allocs/op
PASS
ok  	github.com/gadelkareem/cachita	40.686s
```

## How to

Please go through [examples](./example_test.go) to get an idea how to use this package.

## See also

- [Golang Helpers](https://github.com/gadelkareem/go-helpers)


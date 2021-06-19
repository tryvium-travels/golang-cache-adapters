<p align="center"><img src="https://res.cloudinary.com/tryvium/image/upload/v1551645701/company/logo-circle.png"/></p>

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tryvium-travels/golang-cache-adapters?style=flat-square)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/tryvium-travels/golang-cache-adapters)
[![Go Report Card](https://goreportcard.com/badge/github.com/saniales/golang-crypto-trading-bot?style=flat-square)](https://goreportcard.com/report/github.com/tryvium-travels/golang-cache-adapters)
![GitHub](https://img.shields.io/github/license/tryvium-travels/golang-cache-adapters?style=flat-square)
![Twitter Follow](https://img.shields.io/twitter/follow/tryviumtravels?style=social)
[![Build and Test library](https://github.com/tryvium-travels/golang-cache-adapters/actions/workflows/test-library.yml/badge.svg?style=flat-square)](https://github.com/tryvium-travels/golang-cache-adapters/actions/workflows/test-library.yml)

# Golang Cache Adapters
A set of Cache Adapters for various distributed systems (like Redis) written in Go.

Now with a new [`MultiCacheAdapter`](/multicache) implementation!

## Supported CacheAdapter implementations

- [**MultiCache**](/multicache) -> Leverages the possibility to use multiple cache adapters at the same time, useful to create fallbacks in case one or more of the cache service you specified gives temporary errors.
- [**Redis**](/redis) -> using [`github.com/gomodule/redigo`](github.com/gomodule/redigo)

## Library reference

Just check the [**GoDocs**](https://pkg.go.dev/github.com/tryvium-travels/golang-cache-adapters)

# Install

Just use the standard way

``` bash
go get github.com/tryvium-travels/golang-cache-adapters
```

or use go modules, it should take care of this automagically.

# Usage

The `CacheAdapter` interface offers 2 main ways to access a cache: `Get` and `Set` methods.

You can use `Delete` and `SetTTL` functions as well. For more info, check the docs.

This example creates a new `RedisAdapter` and uses it, but you can replace it with any of the other
supported Adapters.

To know how to use the [`MultiCacheAdapter`](/multicache) check out the docs.

``` go
package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
)

func main() {
	redisURI := "rediss://my-redis-instance-uri:port"

	myRedisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// obtain a redis connection, there
			// are plenty of ways to do that.
			return redis.DialURL(redisURI)
		},
	}

	exampleTTL := time.Hour

	adapter, err := rediscacheadapters.New(myRedisPool, exampleTTL)
	if err != nil {
		// remember to check for errors
		log.Fatalf("Adapter initialization error: %s", err)
	}

	type exampleStruct struct {
		Value string
	}

	exampleKey := "a:redis:key"

	var exampleValue exampleStruct
	err = adapter.Get(exampleKey, &exampleValue)
	if err != nil {
		// remember to check for errors
		log.Fatalf("adapter.Get error: %s", err)
	}

	exampleKey = "another:redis:key"

	// nil TTL represents the default value put in the New function
	err = adapter.Set(exampleKey, exampleValue, nil)
	if err != nil {
		// remember to check for errors
		log.Fatalf("adapter.Get error: %s", err)
	}
}
```

# Contributing

First of all, thank you!

To contribute to our project:

1. Open an Issue signaling bugs or proposals using the provided templates
2. If you are going to add a new adapter, please stick to the *NAMING CONVENTION* for branch, package, folder, file and adapter names.
3. We require a total test coverage (100%) if we are going to accept a new adapter, we can help you if you have any difficulties to achieve that.
4. Please be polite in issue and PR discussions and do not offend other people.

## Naming Convention

The following naming convention applies if you want to contribute:

1. Branch name should follow the following form

   `adapters/{{ADAPTER_NAME}}-#{{GITHUB_ISSUE_ID}}`
   
   Example branch name: `adapters/redis-#1`
2. If you add a new adapter: 
   
   The ***name of the folder*** must match the ***lowercase adapter name*** (e.g. redis/mongodb/memcached)

   We strongly suggest you to create an adapter and a session_adapter go file, although it is not required.

   Example file names:
   
   `redis_adapter.go` and `redis_session_adapter.go`
   
   Example test file names:
   
   `redis_adapter_test.go` and `redis_session_adapter_test.go`.

   There must also be test files with 100% test coverage.
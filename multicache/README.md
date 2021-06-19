<p align="center"><img src="https://res.cloudinary.com/tryvium/image/upload/v1551645701/company/logo-circle.png"/></p>

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tryvium-travels/golang-cache-adapters?style=flat-square)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/tryvium-travels/golang-cache-adapters)
[![Go Report Card](https://goreportcard.com/badge/github.com/saniales/golang-crypto-trading-bot?style=flat-square)](https://goreportcard.com/report/github.com/tryvium-travels/golang-cache-adapters)
![GitHub](https://img.shields.io/github/license/tryvium-travels/golang-cache-adapters?style=flat-square)
![Twitter Follow](https://img.shields.io/twitter/follow/tryviumtravels?style=social)

# Multi Cache Adapter implementation
A `CacheAdapter` implementation that allows to connect and use multiple adapters at the same time.

## Features

- Allows the creation of fallback cache logics, allowing to reach different cache types if one of them is failing for any reason
- Allows, optionally, to track partial failures as warnings, so you can log them using your preferred method

## Usage

Since `MultiCacheAdapter` and `MultiCacheSessionAdapter` implement the respective interfaces (`CacheAdapter` and `CacheSessionAdapter`) you have all the methods at your disposal.

Please refer to the following example for the correct usage:

``` go
package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
    multicacheadapters "github.com/tryvium-travels/golang-cache-adapters/multicache"
)

func main() {
    // Ideally you would want to use 2 DIFFERENT cache types
    // for the sake of the simplicity, we will just create 2
    // adapters of the same type.

	redisURI1 := "rediss://my-redis-instance-uri1:port"
    redisURI2 := "rediss://my-redis-instance-uri2:port"

	myRedisPool1 := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// obtain a redis connection, there
			// are plenty of ways to do that.
			return redis.DialURL(redisURI1)
		},
	}

    myRedisPool2 := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// obtain a redis connection, there
			// are plenty of ways to do that.
			return redis.DialURL(redisURI2)
		},
	}

	exampleTTL := time.Hour

    // first we create the adapters.
	redisAdapter1, err := rediscacheadapters.New(myRedisPool1, exampleTTL)
	if err != nil {
		// remember to check for errors
		log.Fatalf("Redis Adapter 1 initialization error: %s", err)
	}

    redisAdapter2, err := rediscacheadapters.New(myRedisPool1, exampleTTL)
	if err != nil {
		// remember to check for errors
		log.Fatalf("Redis Adapter 2 initialization error: %s", err)
	}

    // then we create a MultiCacheAdapter to use them at the same time.
    adapter, err := multicacheadapters.New(redisAdapter1, redisAdapter2)
    if err != nil {
		// remember to check for errors
		log.Fatalf("Multi Cache Adapter 2 initialization error: %s", err)
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
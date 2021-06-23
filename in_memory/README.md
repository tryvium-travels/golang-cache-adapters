<p align="center"><img src="https://res.cloudinary.com/tryvium/image/upload/v1551645701/company/logo-circle.png"/></p>

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tryvium-travels/golang-cache-adapters?style=flat-square)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/tryvium-travels/golang-cache-adapters)
[![Go Report Card](https://goreportcard.com/badge/github.com/saniales/golang-crypto-trading-bot?style=flat-square)](https://goreportcard.com/report/github.com/tryvium-travels/golang-cache-adapters)
![GitHub](https://img.shields.io/github/license/tryvium-travels/golang-cache-adapters?style=flat-square)
![Twitter Follow](https://img.shields.io/twitter/follow/tryviumtravels?style=social)

# Cache Adapter implementation for In Memory

A `CacheAdapter` implementation that allows to use a RAM-based memory.

Beware that at the moment it has the [**github.com/gomodule/redigo**](github.com/gomodule/redigo) dependency

## Usage

Please refer to the following example for the correct usage:

``` go
package main

import (
	"log"
	"time"

	inmemorycacheadapters "github.com/tryvium-travels/golang-cache-adapters/in_memory"
)

func main() {
	exampleTTL := time.Hour

	adapter, err := inmemorycacheadapters.New(exampleTTL)
	if err != nil {
		// remember to check for errors
		log.Fatalf("Adapter initialization error: %s", err)
	}

	type exampleStruct struct {
		Value string
	}

	exampleKey := "a:in_memory:key"

	var exampleValue exampleStruct
	err = adapter.Get(exampleKey, &exampleValue)
	if err != nil {
		// remember to check for errors
		log.Fatalf("adapter.Get error: %s", err)
	}

	exampleKey = "another:in_memory:key"

	// nil TTL represents the default value put in the New function
	err = adapter.Set(exampleKey, exampleValue, nil)
	if err != nil {
		// remember to check for errors
		log.Fatalf("adapter.Get error: %s", err)
	}
}
```
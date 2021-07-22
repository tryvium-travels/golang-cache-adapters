// Copyright 2021 Tryvium Travels LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rediscacheadapters

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// RedisAdapter is the CacheAdapter implementation for Redis.
type RedisAdapter struct {
	pool       *redis.Pool   // The Redis pool used to create connections.
	defaultTTL time.Duration // The defaultTTL of the Set operations.
}

// New creates a new RedisAdapter from an initialized Redis pool.
func New(pool *redis.Pool, defaultTTL time.Duration) (cacheadapters.CacheAdapter, error) {
	if pool == nil {
		return nil, fmt.Errorf("the Redis Pool cannot be nil")
	}

	if defaultTTL <= 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &RedisAdapter{
		pool:       pool,
		defaultTTL: defaultTTL,
	}, nil
}

// OpenSession opens a new Cache Session.
func (ra *RedisAdapter) OpenSession() (cacheadapters.CacheSessionAdapter, error) {
	conn, err := ra.pool.Dial()
	if err != nil {
		return nil, err
	}

	return NewSession(conn, ra.defaultTTL)
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (ra *RedisAdapter) Get(key string, objectRef interface{}) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Get(key, objectRef)
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (ra *RedisAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Set(key, object, TTL)
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (ra *RedisAdapter) SetTTL(key string, newTTL time.Duration) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.SetTTL(key, newTTL)
}

// Delete deletes a key from the cache.
func (ra *RedisAdapter) Delete(key string) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Delete(key)
}

// Copyright 2023 Tryvium Travels LTD
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

package cacheadapters

import (
	"time"
)

// TTLExpired represents the zero-value of a time expiration.
const TTLExpired time.Duration = 0

// CacheAdapter represents a Cache Mechanism abstraction.
type CacheAdapter interface {
	// OpenSession opens a new Cache Session.
	OpenSession() (CacheSessionAdapter, error)

	cacheOperator
}

// CacheSessionAdapter represents a Cache Session Mechanism abstraction.
type CacheSessionAdapter interface {
	// Close closes the Cache Session.
	Close() error

	cacheOperator
}

// cacheOperator is an intermediary interface to share methods between CacheAdapter and CacheSessionAdapter
type cacheOperator interface {
	// Get obtains a value from the cache using a key, then tries to unmarshal
	// it into the object reference passed as parameter.
	Get(key string, objectRef interface{}) error

	// Set sets a value represented by the object parameter into the cache,
	// with the specified key.
	Set(key string, object interface{}, TTL *time.Duration) error

	// SetTTL marks the specified key new expiration, deletes it via using
	// cacheadapters.TTLExpired or negative duration.
	SetTTL(key string, newTTL time.Duration) error

	// Delete deletes a key from the cache.
	Delete(key string) error
}

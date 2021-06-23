// Copyright 2021 The Tryvium Company LTD
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

package inmemorycacheadapters

import (
	"encoding/json"
	"sync"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// cacheItem is the internal struct
// handling the mechanism of cache expiration.
type cacheItem struct {
	item      json.RawMessage // The actual item in cache.
	expiresAt time.Time       // The expiration time of the item in cache.
}

// cacheData is the container of all the in-memory
// cache used by the adapter.
type cacheData map[string]cacheItem

// InMemoryAdapter is the cache adapter which uses internal memory
// of the process.
type InMemoryAdapter struct {
	defaultTTL time.Duration // The defaultTTL of the Set operations.
	data       cacheData     // The data being stored in the in-memory cache.
	mutex      sync.Mutex    // The mutex locking the operations.
}

// New creates a new InMemoryAdapter from an default TTL.
func New(defaultTTL time.Duration) (cacheadapters.CacheAdapter, error) {
	if defaultTTL <= 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &InMemoryAdapter{
		defaultTTL: defaultTTL,
		data:       make(cacheData),
	}, nil
}

// OpenSession opens a new Cache Session.
// Returns the same adapter because the
// session with the memory is always open.
func (ima *InMemoryAdapter) OpenSession() (cacheadapters.CacheSessionAdapter, error) {
	return ima, nil
}

// Close closes the Cache Session.
// Returns nil because the session with
// the memory is always on and does not
// need to be closed.
func (ima *InMemoryAdapter) Close() error {
	return nil
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (ima *InMemoryAdapter) Get(key string, resultRef interface{}) error {
	if resultRef == nil {
		return cacheadapters.ErrGetRequiresObjectReference
	}

	ima.mutex.Lock()
	valueFromMemory, exists := ima.data[key]
	ima.mutex.Unlock()
	if !exists {
		return cacheadapters.ErrNotFound
	}

	now := time.Now()
	if valueFromMemory.expiresAt.UnixNano() < now.UnixNano() {
		ima.Delete(key)
		return cacheadapters.ErrNotFound
	}

	err := json.Unmarshal(valueFromMemory.item, resultRef)
	if err != nil {
		return err
	}

	return nil
}

// Set sets a value represented by the object parameter into the cache,
// with the specified key.
func (ima *InMemoryAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	if TTL == nil {
		TTL = new(time.Duration)
		*TTL = ima.defaultTTL
	} else if *TTL <= 0 {
		return cacheadapters.ErrInvalidTTL
	}

	now := time.Now()
	expiresAt := now.Add(*TTL)

	content, err := json.Marshal(object)
	if err != nil {
		return err
	}

	ima.mutex.Lock()
	ima.data[key] = cacheItem{
		item:      content,
		expiresAt: expiresAt,
	}
	ima.mutex.Unlock()

	return nil
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (ima *InMemoryAdapter) SetTTL(key string, newTTL time.Duration) error {
	if newTTL <= cacheadapters.TTLExpired {
		return ima.Delete(key)
	}

	ima.mutex.Lock()
	valueFromMemory, exists := ima.data[key]
	ima.mutex.Unlock()
	if !exists {
		return cacheadapters.ErrNotFound
	}

	now := time.Now()
	newExpiresAt := now.Add(newTTL)

	if valueFromMemory.expiresAt.UnixNano() < now.UnixNano() {
		return ima.Delete(key)
	}

	valueFromMemory.expiresAt = newExpiresAt

	ima.mutex.Lock()
	ima.data[key] = valueFromMemory
	ima.mutex.Unlock()

	return nil
}

// Delete deletes a key from the cache.
func (ima *InMemoryAdapter) Delete(key string) error {
	ima.mutex.Lock()
	delete(ima.data, key)
	ima.mutex.Unlock()
	return nil
}

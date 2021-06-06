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

package cacheadapters

import "time"

// MultiCacheAdapter is a cache adapter which uses multiple
// sub-adapters, following a priority given by the index of
// the adapter in the inner array of adapters.
type MultiCacheAdapter struct {
	subAdapters []CacheAdapter // The array of sub-adapters
}

// NewMultiCacheAdapter creates a new multi cache adapter from an
// index-based priority array of cache adapters (called sub-adapters).
//
//     index-based means that the array at the first position(s) will
//     have more priority than those at latter positions.
func NewMultiCacheAdapter(adapters ...CacheAdapter) (*MultiCacheAdapter, error) {
	for _, adapter := range adapters {
		if adapter == nil {
			return nil, ErrNilSubAdapter
		}
	}
	return &MultiCacheAdapter{adapters}, nil
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (mca *MultiCacheAdapter) Get(key string, objectRef interface{}) error {
	return errNotImplemented
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (mca *MultiCacheAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	return errNotImplemented
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (mca *MultiCacheAdapter) SetTTL(key string, newTTL time.Duration) error {
	return errNotImplemented
}

// Delete deletes a key from the cache.
func (mca *MultiCacheAdapter) Delete(key string) error {
	return errNotImplemented
}

// InTransaction allows to execute multiple Cache Sets and Gets in a Transaction, then tries to
// Unmarshal the array of results into the specified array of object references.
func (mca *MultiCacheAdapter) InTransaction(inTransactionFunc InTransactionFunc, objectRefs []interface{}) error {
	return errNotImplemented
}

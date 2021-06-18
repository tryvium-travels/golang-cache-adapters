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

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
)

// MultiCacheAdapter is a cache adapter which uses multiple
// sub-adapters, following a priority given by the index of
// the adapter in the inner array of adapters.
type MultiCacheAdapter struct {
	subAdapters  []CacheAdapter // The array of sub-adapters
	showWarnings bool
	wg           sync.WaitGroup
}

// NewMultiCacheAdapter creates a new multi cache adapter from an
// index-based priority array of cache adapters (called sub-adapters) and a
// flag instructing to show warning (non-fatal) errors.
//
//     index-based means that the array at the first position(s) will
//     have more priority than those at latter positions.
func NewMultiCacheAdapter(adapters ...CacheAdapter) (*MultiCacheAdapter, error) {
	for _, adapter := range adapters {
		if adapter == nil {
			return nil, ErrNilSubAdapter
		}
	}

	return &MultiCacheAdapter{adapters, false, sync.WaitGroup{}}, nil
}

// EnableWarning enable the return of warning errors.
//
// If the error is a warning, you can continue standard execution as
// The operations concluded successfully. You can then log the warning
// using your favourite tool (like sentry).
//
//     Warning errors, if shown need different handling from traditional
//     errors. Use the helper IsWarning(error err) to check for warnings.
//
// Example of handling of warnings:
//
//     err := adapter.Get("key", &objRef)
//     if (cacheadapters.IsWarning(err)) {
//	       // log the error, but use objRef safely
//     } else if err != nil {
//         // log the error and handle a failure
//         // you cannot use objRef safely here
//     }
//     // else use objRef safely without any error
func (mca *MultiCacheAdapter) EnableWarnings() {
	mca.showWarnings = true
}

// DisableWarning disables the return of warning errors.
//
// If the error is a warning, you can continue standard execution as
// The operations concluded successfully. You can then log the warning
// using your favourite tool (like sentry).
//
//     Warning errors, if shown need different handling from traditional
//     errors. Use the helper IsWarning(error err) to check for warnings.
//
// Example of handling of warnings:
//
//     err := adapter.Get("key", &objRef)
//     if (errors.Is(cacheadapter.ErrWarning)) {
//	       // log the error, but use objRef safely
//     } else if err != nil {
//         // log the error and handle a failure
//         // you cannot use objRef safely here
//     }
//     // else use objRef safely without any error
func (mca *MultiCacheAdapter) DisableWarnings() {
	mca.showWarnings = false
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (mca *MultiCacheAdapter) Get(key string, objectRef interface{}) error {
	errs := make([]error, 0, len(mca.subAdapters))
	for _, adapter := range mca.subAdapters {
		var temp json.RawMessage

		err := adapter.Get(key, &temp)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = json.Unmarshal(temp, &objectRef)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		break
	}

	err := multierror.Append(nil, errs...)
	if len(errs) == len(mca.subAdapters) {
		return err
	}

	if mca.showWarnings {
		return multierror.Append(ErrMultiCacheWarning, err)
	}

	return nil
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (mca *MultiCacheAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	errs := make([]error, 0, len(mca.subAdapters))
	for _, adapter := range mca.subAdapters {
		err := adapter.Set(key, object, TTL)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	err := multierror.Append(nil, errs...)
	if len(errs) == len(mca.subAdapters) {
		return err
	}

	if mca.showWarnings {
		return multierror.Append(ErrMultiCacheWarning, err)
	}

	return nil
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (mca *MultiCacheAdapter) SetTTL(key string, newTTL time.Duration) error {
	errs := make([]error, 0, len(mca.subAdapters))
	for _, adapter := range mca.subAdapters {
		err := adapter.SetTTL(key, newTTL)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	err := multierror.Append(nil, errs...)
	if len(errs) == len(mca.subAdapters) {
		return err
	}

	if mca.showWarnings {
		return multierror.Append(ErrMultiCacheWarning, err)
	}

	return nil
}

// Delete deletes a key from the cache.
func (mca *MultiCacheAdapter) Delete(key string) error {
	errs := make([]error, 0, len(mca.subAdapters))
	for _, adapter := range mca.subAdapters {
		err := adapter.Delete(key)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	err := multierror.Append(nil, errs...)
	if len(errs) == len(mca.subAdapters) {
		return err
	}

	if mca.showWarnings {
		return multierror.Append(ErrMultiCacheWarning, err)
	}

	return nil
}

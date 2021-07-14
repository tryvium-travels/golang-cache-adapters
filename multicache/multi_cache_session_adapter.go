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

package multicacheadapters

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// MultiCacheSessionAdapter is a cache adapter which uses multiple
// sub-adapters, following a priority given by the index of
// the adapter in the inner array of adapters.
type MultiCacheSessionAdapter struct {
	subAdapters  []cacheadapters.CacheSessionAdapter // The array of sub-adapters
	showWarnings bool
	wg           sync.WaitGroup
}

// NewSession creates a new multi cache session adapter from an
// index-based priority array of cache adapters (called sub-adapters) and a
// flag instructing to show warning (non-fatal) errors.
//
//     index-based means that the array at the first position(s) will
//     have more priority than those at latter positions.
func NewSession(adapters ...cacheadapters.CacheSessionAdapter) (*MultiCacheSessionAdapter, error) {
	finalAdapters := make([]cacheadapters.CacheSessionAdapter, 0, len(adapters))
	for _, adapter := range adapters {
		if value := reflect.ValueOf(adapter); adapter != nil && value.IsValid() && !value.IsNil() {
			finalAdapters = append(finalAdapters, adapter)
		}
	}

	if len(finalAdapters) == 0 {
		return nil, ErrInvalidSubAdapters
	}

	return &MultiCacheSessionAdapter{finalAdapters, false, sync.WaitGroup{}}, nil
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
func (mcsa *MultiCacheSessionAdapter) EnableWarnings() {
	mcsa.showWarnings = true
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
func (mcsa *MultiCacheSessionAdapter) DisableWarnings() {
	mcsa.showWarnings = false
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (mcsa *MultiCacheSessionAdapter) Get(key string, objectRef interface{}) error {
	errs := make([]error, 0, len(mcsa.subAdapters))
	for _, adapter := range mcsa.subAdapters {
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

	return mcsa.errorOrNil(errs)
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (mcsa *MultiCacheSessionAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	errs := make([]error, 0, len(mcsa.subAdapters))
	for _, adapter := range mcsa.subAdapters {
		err := adapter.Set(key, object, TTL)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return mcsa.errorOrNil(errs)
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (mcsa *MultiCacheSessionAdapter) SetTTL(key string, newTTL time.Duration) error {
	errs := make([]error, 0, len(mcsa.subAdapters))
	for _, adapter := range mcsa.subAdapters {
		err := adapter.SetTTL(key, newTTL)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return mcsa.errorOrNil(errs)
}

// Delete deletes a key from the cache.
func (mcsa *MultiCacheSessionAdapter) Delete(key string) error {
	errs := make([]error, 0, len(mcsa.subAdapters))
	for _, adapter := range mcsa.subAdapters {
		err := adapter.Delete(key)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return mcsa.errorOrNil(errs)
}

// Close closes the Cache Sessions.
func (mcsa *MultiCacheSessionAdapter) Close() error {
	errs := make([]error, 0, len(mcsa.subAdapters))
	for _, adapter := range mcsa.subAdapters {
		err := adapter.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return mcsa.errorOrNil(errs)
}

// errorOrNil parses the accumulated errors into one final
// error, or returns nil if there are none.
func (mcsa *MultiCacheSessionAdapter) errorOrNil(errs []error) error {
	err := multierror.Append(nil, errs...)
	if len(errs) == len(mcsa.subAdapters) {
		return err
	}

	if mcsa.showWarnings {
		return multierror.Append(ErrMultiCacheWarning, err)
	}

	return nil
}

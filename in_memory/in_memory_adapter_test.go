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

package inmemorycacheadapters_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	inmemorycacheadapters "github.com/tryvium-travels/golang-cache-adapters/in_memory"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

type InMemoryAdapterTestSuite struct {
	*suite.Suite
	*testutil.CacheAdapterPartialTestSuite
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func testSleepFunc() func(time.Duration) {
	return func(duration time.Duration) {
		time.Sleep(duration)
	}
}

func newTestAdapterFunc(defaultTTL time.Duration) func() (cacheadapters.CacheAdapter, error) {
	return func() (cacheadapters.CacheAdapter, error) {
		return inmemorycacheadapters.New(defaultTTL)
	}
}

func newTestSessionFunc(t *testing.T, defaultTTL time.Duration) func() (cacheadapters.CacheSessionAdapter, error) {
	return func() (cacheadapters.CacheSessionAdapter, error) {
		cacheAdapter, err := newTestAdapterFunc(defaultTTL)()
		if err != nil {
			t.Error(err)
		}
		return cacheAdapter.OpenSession()
	}
}

// newInMemoryTestSuite creates a new test suite with tests for Redis adapters and sessions.
func newInMemoryTestSuite(t *testing.T, defaultTTL time.Duration) *InMemoryAdapterTestSuite {
	var suite suite.Suite

	return &InMemoryAdapterTestSuite{
		Suite: &suite,
		CacheAdapterPartialTestSuite: &testutil.CacheAdapterPartialTestSuite{
			Suite:      &suite,
			DefaultTTL: defaultTTL,
			NewAdapter: newTestAdapterFunc(defaultTTL),
			NewSession: newTestSessionFunc(t, defaultTTL),
			SleepFunc:  testSleepFunc(),
		},
	}
}

func (suite *InMemoryAdapterTestSuite) SetupSuite() {
	// actually, nothing is required for the In-Memory Cache Adapter
}

func (Test *InMemoryAdapterTestSuite) TearDownSuite() {
	// actually, nothing is required for the In-Memory Cache Adapter
}

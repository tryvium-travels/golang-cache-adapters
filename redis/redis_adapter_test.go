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

package rediscacheadapters_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

func TestRedisAdapterSuite(t *testing.T) {
	defaultTTL := 1 * time.Second
	suite.Run(t, newRedisTestSuite(t, defaultTTL))
}

type RedisAdapterTestSuite struct {
	*suite.Suite
	*testutil.CacheAdapterPartialTestSuite
}

func newTestAdapterFunc(defaultTTL time.Duration) func() (cacheadapters.CacheAdapter, error) {
	return func() (cacheadapters.CacheAdapter, error) {
		return rediscacheadapters.New(testRedisPool, defaultTTL)
	}
}

func testSleepFunc() func(time.Duration) {
	return func(duration time.Duration) {
		localRedisServer.FastForward(duration)
	}
}

// newRedisTestSuite creates a new test suite with tests for Redis adapters and sessions.
func newRedisTestSuite(t *testing.T, defaultTTL time.Duration) *RedisAdapterTestSuite {
	var suite suite.Suite

	return &RedisAdapterTestSuite{
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

func (suite *RedisAdapterTestSuite) SetupSuite() {
	startLocalRedisServer()
}

func (suite *RedisAdapterTestSuite) TearDownSuite() {
	stopLocalRedisServer()
}

func (suite *RedisAdapterTestSuite) TestNew_NilPool() {
	adapter, err := rediscacheadapters.New(nil, -time.Second)
	suite.Require().Nil(adapter, "Should be nil on nil redis Pool")
	suite.Require().Error(err, "Should give error on nil redis Pool")
}

func (suite *RedisAdapterTestSuite) TestNew_NegativeTTL() {
	adapter, err := rediscacheadapters.New(testRedisPool, -time.Second)
	suite.Require().Nil(adapter, "Should be nil on negative time duration for TTL")
	suite.Require().Error(err, "Should give error on negative time duration for TTL")
}

func (suite *RedisAdapterTestSuite) TestGet_InvalidPool() {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	var actual testutil.TestStruct
	err := adapter.Get(testutil.TestKeyForGet, &actual)

	suite.Require().Equal(testutil.TestStruct{}, actual, "Actual should remain empty since the pool is invalid")
	suite.Require().Error(err, "Should error since the pool is invalid")
}

func (suite *RedisAdapterTestSuite) TestSet_InvalidPool() {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	suite.Require().Error(err, "Should error since the pool is invalid")
}

func (suite *RedisAdapterTestSuite) TestOpenSession_InvalidPool() {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	_, err := adapter.OpenSession()
	suite.Require().Error(err, "Should error on invalid session opening")
}

func (suite *RedisAdapterTestSuite) TestSetTTL_InvalidPool() {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForSetTTL, "1")
	suite.Require().NoError(err, "Must not error on setting test var")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, time.Second)
	suite.Require().Error(err, "Should error since the pool is invalid")
}

func (suite *RedisAdapterTestSuite) TestDelete_InvalidPool() {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForDelete, "1")
	suite.Require().NoError(err, "Must not error on setting test var")

	err = adapter.Delete(testutil.TestKeyForDelete)
	suite.Require().Error(err, "Should error since the pool is invalid")
}

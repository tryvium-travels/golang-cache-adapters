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

package inmemorycacheadapters_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	inmemorycacheadapters "github.com/tryvium-travels/golang-cache-adapters/in_memory"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

// InMemoryAdapterTestSuite contains all methods to run tests in a
// isolated suite.
type InMemoryAdapterTestSuite struct {
	*suite.Suite
	*testutil.CacheAdapterPartialTestSuite
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
		adapter, err := inmemorycacheadapters.New(defaultTTL)
		return adapter.(cacheadapters.CacheSessionAdapter), err
	}
}

// newInMemoryTestSuite creates a new test suite with tests for In-Memory adapters and sessions.
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

func TestInMemoryAdapterSuite(t *testing.T) {
	defaultTTL := 1 * time.Second
	suite.Run(t, newInMemoryTestSuite(t, defaultTTL))
}

func (suite *InMemoryAdapterTestSuite) TestNew_NegativeTTL() {
	adapter, err := inmemorycacheadapters.New(-time.Second)
	suite.Require().Nil(adapter, "Should be nil on negative time duration for TTL")
	suite.Require().Error(err, "Should give error on negative time duration for TTL")
}

func (suite *InMemoryAdapterTestSuite) TestNew_InvalidTTL() {
	adapter, err := inmemorycacheadapters.New(testutil.InvalidTTL)
	suite.Require().Nil(adapter, "Should be nil on invalid time duration for TTL")
	suite.Require().Error(err, "Should give error on invalid time duration for TTL")
}

func (suite *InMemoryAdapterTestSuite) TestNew_InvalidTTLZero() {
	adapter, err := inmemorycacheadapters.New(testutil.ZeroTTL)
	suite.Require().Nil(adapter, "Should be nil on zero time duration for TTL")
	suite.Require().Error(err, "Should give error on zero time duration for TTL")
}

func (suite *InMemoryAdapterTestSuite) TestOpenCloseSessionOK() {
	adapter, errAdapter := inmemorycacheadapters.New(testutil.DummyTTL)
	session, errSession := adapter.OpenSession()

	suite.Require().NoError(errAdapter, "Should give no error on creating an adapter")
	suite.Require().NoError(errSession, "Should give no error on opening a session")

	errClose := session.Close()
	suite.Require().NoError(errClose, "Should not error on closing an open session")
}

func (suite *InMemoryAdapterTestSuite) TestSetGetOK() {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual testutil.TestStruct

	err = adapter.Get(testutil.TestKeyForSet, &actual)
	suite.Require().NoError(err, "Should not error on valid get")

	suite.Require().EqualValues(testutil.TestValue, actual, "The value obatined with get should be equal to the test value set before")
}

func (suite *InMemoryAdapterTestSuite) TestDel_ErrMissing() {
	testKeyForDeleteButInvalid := fmt.Sprintf("%s:but-invalid", testutil.TestKeyForDelete)
	adapter, err := inmemorycacheadapters.New(testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	err = adapter.Delete(testKeyForDeleteButInvalid)
	suite.Require().NoError(err, "Should not error on delete with non-existing key")
}

// test a subsequent Delete operation over the same key without Set between these two Delete.
// It differs from common/TestDelete_OK since it's the fail of the Delete operation
// that is tested
func (suite *InMemoryAdapterTestSuite) TestDel_DoubleDel() {
	adapter, err := inmemorycacheadapters.New(testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	err = adapter.Set(testutil.TestKeyForDelete, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForDelete, &actual)
	suite.Require().Equal(testutil.TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")

	err = adapter.Delete(testutil.TestKeyForDelete)
	suite.Require().NoError(err, "Should not error on valid Delete")

	err = adapter.Delete(testutil.TestKeyForDelete)
	suite.Require().NoError(err, "Should not error on subsequent Delete on the same key solely by this")
}

func (suite *InMemoryAdapterTestSuite) TestTTL_SetDeleteExpires() {
	adapter, _ := suite.NewAdapter()

	duration := 250 * time.Millisecond

	err := adapter.Set(testutil.TestKeyForSetTTL, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, duration)
	suite.Require().NoError(err, "Should not error on valid SetTTL")

	suite.SleepFunc(100 * time.Millisecond)

	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForSetTTL, &actual)
	suite.Require().Equal(testutil.TestValue, actual, "The value just set must be equal to the test value (is not expired)")
	suite.Require().NoError(err, "Should not error on valid get (is not expired)")

	suite.SleepFunc(200 * time.Millisecond)

	err = adapter.Delete(testutil.TestKeyForSetTTL)
	suite.Require().NoError(err, "Should not error on Delete after expires since it makes no differences")
}

func (suite *InMemoryAdapterTestSuite) TestTTL_SetNotExistingKey() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")
	testKeyForSetDLLButInvalid := fmt.Sprintf("%s:but-invalid", testutil.TestKeyForSetTTL)

	duration := new(time.Duration)
	*duration = time.Millisecond * 250
	err = adapter.SetTTL(testKeyForSetDLLButInvalid, *duration)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should error on set with not existing key")
}

func (suite *InMemoryAdapterTestSuite) TestTTL_SetOverExpired() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	duration := new(time.Duration)
	*duration = 50 * time.Millisecond
	err = adapter.Set(testutil.TestKeyForSetTTL, testutil.TestValue, duration)
	suite.Require().NoError(err, "Should not error on valid set")

	suite.SleepFunc(100 * time.Millisecond)

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, (*duration)*2)
	suite.Require().ErrorIs(err, nil, "Should not error on setting TTL over expired key, since it's removed")
}

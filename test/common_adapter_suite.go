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

package testutil

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// CacheAdapterTestSuite contains some methods to run common tests in an
// isolated suite.
//
// Use the SetupTest function to create the correct adapter before each test.
// Test also all New functions and add particular edge cases depending on
// the adapter suite in which this one is put as composition.
type CacheAdapterPartialTestSuite struct {
	*suite.Suite

	// The default TTL for all adapters and sessions.
	DefaultTTL time.Duration

	SleepFunc func(time.Duration)

	// The function to create New instances of the adapter.
	// Example with redis adapter
	//
	//   suite := CacheAdapterPartialTestSuite{
	//	     NewAdapter = func() {
	//           redisPool := redis.Pool{
	//	             // ...
	//           }
	//
	//           DefaultTTL := time.Second
	//
	//           return rediscacheadapters.New(redisPool, DefaultTTL)
	//       }
	//   }
	NewAdapter func() (cacheadapters.CacheAdapter, error)

	// The function to create New instances of the adapter session.
	// Example with redis adapter
	//
	//   suite := CacheAdapterPartialTestSuite{
	//	     NewSession = func() {
	//           redisPool := redis.Pool{
	//	             // ...
	//           }
	//
	//           DefaultTTL := time.Second
	//
	//           adapter, _ := rediscacheadapters.New(redisPool, DefaultTTL)
	//           return adapter.OpenSession()
	//       }
	//   }
	NewSession func() (cacheadapters.CacheSessionAdapter, error)
}

func (suite *CacheAdapterPartialTestSuite) TestNew_OK() {
	adapter, err := suite.NewAdapter()
	suite.Require().NotNil(adapter, "Should not create a nil adapter on valid New")
	suite.Require().NoError(err, "Should not give error on valid New")
}

func (suite *CacheAdapterPartialTestSuite) TestGet_OK() {
	adapter, _ := suite.NewAdapter()

	err := adapter.Set(TestKeyForGet, TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	var actual TestStruct
	err = adapter.Get(TestKeyForGet, &actual)

	suite.Require().Equal(TestValue, actual, "Should be the correct value on a correct get and key not expired")
	suite.Require().NoError(err, "Should not return an error on valid object reference")
}

func (suite *CacheAdapterPartialTestSuite) TestGet_NilReference() {
	adapter, _ := suite.NewAdapter()

	err := adapter.Set(TestKeyForGet, TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	err = adapter.Get(TestKeyForGet, nil)
	suite.Require().ErrorIs(err, cacheadapters.ErrGetRequiresObjectReference, "Should return ErrGetRequiresObjectReference on nil object reference")
}

func (suite *CacheAdapterPartialTestSuite) TestGet_NonUnmarshalableReference() {
	adapter, _ := suite.NewAdapter()

	err := adapter.Set(TestKeyForGet, TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	var actual complex128
	err = adapter.Get(TestKeyForGet, &actual)
	suite.Require().Error(err, "Should not return an error on valid object reference")
}

func (suite *CacheAdapterPartialTestSuite) TestGet_InvalidKey() {
	adapter, _ := suite.NewAdapter()
	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", TestKeyForGet)

	var actual TestStruct
	err := adapter.Get(testKeyForGetButInvalid, &actual)
	suite.Require().Equal(TestStruct{}, actual, "Actual should remain empty since the key is invalid")
	suite.Require().ErrorIs(cacheadapters.ErrNotFound, err, "Should be ErrNotFound since the key is invalid")
}

func (suite *CacheAdapterPartialTestSuite) TestSet_OK() {
	adapter, _ := suite.NewAdapter()

	duration := new(time.Duration)
	*duration = time.Second

	err := adapter.Set(TestKeyForSet, TestValue, duration)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual TestStruct
	err = adapter.Get(TestKeyForSet, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")
}

func (suite *CacheAdapterPartialTestSuite) TestSet_OK_NilTTL() {
	adapter, _ := suite.NewAdapter()

	err := adapter.Set(TestKeyForSet, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set with nil TTL")

	var actual TestStruct
	err = adapter.Get(TestKeyForSet, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")
}

func (suite *CacheAdapterPartialTestSuite) TestSet_NonMarshalableReference() {
	adapter, _ := suite.NewAdapter()

	actualNonMarshalable := complex128(1)
	err := adapter.Set(TestKeyForSet, actualNonMarshalable, nil)
	suite.Require().Error(err, "Should error since the value is not unmarshallable")
}

func (suite *CacheAdapterPartialTestSuite) TestSet_InvalidTTL() {
	adapter, _ := suite.NewAdapter()

	invalidDuration := new(time.Duration)
	*invalidDuration = -time.Second

	err := adapter.Set(TestKeyForSet, TestValue, invalidDuration)
	suite.Require().Error(err, "Should error on valid set with invalid TTL")
}

func (suite *CacheAdapterPartialTestSuite) TestSet_CheckTTL() {
	adapter, _ := suite.NewAdapter()

	duration := new(time.Duration)
	*duration = time.Millisecond * 250

	err := adapter.Set(TestKeyForSet, TestValue, duration)
	suite.Require().NoError(err, "Should not error on valid set")

	suite.SleepFunc(*duration * 2)

	var actual TestStruct
	err = adapter.Get(TestKeyForSet, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired")
}

func (suite *CacheAdapterPartialTestSuite) TestSetTTL_OK() {
	adapter, _ := suite.NewAdapter()

	duration := 250 * time.Millisecond

	err := adapter.Set(TestKeyForSetTTL, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	err = adapter.SetTTL(TestKeyForSetTTL, duration)
	suite.Require().NoError(err, "Should not error on valid SetTTL")

	suite.SleepFunc(100 * time.Millisecond)

	var actual TestStruct
	err = adapter.Get(TestKeyForSetTTL, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value (is not expired)")
	suite.Require().NoError(err, "Should not error on valid get (is not expired)")

	suite.SleepFunc(200 * time.Millisecond)

	err = adapter.Get(TestKeyForSetTTL, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired")
}

func (suite *CacheAdapterPartialTestSuite) TestSetTTL_NegativeTTL() {
	adapter, _ := suite.NewAdapter()

	invalidDuration := -time.Second

	err := adapter.SetTTL(TestKeyForSetTTL, invalidDuration)
	suite.Require().NoError(err, "Should not error on valid set with nil TTL")

	var actual TestStruct
	err = adapter.Get(TestKeyForSetTTL, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired (invalidTTL in SetTTL operation)")
}

func (suite *CacheAdapterPartialTestSuite) TestSetTTL_CheckTTL() {
	adapter, _ := suite.NewAdapter()

	duration := new(time.Duration)
	*duration = time.Millisecond * 250

	err := adapter.Set(TestKeyForSet, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	err = adapter.SetTTL(TestKeyForSet, *duration)
	suite.Require().NoError(err, "Should not error on valid SetTTL")

	var actual TestStruct
	err = adapter.Get(TestKeyForSet, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Should not error (is not expired)")

	suite.SleepFunc(*duration)

	err = adapter.Get(TestKeyForSet, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired")
}

func (suite *CacheAdapterPartialTestSuite) TestDelete_OK() {
	adapter, _ := suite.NewAdapter()

	err := adapter.Set(TestKeyForDelete, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual TestStruct
	err = adapter.Get(TestKeyForDelete, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")

	err = adapter.Delete(TestKeyForDelete)
	suite.Require().NoError(err, "Should not error on valid Delete")

	err = adapter.Get(TestKeyForDelete, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after deleted")
}

func (suite *CacheAdapterPartialTestSuite) TestOpenSession_OK() {
	adapter, _ := suite.NewAdapter()

	session, err := adapter.OpenSession()

	suite.Require().NoError(err, "Should not error on valid session opening")
	defer session.Close()
}

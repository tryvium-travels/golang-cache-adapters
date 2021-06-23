// Copyright 2021 The Tryvium Company LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance _ the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// _OUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"fmt"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

func (suite *CacheAdapterPartialTestSuite) TestSessionNew_OK() {
	session, err := suite.NewSession()
	suite.Require().NotNil(session, "Should not create a nil session on valid New")
	suite.Require().NoError(err, "Should not give error on valid New")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionGet_OK() {
	session, _ := suite.NewSession()
	defer session.Close()

	err := session.Set(TestKeyForGet, TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	var actual TestStruct
	err = session.Get(TestKeyForGet, &actual)

	suite.Require().Equal(TestValue, actual, "Should be the correct value on a correct get and key not expired")
	suite.Require().NoError(err, "Should not return an error on valid object reference")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionGet_NilReference() {
	session, _ := suite.NewSession()
	defer session.Close()

	err := session.Set(TestKeyForGet, TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	err = session.Get(TestKeyForGet, nil)
	suite.Require().ErrorIs(cacheadapters.ErrGetRequiresObjectReference, err, "Should return ErrGetRequiresObjectReference on nil object reference")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionGet_NonUnmarshalableReference() {
	session, _ := suite.NewSession()
	defer session.Close()

	err := session.Set(TestKeyForGet, TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	var actual complex128
	err = session.Get(TestKeyForGet, &actual)
	suite.Require().Error(err, "Should not return an error on valid object reference")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionGet_InvalidKey() {
	session, _ := suite.NewSession()
	defer session.Close()
	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", TestKeyForGet)

	var actual TestStruct
	err := session.Get(testKeyForGetButInvalid, &actual)
	suite.Require().Equal(TestStruct{}, actual, "Actual should remain empty since the key is invalid")
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should be ErrNotFound since the key is invalid")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSet_OK() {
	session, _ := suite.NewSession()
	defer session.Close()

	duration := new(time.Duration)
	*duration = time.Second

	err := session.Set(TestKeyForSet, TestValue, duration)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual TestStruct
	err = session.Get(TestKeyForSet, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSet_OK_NilTTL() {
	session, _ := suite.NewSession()
	defer session.Close()

	err := session.Set(TestKeyForSet, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set with nil TTL")

	var actual TestStruct
	err = session.Get(TestKeyForSet, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSet_NonMarshalableReference() {
	session, _ := suite.NewSession()
	defer session.Close()

	actualNonMarshalable := complex128(1)
	err := session.Set(TestKeyForSet, actualNonMarshalable, nil)
	suite.Require().Error(err, "Should error since the value is not unmarshallable")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSet_InvalidTTL() {
	session, _ := suite.NewSession()
	defer session.Close()

	invalidDuration := new(time.Duration)
	*invalidDuration = -time.Second

	err := session.Set(TestKeyForSet, TestValue, invalidDuration)
	suite.Require().Error(err, "Should error on valid set with invalid TTL")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSet_CheckTTL() {
	session, _ := suite.NewSession()
	defer session.Close()

	duration := new(time.Duration)
	*duration = time.Millisecond * 250

	err := session.Set(TestKeyForSet, TestValue, duration)
	suite.Require().NoError(err, "Should not error on valid set")

	suite.SleepFunc(*duration * 2)

	var actual TestStruct
	err = session.Get(TestKeyForSet, &actual)
	suite.Require().Equal(err, cacheadapters.ErrNotFound, "Should not be found after expired")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSetTTL_OK() {
	session, _ := suite.NewSession()
	defer session.Close()

	duration := 250 * time.Millisecond

	err := session.Set(TestKeyForSetTTL, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	err = session.SetTTL(TestKeyForSetTTL, duration)
	suite.Require().NoError(err, "Should not error on valid SetTTL")

	suite.SleepFunc(100 * time.Millisecond)

	var actual TestStruct
	err = session.Get(TestKeyForSetTTL, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value (is not expired)")
	suite.Require().NoError(err, "Should not error on valid get (is not expired)")

	suite.SleepFunc(200 * time.Millisecond)

	err = session.Get(TestKeyForSet, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSetTTL_NegativeTTL() {
	session, _ := suite.NewSession()
	defer session.Close()

	invalidDuration := -time.Second

	err := session.SetTTL(TestKeyForSetTTL, invalidDuration)
	suite.Require().Error(err, "Should error on invalid set with invalid TTL")

	var actual TestStruct
	err = session.Get(TestKeyForSetTTL, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired (invalidTTL in SetTTL operation)")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionSetTTL_CheckTTL() {
	session, _ := suite.NewSession()
	defer session.Close()

	duration := new(time.Duration)
	*duration = time.Millisecond * 250

	err := session.Set(TestKeyForSet, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	err = session.SetTTL(TestKeyForSet, *duration)
	suite.Require().NoError(err, "Should not error on valid SetTTL")

	var actual TestStruct
	err = session.Get(TestKeyForSet, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Should not error (is not expired)")

	suite.SleepFunc(*duration)

	err = session.Get(TestKeyForSet, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after expired")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionDelete_OK() {
	session, _ := suite.NewSession()
	defer session.Close()

	err := session.Set(TestKeyForDelete, TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual TestStruct
	err = session.Get(TestKeyForDelete, &actual)
	suite.Require().Equal(TestValue, actual, "The value just set must be equal to the test value")
	suite.Require().NoError(err, "Value should be valid, hence no error")

	err = session.Delete(TestKeyForDelete)
	suite.Require().NoError(err, "Should not error on valid Delete")

	err = session.Get(TestKeyForDelete, &actual)
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should not be found after deleted")
}

func (suite *CacheAdapterPartialTestSuite) TestSessionClose_OK() {
	session, _ := suite.NewSession()
	defer session.Close()

	err := session.Close()

	suite.Require().NoError(err, "Should not error on valid session opening")
}

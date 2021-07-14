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

package multicacheadapters_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	multicacheadapters "github.com/tryvium-travels/golang-cache-adapters/multicache"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

func TestMultiCacheAdapterSuite(t *testing.T) {
	suite.Run(t, new(MultiCacheAdapterTestSuite))
}

// MultiCacheAdapterTestSuite contains all methods to run tests in a
// isolated suite.
type MultiCacheAdapterTestSuite struct {
	suite.Suite
	firstDummyAdapter  *mockMultiCacheAdapter
	secondDummyAdapter *mockMultiCacheAdapter
	thirdDummyAdapter  *mockMultiCacheAdapter
}

func (suite *MultiCacheAdapterTestSuite) TestNewOK() {
	adapter, err := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	suite.NotNil(adapter, "Should not be nil if New is ok")
	suite.NoError(err, "Should not give error on valid New")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestNewNoAdapters() {
	adapter, err := multicacheadapters.New()
	suite.Nil(adapter, "Should be nil if New is without adapters")
	suite.ErrorIs(err, multicacheadapters.ErrInvalidSubAdapters, "Should give ErrNilSubadapter on New without adapters")
}

func (suite *MultiCacheAdapterTestSuite) SetupTest() {
	suite.firstDummyAdapter = newmockMultiCacheAdapter()
	suite.secondDummyAdapter = newmockMultiCacheAdapter()
	suite.thirdDummyAdapter = newmockMultiCacheAdapter()
}

func (suite *MultiCacheAdapterTestSuite) TestNewWithNilAdapter() {
	adapter, err := multicacheadapters.New(nil, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	suite.NotNil(adapter, "Should not be nil if one is missing")
	suite.NoError(err, "Should not give error on only nil first adapter")

	adapter, err = multicacheadapters.New(suite.firstDummyAdapter, nil, suite.thirdDummyAdapter)
	suite.NotNil(adapter, "Should not be nil if one is missing")
	suite.NoError(err, "Should not give error on only nil second adapter")

	adapter, err = multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, nil)
	suite.NotNil(adapter, "Should not be nil if one is missing")
	suite.NoError(err, "Should not give error on only nil third adapter")

	adapter, err = multicacheadapters.New(nil, nil, suite.thirdDummyAdapter)
	suite.NotNil(adapter, "Should not be nil if two is missing")
	suite.NoError(err, "Should not give error on only 2 nil adapter")

	adapter, err = multicacheadapters.New(nil, nil, nil)
	suite.Nil(adapter, "Should be nil if all are missing")
	suite.Equal(multicacheadapters.ErrInvalidSubAdapters, err, "Should give error on all nil adapter")
}

func (suite *MultiCacheAdapterTestSuite) TestGetOK() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheAdapterTestSuite) TestGetUsingPriorityOnPartialFail1() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheAdapterTestSuite) TestGetUsingPriorityOnPartialFail2() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheAdapterTestSuite) TestGetUsingPriorityOnPartialFail3() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheAdapterTestSuite) TestGetUsingPriorityOnPartialFailMulti() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheAdapterTestSuite) TestGetOnTotalFailAndDisabledWarnings() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.Error(err, "Should error on total failing Get")
}

func (suite *MultiCacheAdapterTestSuite) TestGetUsingPriorityOnPartialFailAndWarnings() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error on valid Get with a warning")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheAdapterTestSuite) TestGetWithNilReference() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(cacheadapters.ErrGetRequiresObjectReference)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(cacheadapters.ErrGetRequiresObjectReference)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(cacheadapters.ErrGetRequiresObjectReference)

	err := adapter.Get(testutil.TestKeyForGet, nil)
	suite.Error(err, "Should error on Get with an empty reference")
}

func (suite *MultiCacheAdapterTestSuite) TestGetWithNonUnmarshalableReference() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	actual := complex128(1)

	var dummyRawMessage json.RawMessage

	// forcing to return nil simulates a wrong unmarshal handling, corrected by the multi adapter.
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.Error(err, "Should error on Get with a non unmarshalable reference")
}

func (suite *MultiCacheAdapterTestSuite) TestSetOK() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	fakeTTL := time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &fakeTTL)
	suite.NoError(err, "Should not error on OK Set")
}

func (suite *MultiCacheAdapterTestSuite) TestSetOKWithNilTTL() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var nilDuration *time.Duration
	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nilDuration).Once().Return(nil)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nilDuration).Once().Return(nil)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nilDuration).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	suite.NoError(err, "Should not error on OK Set with nil value, but should replace it with defaultTTL")
}

func (suite *MultiCacheAdapterTestSuite) TestSetWithInvalidTTL() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &invalidTTL)
	suite.NoError(err, "Should error on Set with invalid TTL")
}

func (suite *MultiCacheAdapterTestSuite) TestSetWithError() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	var actual complex128

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Set(testutil.TestKeyForSet, actual, &invalidTTL)
	suite.Error(err, "Should error on non total fail value in Set")
}

func (suite *MultiCacheAdapterTestSuite) TestSetWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	var actual complex128

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, actual, &invalidTTL)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in Set")
}

func (suite *MultiCacheAdapterTestSuite) TestSetTTLOK() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	fakeTTL := time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, fakeTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, fakeTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, fakeTTL).Once().Return(nil)

	err := adapter.SetTTL(testutil.TestKeyForSet, fakeTTL)
	suite.NoError(err, "Should not error on OK SetTTL")
}

func (suite *MultiCacheAdapterTestSuite) TestSetTTLWithInvalidTTL() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)

	err := adapter.SetTTL(testutil.TestKeyForSet, invalidTTL)
	suite.NoError(err, "Should error on SetTTL with invalid TTL")
}

func (suite *MultiCacheAdapterTestSuite) TestSetTTLWithError() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.SetTTL(testutil.TestKeyForSet, invalidTTL)
	suite.Error(err, "Should error on non total fail value in SetTTL")
}

func (suite *MultiCacheAdapterTestSuite) TestSetTTLWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)

	err := adapter.SetTTL(testutil.TestKeyForSet, invalidTTL)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in SetTTL")
}

func (suite *MultiCacheAdapterTestSuite) TestDeleteOK() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	suite.firstDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)
	suite.secondDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)
	suite.thirdDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)

	err := adapter.Delete(testutil.TestKeyForDelete)
	suite.NoError(err, "Should not error on OK Delete")
}

func (suite *MultiCacheAdapterTestSuite) TestDeleteWithError() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	suite.firstDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Delete(testutil.TestKeyForDelete)
	suite.Error(err, "Should error on non total fail value in Delete")
}

func (suite *MultiCacheAdapterTestSuite) TestDeleteWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	suite.firstDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)

	err := adapter.Delete(testutil.TestKeyForDelete)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in Delete")
}

func (suite *MultiCacheAdapterTestSuite) TestOpenSessionOK() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	suite.firstDummyAdapter.On("OpenSession").Once().Return(newmockMultiCacheSessionAdapter(), nil)
	suite.secondDummyAdapter.On("OpenSession").Once().Return(newmockMultiCacheSessionAdapter(), nil)
	suite.thirdDummyAdapter.On("OpenSession").Once().Return(newmockMultiCacheSessionAdapter(), nil)

	sessionAdapter, err := adapter.OpenSession()
	suite.NotNil(sessionAdapter, "Session adapter should be initialized")
	suite.NoError(err, "Should not error on OK Open Session")
}

func (suite *MultiCacheAdapterTestSuite) TestOpenSessionWithError() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	suite.firstDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), testutil.ErrTestingFailureCheck)

	_, err := adapter.OpenSession()
	suite.Error(err, "Should error on non total fail value in Open Session")
}

func (suite *MultiCacheAdapterTestSuite) TestOpenSessionWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	suite.firstDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("OpenSession").Once().Return(newmockMultiCacheSessionAdapter(), nil)

	_, err := adapter.OpenSession()
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in Open Session")
}

func (suite *MultiCacheAdapterTestSuite) TestOpenSessionWithPartialErrorNilSubadapter() {
	adapter, _ := multicacheadapters.New(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	suite.firstDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), nil)
	suite.secondDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), nil)
	suite.thirdDummyAdapter.On("OpenSession").Once().Return((*mockMultiCacheSessionAdapter)(nil), nil)

	_, err := adapter.OpenSession()
	suite.ErrorIs(err, multicacheadapters.ErrInvalidSubAdapters, "Should error unitialized subSessionAdapters")
}

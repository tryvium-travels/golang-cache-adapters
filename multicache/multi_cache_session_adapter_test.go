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

package multicacheadapters_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	multicacheadapters "github.com/tryvium-travels/golang-cache-adapters/multicache"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

type mockMultiCacheSessionAdapter struct {
	mock.Mock
	*multicacheadapters.MultiCacheSessionAdapter
}

func (mca *mockMultiCacheSessionAdapter) Get(key string, objectRef interface{}) error {
	args := mca.Called(key, objectRef)

	json.Unmarshal([]byte(testutil.TestValueJSON), &objectRef)

	return args.Error(0)
}

func (mca *mockMultiCacheSessionAdapter) Set(key string, object interface{}, newTTL *time.Duration) error {
	args := mca.Called(key, object, newTTL)

	return args.Error(0)
}

func (mca *mockMultiCacheSessionAdapter) SetTTL(key string, newTTL time.Duration) error {
	args := mca.Called(key, newTTL)

	return args.Error(0)
}

func (mca *mockMultiCacheSessionAdapter) Delete(key string) error {
	args := mca.Called(key)

	return args.Error(0)
}

func newmockMultiCacheSessionAdapter() *mockMultiCacheSessionAdapter {
	return &mockMultiCacheSessionAdapter{}
}

func TestMultiCacheSessionAdapterSuite(t *testing.T) {
	suite.Run(t, new(MultiCacheSessionAdapterTestSuite))
}

// MultiCacheSessionAdapterTestSuite contains all methods to run tests in a
// isolated suite.
type MultiCacheSessionAdapterTestSuite struct {
	suite.Suite
	firstDummyAdapter  *mockMultiCacheSessionAdapter
	secondDummyAdapter *mockMultiCacheSessionAdapter
	thirdDummyAdapter  *mockMultiCacheSessionAdapter
}

func (suite *MultiCacheSessionAdapterTestSuite) TestNewOK() {
	_, err := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	suite.NoError(err, "Should not give error on valid New")
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *MultiCacheSessionAdapterTestSuite) SetupTest() {
	suite.firstDummyAdapter = newmockMultiCacheSessionAdapter()
	suite.secondDummyAdapter = newmockMultiCacheSessionAdapter()
	suite.thirdDummyAdapter = newmockMultiCacheSessionAdapter()
}

func (suite *MultiCacheSessionAdapterTestSuite) TestNewWithNilAdapter() {
	_, err := multicacheadapters.NewSession(nil, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	suite.Equal(multicacheadapters.ErrNilSubAdapter, err, "Should give error on nil first adapter")

	_, err = multicacheadapters.NewSession(suite.firstDummyAdapter, nil, suite.thirdDummyAdapter)
	suite.Equal(multicacheadapters.ErrNilSubAdapter, err, "Should give error on nil second adapter")

	_, err = multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, nil)
	suite.Equal(multicacheadapters.ErrNilSubAdapter, err, "Should give error on nil third adapter")

	_, err = multicacheadapters.NewSession(nil, nil, suite.thirdDummyAdapter)
	suite.Equal(multicacheadapters.ErrNilSubAdapter, err, "Should give error on nil first and second adapter")

	_, err = multicacheadapters.NewSession(nil, nil, nil)
	suite.Equal(multicacheadapters.ErrNilSubAdapter, err, "Should give error on all nil adapter")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetOK() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetUsingPriorityOnPartialFail1() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetUsingPriorityOnPartialFail2() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetUsingPriorityOnPartialFail3() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetUsingPriorityOnPartialFailMulti() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.NoError(err, "Should not error on valid Get")

	suite.Equal(testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetOnTotalFailAndDisabledWarnings() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	var actual testutil.TestStruct

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.Error(err, "Should error on total failing Get")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetUsingPriorityOnPartialFailAndWarnings() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
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

func (suite *MultiCacheSessionAdapterTestSuite) TestGetWithNilReference() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var dummyRawMessage json.RawMessage
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(cacheadapters.ErrGetRequiresObjectReference)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(cacheadapters.ErrGetRequiresObjectReference)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(cacheadapters.ErrGetRequiresObjectReference)

	err := adapter.Get(testutil.TestKeyForGet, nil)
	suite.Error(err, "Should error on Get with an empty reference")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestGetWithNonUnmarshalableReference() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	actual := complex128(1)

	var dummyRawMessage json.RawMessage

	// forcing to return nil simulates a wrong unmarshal handling, corrected by the multi adapter.
	suite.firstDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.secondDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)
	suite.thirdDummyAdapter.On("Get", testutil.TestKeyForGet, &dummyRawMessage).Once().Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	suite.Error(err, "Should error on Get with a non unmarshalable reference")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetOK() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	fakeTTL := time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &fakeTTL)
	suite.NoError(err, "Should not error on OK Set")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetOKWithNilTTL() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	var nilDuration *time.Duration
	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nilDuration).Once().Return(nil)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nilDuration).Once().Return(nil)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nilDuration).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	suite.NoError(err, "Should not error on OK Set with nil value, but should replace it with defaultTTL")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetWithInvalidTTL() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &invalidTTL)
	suite.NoError(err, "Should error on Set with invalid TTL")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetWithError() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	var actual complex128

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Set(testutil.TestKeyForSet, actual, &invalidTTL)
	suite.Error(err, "Should error on non total fail value in Set")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	var actual complex128

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Set", testutil.TestKeyForSet, actual, &invalidTTL).Once().Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, actual, &invalidTTL)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in Set")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetTTLOK() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	fakeTTL := time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, fakeTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, fakeTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, fakeTTL).Once().Return(nil)

	err := adapter.SetTTL(testutil.TestKeyForSet, fakeTTL)
	suite.NoError(err, "Should not error on OK SetTTL")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetTTLWithInvalidTTL() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)

	err := adapter.SetTTL(testutil.TestKeyForSet, invalidTTL)
	suite.NoError(err, "Should error on SetTTL with invalid TTL")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetTTLWithError() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.SetTTL(testutil.TestKeyForSet, invalidTTL)
	suite.Error(err, "Should error on non total fail value in SetTTL")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestSetTTLWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	invalidTTL := -time.Second

	suite.firstDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("SetTTL", testutil.TestKeyForSet, invalidTTL).Once().Return(nil)

	err := adapter.SetTTL(testutil.TestKeyForSet, invalidTTL)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in SetTTL")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestDeleteOK() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	suite.firstDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)
	suite.secondDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)
	suite.thirdDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)

	err := adapter.Delete(testutil.TestKeyForDelete)
	suite.NoError(err, "Should not error on OK Delete")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestDeleteWithError() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	suite.firstDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Delete(testutil.TestKeyForDelete)
	suite.Error(err, "Should error on non total fail value in Delete")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestDeleteWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	suite.firstDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Delete", testutil.TestKeyForDelete).Once().Return(nil)

	err := adapter.Delete(testutil.TestKeyForDelete)
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non marshalable value in Delete")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestCloseOK() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)

	suite.firstDummyAdapter.On("Close").Once().Return(nil)
	suite.secondDummyAdapter.On("Close").Once().Return(nil)
	suite.thirdDummyAdapter.On("Close").Once().Return(nil)

	err := adapter.Close()
	suite.NoError(err, "Should not error on OK Close")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestCloseWithError() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.DisableWarnings()

	suite.firstDummyAdapter.On("Close").Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Close").Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Close").Once().Return(testutil.ErrTestingFailureCheck)

	err := adapter.Close()
	suite.Error(err, "Should error on non total fail value in Close")
}

func (suite *MultiCacheSessionAdapterTestSuite) TestCloseWithPartialErrorAndWarnings() {
	adapter, _ := multicacheadapters.NewSession(suite.firstDummyAdapter, suite.secondDummyAdapter, suite.thirdDummyAdapter)
	adapter.EnableWarnings()

	suite.firstDummyAdapter.On("Close").Once().Return(testutil.ErrTestingFailureCheck)
	suite.secondDummyAdapter.On("Close").Once().Return(testutil.ErrTestingFailureCheck)
	suite.thirdDummyAdapter.On("Close").Once().Return(nil)

	err := adapter.Close()
	suite.Error(err, "Should error on non closable connection")
	suite.ErrorIs(err, multicacheadapters.ErrMultiCacheWarning, "Should error with warning on non closable connection")
}

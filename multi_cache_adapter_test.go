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

package cacheadapters_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

type mockCacheAdapter struct {
	mock.Mock
	cacheadapters.CacheAdapter
}

func newMockCacheAdapter() *mockCacheAdapter {
	return &mockCacheAdapter{
		CacheAdapter: nil,
	}
}

var (
	firstDummyAdapter  *mockCacheAdapter
	secondDummyAdapter *mockCacheAdapter
	thirdDummyAdapter  *mockCacheAdapter
)

// TestMain adds Global test setups and teardowns.
func TestMain(m *testing.M) {
	firstDummyAdapter = newMockCacheAdapter()
	secondDummyAdapter = newMockCacheAdapter()
	thirdDummyAdapter = newMockCacheAdapter()

	code := m.Run()

	os.Exit(code)
}

func TestNewMultiCacheAdapterOK(t *testing.T) {
	_, err := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)
	require.NoError(t, err, "Should not give error on valid New")
}

func TestNewWithNilAdapter(t *testing.T) {
	_, err := cacheadapters.NewMultiCacheAdapter(nil, secondDummyAdapter, thirdDummyAdapter)
	require.Equal(t, cacheadapters.ErrNilSubAdapter, err, "Should give error on nil first adapter")

	_, err = cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, nil, thirdDummyAdapter)
	require.Equal(t, cacheadapters.ErrNilSubAdapter, err, "Should give error on nil second adapter")

	_, err = cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, nil)
	require.Equal(t, cacheadapters.ErrNilSubAdapter, err, "Should give error on nil third adapter")

	_, err = cacheadapters.NewMultiCacheAdapter(nil, nil, thirdDummyAdapter)
	require.Equal(t, cacheadapters.ErrNilSubAdapter, err, "Should give error on nil first and second adapter")

	_, err = cacheadapters.NewMultiCacheAdapter(nil, nil, nil)
	require.Equal(t, cacheadapters.ErrNilSubAdapter, err, "Should give error on all nil adapter")
}

func TestGetOK(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual testutil.TestStruct

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.NoError(t, err, "Should not error on valid Get")

	require.Equal(t, testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func TestGetUsingPriorityOnPartialFail1(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual testutil.TestStruct

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.NoError(t, err, "Should not error on valid Get")

	require.Equal(t, testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func TestGetUsingPriorityOnPartialFail2(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual testutil.TestStruct

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.NoError(t, err, "Should not error on valid Get")

	require.Equal(t, testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func TestGetUsingPriorityOnPartialFail3(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual testutil.TestStruct

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.NoError(t, err, "Should not error on valid Get")

	require.Equal(t, testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func TestGetUsingPriorityOnPartialFailMulti(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual testutil.TestStruct

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(nil)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.NoError(t, err, "Should not error on valid Get")

	require.Equal(t, testutil.TestValue.Value, actual.Value, "Should be equal to the provided test value")
}

func TestGetOnTotalFail(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual testutil.TestStruct

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, testutil.TestStruct{}).Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.Error(t, err, "Should error on total failing Get")
}

func TestGetWithNilReference(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, nil).Return(cacheadapters.ErrGetRequiresObjectReference)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, nil).Return(cacheadapters.ErrGetRequiresObjectReference)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, nil).Return(cacheadapters.ErrGetRequiresObjectReference)

	err := adapter.Get(testutil.TestKeyForGet, nil)
	require.Error(t, err, "Should error on Get with an empty reference (and not in a transaction)")
}

func TestGetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	actual := complex128(1)

	firstDummyAdapter.On("Get", testutil.TestKeyForGet, nil).Return(testutil.ErrTestingFailureCheck)
	secondDummyAdapter.On("Get", testutil.TestKeyForGet, nil).Return(testutil.ErrTestingFailureCheck)
	thirdDummyAdapter.On("Get", testutil.TestKeyForGet, nil).Return(testutil.ErrTestingFailureCheck)

	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.Error(t, err, "Should error on Get with a non unmarshalable reference")
}

func TestSetOK(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	fakeTTL := time.Second

	firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Return(nil)
	secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Return(nil)
	thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &fakeTTL).Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &fakeTTL)
	require.NoError(t, err, "Should not error on OK Set")
}

func TestSetOKWithNilTTL(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nil).Return(nil)
	secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nil).Return(nil)
	thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, nil).Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	require.Error(t, err, "Should not error on OK Set with nil value, but should replace it with defaultTTL")
}

func TestSetWithInvalidTTL(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	invalidTTL := -time.Second

	firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Return(nil)
	secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Return(nil)
	thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &invalidTTL)
	require.NoError(t, err, "Should error on Set with invalid TTL")
}

func TestSetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := cacheadapters.NewMultiCacheAdapter(firstDummyAdapter, secondDummyAdapter, thirdDummyAdapter)

	var actual complex128

	invalidTTL := -time.Second

	firstDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Return(nil)
	secondDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Return(nil)
	thirdDummyAdapter.On("Set", testutil.TestKeyForSet, testutil.TestValue, &invalidTTL).Return(nil)

	err := adapter.Set(testutil.TestKeyForSet, &actual, &invalidTTL)
	require.NoError(t, err, "Should error on non marshalable value in Set")
}

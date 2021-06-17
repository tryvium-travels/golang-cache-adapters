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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	inmemorycacheadapters "github.com/tryvium-travels/golang-cache-adapters/in_memory"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

func TestNewInMemoryCacheAdapterOK(t *testing.T) {
	_, err := inmemorycacheadapters.New(testutil.DummyTTL)
	require.NoError(t, err, "Should not give error on valid New")
}

func TestNewInMemoryCacheAdapterWithInvalidTTL(t *testing.T) {
	_, err := inmemorycacheadapters.New(testutil.InvalidTTL)
	require.Error(t, err, "Should give error on valid New with invalid TTL")
}

func TestNewInMemoryCacheAdapterWithInvalidTTLZero(t *testing.T) {
	_, err := inmemorycacheadapters.New(testutil.ZeroTTL)
	require.Error(t, err, "Should give error on valid New with zero TTL")
}

func TestOpenCloseSessionOK(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)
	session, _ := adapter.OpenSession()

	err := session.Close()
	require.NoError(t, err, "Should not error on closing an open session")
}

func TestSetGetOK(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &testutil.DummyTTL)
	require.NoError(t, err, "Should not error on valid set")

	var actual testutil.TestStruct

	err = adapter.Get(testutil.TestKeyForSet, &actual)
	require.NoError(t, err, "Should not error on valid get")

	require.EqualValues(t, testutil.TestValue, actual, "The value obatined with get should be equal to the test value set before")
}

func TestSetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)

	var actual complex128

	err := adapter.Get(testutil.TestKeyForSet, &actual)
	require.Error(t, err, "Should error on non unmarshalable get")
}

// ----------------------------------------

func TestInTransaction_GetFailBeforeSet(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)
	data := []interface{}{
		testutil.TestStruct{Value: "test1"},
		testutil.TestStruct{Value: "test2"},
	}
	err := adapter.InTransaction(testutil.InTransactionFunc, data)
	require.Error(t, err, "Should error on inner Get since the cache is empty")
}

func TestInTransaction_GetSetOk(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)
	data := []interface{}{
		testutil.TestStruct{Value: "test1"},
		testutil.TestStruct{Value: "test2"},
	}
	objToGet := testutil.TestStruct{Value: "3-get"}
	err := adapter.Set(testutil.TestKeyForGet, objToGet, nil)
	require.NoError(t, err, "Should not give error on setting the value")

	err = adapter.InTransaction(
		func(session cacheadapters.CacheSessionAdapter) error {
			getHolder := &testutil.TestStruct{Value: "hold me"}
			err2 := session.Get(testutil.TestKeyForGet, getHolder)
			if err2 != nil {
				return err2
			}

			testValue2 := testutil.TestStruct{
				Value: "222",
			}

			err2 = session.Set(testutil.TestKeyForSet, testValue2, nil)
			if err2 != nil {
				return err2
			}

			return nil
		}, data)

	require.NoError(t, err, "Should not give errors on InTransaction since Get inside InTransaction should succeed")

	if data[0] == nil {
		errorText := "Processed data at index 0 is nil"
		require.NoError(t, fmt.Errorf(errorText), errorText)
	}
	if data[1] == nil {
		errorText := "Processed data at index 1 is nil"
		require.NoError(t, fmt.Errorf(errorText), errorText)
	}

	fmt.Printf("data[0]: %+v\n", data[0])
	fmt.Printf("data[0].([]interface{})[0]: %+v\n", data[0].([]interface{})[0])
	fmt.Printf("all data: %+v\n", data)
	require.Equal(t, "3-get",
		((data[0].([]interface{})[0]).(*testutil.TestStruct)).Value,
		"Processed data at index 1 should have the value \"3-get\"",
	)
	require.Equal(t, "222",
		((data[1].([]interface{})[0]).(*testutil.TestStruct)).Value,
		"Processed data at index 1 should have the value \"222\"",
	)
}

func TestInTransaction_GetSetOkTrimming(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)
	data := []interface{}{
		testutil.TestStruct{Value: "test1"},
		testutil.TestStruct{Value: "test2"},
		testutil.TestStruct{Value: "will be trimmed"},
	}
	err := adapter.Set(testutil.TestKeyForGet, testutil.TestStruct{Value: "3"}, nil)
	require.NoError(t, err, "Should not give error on setting the value")

	err = adapter.InTransaction(
		func(session cacheadapters.CacheSessionAdapter) error {
			getHolder := &testutil.TestStruct{Value: "hold me"}
			err2 := session.Get(testutil.TestKeyForGet, getHolder)
			if err2 != nil {
				return err2
			}

			testValue2 := testutil.TestStruct{
				Value: "222",
			}

			err2 = session.Set(testutil.TestKeyForSet, testValue2, nil)
			if err2 != nil {
				return err2
			}

			return nil
		}, data)

	require.NoError(t, err, "Should not give errors on InTransaction since Get inside InTransaction should succeed")

	if data[0] == nil {
		errorText := "Processed data at index 0 is nil"
		require.NoError(t, fmt.Errorf(errorText), errorText)
	}
	if data[1] == nil {
		errorText := "Processed data at index 1 is nil"
		require.NoError(t, fmt.Errorf(errorText), errorText)
	}
	if data[2] != nil {
		errorText := "Exceeding processed data should be nil"
		require.NoError(t, fmt.Errorf(errorText), errorText)
	}

	require.Equal(t, "3",
		((data[0].([]interface{})[0]).(*testutil.TestStruct)).Value,
		"Processed data at index 1 should have the value \"3\"",
	)
	require.Equal(t, "222",
		((data[1].([]interface{})[0]).(*testutil.TestStruct)).Value,
		"Processed data at index 1 should have the value \"222\"",
	)
}

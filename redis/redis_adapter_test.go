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

package rediscacheadapters_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

func TestNewOK(t *testing.T) {
	_, err := rediscacheadapters.New(testRedisPool, time.Second)
	require.NoError(t, err, "Should not give error on valid New")
}

func TestNewWithNilPool(t *testing.T) {
	_, err := rediscacheadapters.New(nil, -time.Second)
	require.Error(t, err, "Should give error on nil redis Pool")
}

func TestNewWithNegativeDuration(t *testing.T) {
	_, err := rediscacheadapters.New(testRedisPool, -time.Second)
	require.Error(t, err, "Should give error on negative time Duration for TTL")
}

func TestGetOK(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	var actual testutil.TestStruct
	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.Equal(t, testutil.TestValue, actual, "Should be the correct value on a correct get and key not expired")
	require.NoError(t, err, "Should not return an error on valid object reference")
}

func TestGetWithNilReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := adapter.Get(testutil.TestKeyForGet, nil)
	require.Equal(t, cacheadapters.ErrGetRequiresObjectReference, err, "Should return ErrGetRequiresObjectReference on nil object reference")
}

func TestGetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	actual := complex128(1)
	err := adapter.Get(testutil.TestKeyForGet, &actual)
	require.Error(t, err, "Should return an error on non unmarshalable object reference")
}

func TestGetWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	var actual testutil.TestStruct
	err := adapter.Get(testutil.TestKeyForGet, &actual)

	require.Equal(t, testutil.TestStruct{}, actual, "Actual should remain empty since the pool is invalid")
	require.Error(t, err, "Should error since the pool is invalid")
}

func TestGetWithInvalidKey(t *testing.T) {
	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", testutil.TestKeyForGet)

	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	var actual testutil.TestStruct
	err := adapter.Get(testKeyForGetButInvalid, &actual)

	require.Equal(t, testutil.TestStruct{}, actual, "Actual should remain empty since the key is invalid")
	require.Equal(t, cacheadapters.ErrNotFound, err, "Should be ErrNotFound since the key is invalid")
}

func TestOpenSessionOK(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	session, err := adapter.OpenSession()
	require.NoError(t, err, "Should not error on valid session opening")
	defer session.Close()
}

func TestSetOK(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	duration := new(time.Duration)
	*duration = time.Second

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, duration)
	require.NoError(t, err, "Should not error on valid set")

	testValueContent, err := localRedisServer.Get(testutil.TestKeyForSet)
	require.NoError(t, err, "Value just set must exist, hence no error")

	var actual testutil.TestStruct
	err = json.Unmarshal([]byte(testValueContent), &actual)
	require.NoError(t, err, "Value just set be a valid JSON, hence no error")

	require.Equal(t, testutil.TestValue, actual, "The value just set must be equal to the test value")
}

func TestSetOKWithNilTTL(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	require.NoError(t, err, "Should not error on valid set")

	testValueContent, err := localRedisServer.Get(testutil.TestKeyForSet)
	require.NoError(t, err, "Value just set must exist, hence no error")

	var actual testutil.TestStruct
	err = json.Unmarshal([]byte(testValueContent), &actual)
	require.NoError(t, err, "Value just set be a valid JSON, hence no error")

	require.Equal(t, testutil.TestValue, actual, "The value just set must be equal to the test value")
}

func TestSetWithInvalidTTL(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	duration := new(time.Duration)
	*duration = -time.Second

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, duration)
	require.Error(t, err, "Should error on invalid duration")
}

func TestSetWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	require.Error(t, err, "Should error since the pool is invalid")
}

func TestSetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	actualNonUnmarshallable := complex128(1)
	err := adapter.Set(testutil.TestKeyForSet, actualNonUnmarshallable, nil)
	require.Error(t, err, "Should error since the value is not unmarshallable")
}

func TestOpenSessionWithInvalidRedisPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	_, err := adapter.OpenSession()
	require.Error(t, err, "Should error on invalid session opening")
}

func TestSetTTLOK(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForSetTTL, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, time.Second*5)
	require.NoError(t, err, "Must not error on setting the expiration")

	// goes into the future when the key is expired
	localRedisServer.FastForward(time.Second * 6)

	_, err = localRedisServer.Get(testutil.TestKeyForSetTTL)
	require.Equal(t, miniredis.ErrKeyNotFound, err, "Must not find the expired key")
}

func TestSetTTLExpired(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForSetTTL, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, cacheadapters.TTLExpired)
	require.NoError(t, err, "Must not error on setting the expiration")

	_, err = localRedisServer.Get(testutil.TestKeyForSetTTL)
	require.Equal(t, miniredis.ErrKeyNotFound, err, "Must not find the expired key")
}

func TestSetTTLWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForSetTTL, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, time.Second)
	require.Error(t, err, "Should error since the pool is invalid")
}

func TestDeleteOK(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForDelete, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = adapter.Delete(testutil.TestKeyForDelete)
	require.NoError(t, err, "Must not error on deleting the key")

	_, err = localRedisServer.Get(testutil.TestKeyForDelete)
	require.Equal(t, miniredis.ErrKeyNotFound, err, "Must not find the deleted key")
}

func TestDeleteWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForDelete, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = adapter.Delete(testutil.TestKeyForDelete)
	require.Error(t, err, "Should error since the pool is invalid")
}

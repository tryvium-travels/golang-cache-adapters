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
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
)

func initConnection(t *testing.T) redis.Conn {
	conn, err := testRedisPool.Dial()
	if err != nil {
		t.Skip("Skipped because connection has not been created properly")
	}

	return conn
}

func TestNewSessionOK(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()
	_, err := rediscacheadapters.NewSession(conn, time.Second)
	require.NoError(t, err, "Should not give error on valid NewSession")
}

func TestNewSessionWithNilConn(t *testing.T) {
	_, err := rediscacheadapters.NewSession(nil, time.Second)
	require.Error(t, err, "Should give error on nil redis connection")
}

func TestNewSessionWithNegativeDuration(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()
	_, err := rediscacheadapters.NewSession(conn, -time.Second)
	require.Error(t, err, "Should give error on negative time Duration for TTL when creating a session")
}

func TestSessionGetOK(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	var actual testStruct
	err := session.Get(testKeyForGet, &actual)
	require.Equal(t, testValue, actual, "Should be the correct value on a correct get and key not expired")
	require.NoError(t, err, "Should not return an error on valid object reference")
}

func TestSessionGetWithNilReference(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := session.Get(testKeyForGet, nil)
	require.Equal(t, cacheadapters.ErrGetRequiresObjectReference, err, "Should return ErrGetRequiresObjectReference on nil object reference")
}

func TestSessionGetWithNonUnmarshalableReference(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	actual := complex128(1)
	err := session.Get(testKeyForGet, &actual)
	require.Error(t, err, "Should return an error on non unmarshalable object reference")
}

func TestSessionGetWithInvalidConnection(t *testing.T) {
	conn := initConnection(t)

	// by closing the connection we make it invalid
	conn.Close()

	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", testKeyForGet)

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	var actual testStruct
	err := session.Get(testKeyForGetButInvalid, &actual)

	require.Equal(t, testStruct{}, actual, "Actual should remain empty since the connection is invalid (already closed)")
	require.Error(t, err, "Should error since the connection is invalid (already closed)")
}

func TestSessionGetWithInvalidKey(t *testing.T) {
	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", testKeyForGet)

	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	var actual testStruct
	err := session.Get(testKeyForGetButInvalid, &actual)

	require.Equal(t, testStruct{}, actual, "Actual should remain empty since the key is invalid")
	require.Equal(t, cacheadapters.ErrNotFound, err, "Should be ErrNotFound since the key is invalid")
}

func TestSessionSetOK(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	duration := new(time.Duration)
	*duration = time.Second

	err := session.Set(testKeyForSet, testValue, duration)
	require.NoError(t, err, "Should not error on valid set")

	testValueContent, err := localRedisServer.Get(testKeyForSet)
	require.NoError(t, err, "Value just set must exist, hence no error")

	var actual testStruct
	err = json.Unmarshal([]byte(testValueContent), &actual)
	require.NoError(t, err, "Value just set be a valid JSON, hence no error")

	require.Equal(t, testValue, actual, "The value just set must be equal to the test value")
}

func TestSessionSetOKWithNilTTL(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := session.Set(testKeyForSet, testValue, nil)
	require.NoError(t, err, "Should not error on valid set")

	testValueContent, err := localRedisServer.Get(testKeyForSet)
	require.NoError(t, err, "Value just set must exist, hence no error")

	var actual testStruct
	err = json.Unmarshal([]byte(testValueContent), &actual)
	require.NoError(t, err, "Value just set be a valid JSON, hence no error")

	require.Equal(t, testValue, actual, "The value just set must be equal to the test value")
}

func TestSessionSetWithInvalidConnection(t *testing.T) {
	conn := initConnection(t)

	// by closing the connection we make it invalid
	conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := session.Set(testKeyForSet, testValue, nil)
	require.Error(t, err, "Should error since the connection is invalid (already closed)")
}

func TestSessionSetWithNonUnmarshalableReference(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	actualNonUnmarshallable := complex128(1)
	err := session.Set(testKeyForSet, actualNonUnmarshallable, nil)
	require.Error(t, err, "Should error since the value is not unmarshallable")
}
func TestSessionSetWithNegativeTTL(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	duration := new(time.Duration)
	*duration = -time.Second

	err := session.Set(testKeyForSet, testValue, duration)
	require.Error(t, err, "Should give error on setting a value with negative time Duration for TTL")
}

func TestSessionSetTTLOK(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := localRedisServer.Set(testKeyForSetTTL, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = session.SetTTL(testKeyForSetTTL, time.Second*5)
	require.NoError(t, err, "Must not error on setting the expiration")

	// goes into the future when the key is expired
	localRedisServer.FastForward(time.Second * 6)

	_, err = localRedisServer.Get(testKeyForSetTTL)
	require.Equal(t, miniredis.ErrKeyNotFound, err, "Must not find the expired key")
}

func TestSessionSetTTLExpired(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := localRedisServer.Set(testKeyForSetTTL, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = session.SetTTL(testKeyForSetTTL, cacheadapters.TTLExpired)
	require.NoError(t, err, "Must not error on setting the expiration")

	_, err = localRedisServer.Get(testKeyForSetTTL)
	require.Equal(t, miniredis.ErrKeyNotFound, err, "Must not find the expired key")
}

func TestSessionSetTTLWithInvalidConnection(t *testing.T) {
	conn := initConnection(t)

	// by closing the connection we make it invalid
	conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := localRedisServer.Set(testKeyForSetTTL, "1")
	require.NoError(t, err, "Must not error on setting test var")

	err = session.SetTTL(testKeyForSetTTL, time.Second)
	require.Error(t, err, "Should error since the conn is invalid")
}

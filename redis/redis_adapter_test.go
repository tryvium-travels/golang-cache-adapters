package rediscacheadapters_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
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

	var actual testStruct
	err := adapter.Get(testKeyForGet, &actual)
	require.Equal(t, testValue, actual, "Should be the correct value on a correct get and key not expired")
	require.NoError(t, err, "Should not return an error on valid object reference")
}

func TestGetWithNilReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := adapter.Get(testKeyForGet, nil)
	require.Equal(t, cacheadapters.ErrGetRequiresObjectReference, err, "Should return ErrGetRequiresObjectReference on nil object reference")
}

func TestGetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	actual := complex128(1)
	err := adapter.Get(testKeyForGet, &actual)
	require.Error(t, err, "Should return an error on non unmarshalable object reference")
}

func TestGetWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	var actual testStruct
	err := adapter.Get(testKeyForGet, &actual)

	require.Equal(t, testStruct{}, actual, "Actual should remain empty since the pool is invalid")
	require.Error(t, err, "Should error since the pool is invalid")
}

func TestGetWithInvalidKey(t *testing.T) {
	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", testKeyForGet)

	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	var actual testStruct
	err := adapter.Get(testKeyForGetButInvalid, &actual)

	require.Equal(t, testStruct{}, actual, "Actual should remain empty since the key is invalid")
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

	err := adapter.Set(testKeyForSet, testValue, duration)
	require.NoError(t, err, "Should not error on valid set")

	testValueContent, err := localRedisServer.Get(testKeyForSet)
	require.NoError(t, err, "Value just set must exist, hence no error")

	var actual testStruct
	err = json.Unmarshal([]byte(testValueContent), &actual)
	require.NoError(t, err, "Value just set be a valid JSON, hence no error")

	require.Equal(t, testValue, actual, "The value just set must be equal to the test value")
}

func TestSetOKWithNilTTL(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := adapter.Set(testKeyForSet, testValue, nil)
	require.NoError(t, err, "Should not error on valid set")

	testValueContent, err := localRedisServer.Get(testKeyForSet)
	require.NoError(t, err, "Value just set must exist, hence no error")

	var actual testStruct
	err = json.Unmarshal([]byte(testValueContent), &actual)
	require.NoError(t, err, "Value just set be a valid JSON, hence no error")

	require.Equal(t, testValue, actual, "The value just set must be equal to the test value")
}

func TestSetWithInvalidTTL(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	duration := new(time.Duration)
	*duration = -time.Second

	err := adapter.Set(testKeyForSet, testValue, duration)
	require.Error(t, err, "Should error on invalid duration")
}

func TestSetWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	err := adapter.Set(testKeyForSet, testValue, nil)
	require.Error(t, err, "Should error since the pool is invalid")
}

func TestSetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	actualNonUnmarshallable := complex128(1)
	err := adapter.Set(testKeyForSet, actualNonUnmarshallable, nil)
	require.Error(t, err, "Should error since the value is not unmarshallable")
}

func TestOpenSessionWithInvalidRedisPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	_, err := adapter.OpenSession()
	require.Error(t, err, "Should error on invalid session opening")
}

func TestInTransactionOK(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	refs := []interface{}{
		&testStruct{},
		nil,
	}

	err := adapter.InTransaction(inTransactionFunc, refs)
	require.NoError(t, err, "Should not error with a valid transaction")

	expectedValues := []interface{}{
		&testValue,
		nil,
	}

	require.Equal(t, expectedValues, refs, "Should be equal to the result")
}

func TestInTransactionWithNilObjectReference(t *testing.T) {
	adapter, _ := rediscacheadapters.New(testRedisPool, time.Second)

	err := adapter.InTransaction(inTransactionFunc, nil)
	require.Equal(t, err, cacheadapters.ErrGetRequiresObjectReference, "Should error with nil object references in a transaction")
}

func TestInTransactionWithInvalidPool(t *testing.T) {
	adapter, _ := rediscacheadapters.New(invalidRedisPool, time.Second)

	refs := []interface{}{
		&testStruct{},
		nil,
	}

	err := adapter.InTransaction(inTransactionFunc, refs)
	require.Error(t, err, "Should error since the pool is invalid")
}
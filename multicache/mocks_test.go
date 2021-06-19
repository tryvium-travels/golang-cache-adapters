package multicacheadapters_test

import (
	"encoding/json"
	"time"

	"github.com/stretchr/testify/mock"
	multicacheadapters "github.com/tryvium-travels/golang-cache-adapters/multicache"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

type mockMultiCacheAdapter struct {
	mock.Mock
	*multicacheadapters.MultiCacheAdapter
}

func (mca *mockMultiCacheAdapter) Get(key string, objectRef interface{}) error {
	args := mca.Called(key, objectRef)

	json.Unmarshal([]byte(testutil.TestValueJSON), &objectRef)

	return args.Error(0)
}

func (mca *mockMultiCacheAdapter) Set(key string, object interface{}, newTTL *time.Duration) error {
	args := mca.Called(key, object, newTTL)

	return args.Error(0)
}

func (mca *mockMultiCacheAdapter) SetTTL(key string, newTTL time.Duration) error {
	args := mca.Called(key, newTTL)

	return args.Error(0)
}

func (mca *mockMultiCacheAdapter) Delete(key string) error {
	args := mca.Called(key)

	return args.Error(0)
}

func newmockMultiCacheAdapter() *mockMultiCacheAdapter {
	return &mockMultiCacheAdapter{}
}

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

func (mca *mockMultiCacheSessionAdapter) Close() error {
	args := mca.Called()

	return args.Error(0)
}

func newmockMultiCacheSessionAdapter() *mockMultiCacheSessionAdapter {
	return &mockMultiCacheSessionAdapter{}
}

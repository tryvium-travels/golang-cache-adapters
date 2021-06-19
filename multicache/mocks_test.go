package multicacheadapters_test

import (
	"encoding/json"
	"time"

	"github.com/stretchr/testify/mock"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

type mockMultiCacheAdapter struct {
	mock.Mock
	initialized bool
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

func (mca *mockMultiCacheAdapter) OpenSession() (cacheadapters.CacheSessionAdapter, error) {
	args := mca.Called()

	return args.Get(0).(*mockMultiCacheSessionAdapter), args.Error(1)
}

func newmockMultiCacheAdapter() *mockMultiCacheAdapter {
	return &mockMultiCacheAdapter{initialized: true}
}

type mockMultiCacheSessionAdapter struct {
	mock.Mock
	initialized bool
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
	return &mockMultiCacheSessionAdapter{initialized: true}
}

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
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/mock"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

var (
	localRedisServer *miniredis.Miniredis // The local in-memory redis instance
	testRedisPool    *redis.Pool          // The pool used in all the tests, except for the "InvalidPool" ones.
	invalidRedisPool *redis.Pool          // The pool used when in need to test invalid connection behaviours.
)

// erroringMockedRedisConn mocks the redis call to increase code coverage
type erroringMockedRedisConn struct {
	mock.Mock
	redis.Conn
}

type erroringMockedSENDMULTIRedisConn struct {
	*erroringMockedRedisConn
}

func (emc *erroringMockedSENDMULTIRedisConn) Send(commandName string, args ...interface{}) error {
	mockArgs := emc.Called(append([]interface{}{commandName}, args...)...)
	return mockArgs.Error(0)
}

type erroringMockedDOEXECRedisConn struct {
	*erroringMockedSENDMULTIRedisConn
}

func (emc *erroringMockedDOEXECRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	mockArgs := emc.Called(append([]interface{}{commandName}, args...)...)
	return mockArgs.Get(0), mockArgs.Error(1)
}

func setGetExFloat64InTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	err := session.Set(testutil.TestKeyForSet, 2.5, nil)
	if err != nil {
		return err
	}

	err = session.Get(testutil.TestKeyForSet, nil)
	if err != nil {
		return err
	}

	err = session.SetTTL(testutil.TestKeyForSet, time.Second)
	if err != nil {
		return err
	}

	return nil
}

// startLocalRedisServer starts a local, in-memory redis instance for the tests.
func startLocalRedisServer() {
	var err error

	invalidRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return nil, errors.New("TESTING INVALID DIAL FROM POOL")
		},
	}

	localRedisServer, err = miniredis.Run()
	if err != nil {
		log.Fatalf("Cannot start local redis server: %s", err)
	}

	localRedisServer.Select(0)

	testValueContent, err := json.Marshal(testutil.TestValue)
	if err != nil {
		log.Fatalf("Cannot set initial testKeyForGet on local redis: %s", err)
	}

	// set initial value for testKeyForGet
	err = localRedisServer.Set(testutil.TestKeyForGet, string(testValueContent))
	if err != nil {
		log.Fatalf("Cannot set initial testKeyForGet on local redis: %s", err)
	}

	testRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", localRedisServer.Addr())
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}
}

// stopLocalRedisServer stops the previously started,
// local, in-memory Redis instance.
func stopLocalRedisServer() {
	localRedisServer.Close()
}

// TestMain adds Global test setups and teardowns.
func TestMain(m *testing.M) {
	startLocalRedisServer()
	defer stopLocalRedisServer()

	code := m.Run()

	os.Exit(code)
}

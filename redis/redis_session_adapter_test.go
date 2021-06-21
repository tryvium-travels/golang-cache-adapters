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
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	rediscacheadapters "github.com/tryvium-travels/golang-cache-adapters/redis"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

func newTestSessionFunc(t *testing.T, defaultTTL time.Duration) func() (cacheadapters.CacheSessionAdapter, error) {
	return func() (cacheadapters.CacheSessionAdapter, error) {
		conn, err := testRedisPool.Dial()
		if err != nil {
			t.Error(err)
		}

		return rediscacheadapters.NewSession(conn, defaultTTL)
	}
}

func (suite *RedisAdapterTestSuite) initCustomConnection() redis.Conn {
	conn, err := testRedisPool.Dial()
	if err != nil {
		suite.T().Fatal("Skipped because connection has not been created properly")
	}

	return conn
}

func (suite *RedisAdapterTestSuite) TestNewSession_NilConnection() {
	session, err := rediscacheadapters.NewSession(nil, time.Second)
	suite.Require().Nil(session, "Should be nil if I pass a nil redis connection")
	suite.Require().Equal(cacheadapters.ErrInvalidConnection, err, "Should give error on nil redis connection")
}

func (suite *RedisAdapterTestSuite) TestNewSession_NegativeDuration() {
	conn := suite.initCustomConnection()
	defer conn.Close()

	session, err := rediscacheadapters.NewSession(conn, -time.Second)
	suite.Require().Nil(session, "Should be nil on negative time Duration for TTL when creating a session")
	suite.Require().Error(err, "Should give error on negative time Duration for TTL when creating a session")
}

func (suite *RedisAdapterTestSuite) TestSessionGet_InvalidConnection() {
	conn := suite.initCustomConnection()

	// by closing the connection we make it invalid
	conn.Close()

	testKeyForGetButInvalid := fmt.Sprintf("%s:but-invalid", testutil.TestKeyForGet)

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	var actual testutil.TestStruct
	err := session.Get(testKeyForGetButInvalid, &actual)

	suite.Require().Equal(testutil.TestStruct{}, actual, "Actual should remain empty since the connection is invalid (already closed)")
	suite.Require().Error(err, "Should error since the connection is invalid (already closed)")
}

func (suite *RedisAdapterTestSuite) TestSessionSet_InvalidConnection() {
	conn := suite.initCustomConnection()

	// by closing the connection we make it invalid
	conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := session.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	suite.Require().Error(err, "Should error since the connection is invalid (already closed)")
}

func (suite *RedisAdapterTestSuite) TestSessionSetTTL_InvalidConnection() {
	conn := suite.initCustomConnection()

	// by closing the connection we make it invalid
	conn.Close()

	session, _ := rediscacheadapters.NewSession(conn, time.Second)

	err := localRedisServer.Set(testutil.TestKeyForSetTTL, "1")
	suite.Require().NoError(err, "Must not error on setting test var")

	err = session.SetTTL(testutil.TestKeyForSetTTL, time.Second)
	suite.Require().Error(err, "Should error since the conn is invalid")
}

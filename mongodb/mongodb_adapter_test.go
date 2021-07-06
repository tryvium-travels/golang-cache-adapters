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

package mongodbcacheadapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	mongodbcacheadapters "github.com/tryvium-travels/golang-cache-adapters/mongodb"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDBAdapterSuite(t *testing.T) {
	defaultTTL := 1 * time.Second
	suite.Run(t, newMongoDBTestSuite(t, defaultTTL))
}

type MongoDBTestSuite struct {
	*suite.Suite
	*testutil.CacheAdapterPartialTestSuite
}

func (suite *MongoDBTestSuite) SetupSuite() {
	startLocalMongoDBServer()
}

func (suite *MongoDBTestSuite) TearDownSuite() {
	stopLocalMongoDBServer()
}

func newTestAdapterFunc(defaultTTL time.Duration) func() (cacheadapters.CacheAdapter, error) {
	return func() (cacheadapters.CacheAdapter, error) {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
		if err != nil {
			panic(err)
		}

		return mongodbcacheadapters.New(client, testDatabase, testCollection, defaultTTL)
	}
}

func testSleepFunc() func(time.Duration) {
	return func(duration time.Duration) {
		time.Sleep(duration)
	}
}

func newTestSessionFunc(t *testing.T, defaultTTL time.Duration) func() (cacheadapters.CacheSessionAdapter, error) {
	return func() (cacheadapters.CacheSessionAdapter, error) {
		mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
		if err != nil {
			panic(err)
		}

		mongoSession, err := mongoClient.StartSession()
		if err != nil {
			return nil, err
		}

		database := mongoClient.Database(testDatabase)
		if database == nil {
			return nil, mongodbcacheadapters.ErrNilDatabase
		}

		collection := database.Collection(testCollection)
		if collection == nil {
			return nil, mongodbcacheadapters.ErrNilCollection
		}

		sessionAdapter, err := mongodbcacheadapters.NewSession(&mongoSession, collection, defaultTTL)
		if err != nil {
			return nil, err
		}
		return sessionAdapter, nil
	}
}

// newMongoDBTestSuite creates a new test suite with tests for MongoDB adapters and sessions.
func newMongoDBTestSuite(t *testing.T, defaultTTL time.Duration) *MongoDBTestSuite {
	var suite suite.Suite

	return &MongoDBTestSuite{
		Suite: &suite,
		CacheAdapterPartialTestSuite: &testutil.CacheAdapterPartialTestSuite{
			Suite:      &suite,
			DefaultTTL: defaultTTL,
			NewAdapter: newTestAdapterFunc(defaultTTL),
			NewSession: newTestSessionFunc(t, defaultTTL),
			SleepFunc:  testSleepFunc(),
		},
	}
}

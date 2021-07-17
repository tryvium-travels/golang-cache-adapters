// Copyright 2021 Tryvium Travels LTD
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
	"fmt"
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
	suite.Run(t, newMongoDBAdapterTestSuite(t, defaultTTL))
}

type MongoDBAdapterTestSuite struct {
	*suite.Suite
	*testutil.CacheAdapterPartialTestSuite
}

func (suite *MongoDBAdapterTestSuite) SetupSuite() {
	startLocalMongoDBServer()
}

func (suite *MongoDBAdapterTestSuite) TearDownSuite() {
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

		sessionAdapter, err := mongodbcacheadapters.NewSession(mongoSession, collection, defaultTTL)
		if err != nil {
			return nil, err
		}
		return sessionAdapter, nil
	}
}

// newMongoDBAdapterTestSuite creates a new test suite with tests for MongoDB adapters and sessions.
func newMongoDBAdapterTestSuite(t *testing.T, defaultTTL time.Duration) *MongoDBAdapterTestSuite {
	var suite suite.Suite

	return &MongoDBAdapterTestSuite{
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

func (suite *MongoDBAdapterTestSuite) TestNew_NilClient() {
	adapter, err := mongodbcacheadapters.New(nil, testDatabase, testCollection, testDefaultTTL)
	suite.Require().Nil(adapter, "Should be nil on nil mongo client")
	suite.Require().Error(err, "Should give error on nil mongo client")
}

func (suite *MongoDBAdapterTestSuite) TestNew_DisconnectedClient() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")

	err = client.Disconnect(context.Background())
	suite.Require().NoError(err, "Should not give error on disconnecting a mongo client since no operations has been performed after the connection")

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testDefaultTTL)
	suite.Require().Nil(adapter, "Should be nil on invalid (disconnected) mongo client")
	suite.Require().Error(err, "Should give error on invalid (disconnected) mongo client")
}

func (suite *MongoDBAdapterTestSuite) TestNew_InvalidDatabase() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")
	testDatabaseNameButInvalid := "" // fmt.Sprintf("%s:but-invalid", testDatabase)

	adapter, err := mongodbcacheadapters.New(client, testDatabaseNameButInvalid, testCollection, testutil.DummyTTL)
	suite.Require().Nil(adapter, "Should be nil on invalid database name")
	suite.Require().Error(err, "Should give error on invalid database name")
}

func (suite *MongoDBAdapterTestSuite) TestNew_EmptyDatabaseName() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")
	testDatabaseNameButEmpty := ""

	adapter, err := mongodbcacheadapters.New(client, testDatabaseNameButEmpty, testCollection, testutil.DummyTTL)
	suite.Require().Nil(adapter, "Should be nil on empty database name")
	suite.Require().Error(err, "Should give error on empty database name")
}

func (suite *MongoDBAdapterTestSuite) TestNew_WhiteDatabaseName() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")
	testDatabaseNameButWhite := "       "

	adapter, err := mongodbcacheadapters.New(client, testDatabaseNameButWhite, testCollection, testutil.DummyTTL)
	suite.Require().Nil(adapter, "Should be nil on white database name")
	suite.Require().Error(err, "Should give error on white database name")
}

func (suite *MongoDBAdapterTestSuite) TestNew_EmptyCollectionName() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")
	testCollectionNameButEmpty := ""

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollectionNameButEmpty, testutil.DummyTTL)
	suite.Require().Nil(adapter, "Should be nil on empty collection name")
	suite.Require().Error(err, "Should give error on empty collection name")
	suite.Require().ErrorIs(err, mongodbcacheadapters.ErrInvalidCollectionName)
}

func (suite *MongoDBAdapterTestSuite) TestNew_WhiteCollectionName() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")
	testCollectionNameButWhite := " "

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollectionNameButWhite, testutil.DummyTTL)
	suite.Require().Nil(adapter, "Should be nil on white collection name")
	suite.Require().Error(err, "Should give error on white collection name")
	suite.Require().ErrorIs(err, mongodbcacheadapters.ErrInvalidCollectionName)
}

func (suite *MongoDBAdapterTestSuite) TestNew_ZeroTTL() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.ZeroTTL)
	suite.Require().Nil(adapter, "Should be nil on zero TTL")
	suite.Require().Error(err, "Should give error on zero TTL")
}

func (suite *MongoDBAdapterTestSuite) TestNew_InvalidTTL() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.InvalidTTL)
	suite.Require().Nil(adapter, "Should be nil on invalid TTL")
	suite.Require().Error(err, "Should give error on invalid TTL")
}

func (suite *MongoDBAdapterTestSuite) TestOpenSessionOK() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	sessionAdapter, err := adapter.OpenSession()

	suite.Require().NoError(err, "Should not give error on valid adapter")
	suite.Require().NotNil(sessionAdapter, "Should not be nil on valid adapter")
}
func (suite *MongoDBAdapterTestSuite) TestOpenCloseSession_OK() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	sessionAdapter, err := adapter.OpenSession()

	suite.Require().NoError(err, "Should not give error on valid adapter")
	suite.Require().NotNil(sessionAdapter, "Should not be nil on valid adapter")

	err = sessionAdapter.Close()
	suite.Require().NoError(err, "Should not give error on closing a valid session adapter")
}

func (suite *MongoDBAdapterTestSuite) TestOpenSession_ErrOnDisconnectedClient() {
	//	client := newMockMongoClientAdapter()
	//	client.On("StartSession").Return(nil, testutil.ErrTestingFailureCheck).Once()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	client.Disconnect(context.Background())

	sessionAdapter, err := adapter.OpenSession()
	suite.Require().Error(err, "Should give error on invalid connection")
	suite.Require().Nil(sessionAdapter, "Should be nil on invalid connection")
}

func (suite *MongoDBAdapterTestSuite) TestDel_ErrMissing() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	testKeyForDeleteButInvalid := fmt.Sprintf("%s:but-invalid", testutil.TestKeyForDelete)
	err = adapter.Delete(testKeyForDeleteButInvalid)
	suite.Require().NoError(err, "Should not error on delete with non-existing key")
}

// test a subsequent Delete operation over the same key without Set between these two Delete.
// It differs from common/TestDelete_OK since it's the fail of the Delete operation
// that is tested
func (suite *MongoDBAdapterTestSuite) TestDel_DoubleDel() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	err = adapter.Set(testutil.TestKeyForDelete, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForDelete, &actual)
	suite.Require().NoError(err, "Value should be valid, hence no error")
	suite.Require().Equal(testutil.TestValue, actual, "The value just set must be equal to the test value")

	err = adapter.Delete(testutil.TestKeyForDelete)
	suite.Require().NoError(err, "Should not error on valid Delete")

	err = adapter.Delete(testutil.TestKeyForDelete)
	suite.Require().NoError(err, "Should not error on subsequent Delete on the same key solely by this")

}
func (suite *MongoDBAdapterTestSuite) TestTTL_SetDeleteExpires() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	duration := 250 * time.Millisecond

	err = adapter.Set(testutil.TestKeyForSetTTL, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should not error on valid set")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, duration)
	suite.Require().NoError(err, "Should not error on valid SetTTL")

	suite.SleepFunc(100 * time.Millisecond)

	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForSetTTL, &actual)
	suite.Require().Equal(testutil.TestValue, actual, "The value just set must be equal to the test value (is not expired)")
	suite.Require().NoError(err, "Should not error on valid get (is not expired)")

	suite.SleepFunc(200 * time.Millisecond)

	err = adapter.Delete(testutil.TestKeyForSetTTL)
	suite.Require().NoError(err, "Should not error on Delete after expires since it makes no differences")

}
func (suite *MongoDBAdapterTestSuite) TestTTL_SetNotExistingKey() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	testKeyForSetDLLButInvalid := fmt.Sprintf("%s:but-invalid", testutil.TestKeyForSetTTL)

	duration := new(time.Duration)
	*duration = time.Millisecond * 250
	err = adapter.SetTTL(testKeyForSetDLLButInvalid, *duration)
	suite.Require().Error(err, "Should error on set with not existing key")
	suite.Require().ErrorIs(err, cacheadapters.ErrNotFound, "Should error on set with not existing key")
}

func (suite *MongoDBAdapterTestSuite) TestTTL_SetOverExpired() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	duration := new(time.Duration)
	*duration = 50 * time.Millisecond
	err = adapter.Set(testutil.TestKeyForSetTTL, testutil.TestValue, duration)
	suite.Require().NoError(err, "Should not error on valid set")

	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForSetTTL, &actual)
	suite.Require().NoError(err, "Should not error on valid get (is not expired)")
	suite.Require().Equal(testutil.TestValue, actual, "The value just set must be equal to the test value (is not expired)")

	suite.SleepFunc(100 * time.Millisecond)

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, (*duration)*2)
	suite.Require().NoError(err, "Should not error on setting TTL over expired key, since it's removed")
	suite.Require().ErrorIs(err, nil, "Should not error on setting TTL over expired key, since it's removed")
}

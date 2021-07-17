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

	mongodbcacheadapters "github.com/tryvium-travels/golang-cache-adapters/mongodb"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

func (suite *MongoDBAdapterTestSuite) TestOpenSession_OK() {
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a valid adapter")
	suite.Require().NotNil(adapter, "Should be successfully created")

	sessionAdapter, err := adapter.OpenSession()
	suite.Require().NoError(err, "Should not error on creating a valid session adapter")
	suite.Require().NotNil(sessionAdapter, "Should be successfully created")
}

/*
func (suite *MongoDBAdapterTestSuite) TestNewSession_OK(){
	adapter, err := suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a valid adapter")
	suite.Require().NotNil(adapter, "Should be successfully created")


	mongoSession, err := adapter.client.StartSession()
	if err != nil {
		return nil, err
	}
}
*/

func (suite *MongoDBAdapterTestSuite) TestGet_ErrClientDisconnected() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	if err != nil {
		panic(err)
	}
	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	err = adapter.Set(testutil.TestKeyForGet, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Get")

	// disconnect the client, so that the OpenSession inside the Get will (and has to) fail
	err = client.Disconnect(context.Background())
	suite.Require().NoError(err, "Should not error on valid mongodb client Disconnect")

	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForGet, &actual)
	suite.Require().Equal(testutil.TestStruct{}, actual, "Should not be extracted value using a Get doomed to fail")
	suite.Require().Error(err, "Should return an error on Get with invalid connection")
	suite.Require().ErrorIs(err, mongodbcacheadapters.ErrSessionClosed, "The error should indicate a session closed")
}

func (suite *MongoDBAdapterTestSuite) TestSet_ErrClientDisconnected() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	if err != nil {
		panic(err)
	}
	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	// disconnect the client, so that the OpenSession inside the Set will (and has to) fail
	err = client.Disconnect(context.Background())
	suite.Require().NoError(err, "Should not error on valid mongodb client Disconnect")

	err = adapter.Set(testutil.TestKeyForSet, testutil.TestValue, nil)
	suite.Require().Error(err, "Should return an error on Set due to a closed connection causing OpenSession to fail")
	suite.Require().ErrorIs(err, mongodbcacheadapters.ErrSessionClosed, "The error should indicate a session closed")
}

func (suite *MongoDBAdapterTestSuite) TestDelete_ErrClientDisconnected() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	if err != nil {
		panic(err)
	}
	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	err = adapter.Set(testutil.TestKeyForDelete, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the Delete")

	// do a check just to be shure of the next tests, even if individual tests for Set exists
	var actual testutil.TestStruct
	err = adapter.Get(testutil.TestKeyForDelete, &actual)
	suite.Require().Equal(testutil.TestValue, actual, "The Set operation should be successful")
	suite.Require().NoError(err, "Should succeed the Set in order to test the Delete")

	// disconnect the client, so that the OpenSession inside the Delete will (and has to) fail
	err = client.Disconnect(context.Background())
	suite.Require().NoError(err, "Should not error on valid mongodb client Disconnect")

	err = adapter.Delete(testutil.TestKeyForDelete)
	suite.Require().Error(err, "Should return an error on Delete with invalid connection")
	suite.Require().ErrorIs(err, mongodbcacheadapters.ErrSessionClosed, "The error should indicate a session closed")

	adapter, err = suite.NewAdapter()
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	actual.Value = "some vole surely not stored"
	err = adapter.Get(testutil.TestKeyForDelete, &actual)
	suite.Require().NoError(err, "The entry should still be stored since the Delete should have failed")
	suite.Require().Equal(testutil.TestValue, actual, "The original value should be still stored")

}

func (suite *MongoDBAdapterTestSuite) TestSetTTL_ErrClientDisconnected() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	if err != nil {
		panic(err)
	}
	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid adapter.")

	err = adapter.Set(testutil.TestKeyForSetTTL, testutil.TestValue, nil)
	suite.Require().NoError(err, "Should perform the Set in order to test the SetTTL")

	// disconnect the client, so that the OpenSession inside the SetTTL will (and has to) fail
	err = client.Disconnect(context.Background())
	suite.Require().NoError(err, "Should not error on valid mongodb client Disconnect")

	err = adapter.SetTTL(testutil.TestKeyForSetTTL, testutil.DummyTTL)
	suite.Require().Error(err, "Should return an error on SetTTL with invalid connection")
	suite.Require().ErrorIs(err, mongodbcacheadapters.ErrSessionClosed, "The error should indicate a session closed")
}

func (suite *MongoDBAdapterTestSuite) TestOpenSession_InvalidCollectionName() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(client, "Should instantiate a valid mongo client")
	testCollectionNameButInvalid := "system. $ A║ÄîxFrR║Æ"

	adapter, err := mongodbcacheadapters.New(client, testDatabase, testCollectionNameButInvalid, testutil.DummyTTL)
	suite.Require().NotNil(adapter, "Should not be nil: the adapter is still created even if on invalid collection name")
	suite.Require().NoError(err, "Should not give error even if on invalid collection name")

	sessionAdapter, err := adapter.OpenSession()
	suite.Require().Nil(sessionAdapter, "Should be nil on invalid collection name")
	suite.Require().Error(err, "Should give error on invalid collection name")
}

// ----------------------------------

/*
func (suite *MongoDBAdapterTestSuite) TestGet_ErrNotDecodableStoredValue(){
	sessionAdapter, err := suite.NewSession()
	// TODO should insert into the database a corrupted "codificated" cacheItem so that the result.Decode() will fail
	// or interrupt the update process meanwhile the saving operation is running, leaving th database in a corrupted state
}
*/

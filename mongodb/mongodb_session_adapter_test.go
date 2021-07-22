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
	"time"

	"github.com/stretchr/testify/mock"
	mongodbcacheadapters "github.com/tryvium-travels/golang-cache-adapters/mongodb"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (suite *MongoDBAdapterTestSuite) TestNewSession_NilSession() {
	session, err := mongodbcacheadapters.NewSession(nil, testutil.DummyTTL)
	suite.Require().Nil(session, "Should be nil on invalid collection")
	suite.Require().Error(err, "Should give error on invalid collection")
}

func (suite *MongoDBAdapterTestSuite) TestNewSession_InvalidTTL() {
	session, err := mongodbcacheadapters.NewSession(newMockMongoCollection(nil, true, true, true, true), testutil.InvalidTTL)
	suite.Require().Nil(session, "Should be nil on invalid TTL")
	suite.Require().Error(err, "Should give error on invalid TTL")
}

func (suite *MongoDBAdapterTestSuite) TestSessionGet_DecodeError() {
	session, err := suite.NewSession()
	suite.Require().NoError(err, "Should not error on creating a new valid session adapter.")

	var actual complex128
	err = session.Get(testutil.TestKeyForGet, &actual)
	suite.Require().Error(err, "Should error on getting a non decodable value")
}

func (suite *MongoDBAdapterTestSuite) TestSessionGet_DecodeCacheStructError() {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(mongoClient, "Should instantiate a valid mongo client")

	collection := mongoClient.Database(testDatabase).Collection(testCollection)

	session, err := mongodbcacheadapters.NewSession(collection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should create a valid session adapter")

	invalidCacheItem := bson.M{
		"key":        testutil.TestKeyForGet,
		"expires_at": "INVALID",
	}
	_, err = collection.InsertOne(context.Background(), invalidCacheItem)
	suite.Require().NoError(err, "Must insert the invalid item for the test to work")

	var actual testutil.TestStruct
	err = session.Get(testutil.TestKeyForGet, &actual)
	suite.Require().Error(err, "Should error on getting a non decodable value")
}

func (suite *MongoDBAdapterTestSuite) TestSessionSet_EncodeError() {
	session, err := suite.NewSession()
	suite.Require().NoError(err, "Should not error on creating a new valid session adapter.")

	err = session.Set(testutil.TestKeyForGet, complex128(1), nil)
	suite.Require().Error(err, "Should error on setting a non encodable value")
}

func (suite *MongoDBAdapterTestSuite) TestSessionSet_UpdateError() {
	mockedCollection := newMockMongoCollection(nil, true, true, true, true)
	session, err := mongodbcacheadapters.NewSession(mockedCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid session adapter.")

	mockedCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, testutil.ErrTestingFailureCheck).Once()

	var actual testutil.TestStruct
	err = session.Set(testutil.TestKeyForSet, &actual, nil)
	suite.Require().Error(err, "Should error because of the mocked collection")
}

func (suite *MongoDBAdapterTestSuite) TestSessionSetTTL_FindError() {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(mongoClient, "Should instantiate a valid mongo client")

	collection := mongoClient.Database(testDatabase).Collection(testCollection)

	mockedCollection := newMockMongoCollection(collection, true, true, true, true)
	session, err := mongodbcacheadapters.NewSession(mockedCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid session adapter.")

	mockedCollection.On("FindOne", mock.Anything, mock.Anything).Return(nil, testutil.ErrTestingFailureCheck).Once()

	err = session.SetTTL(testutil.TestKeyForSetTTL, testutil.DummyTTL)
	suite.Require().Error(err, "Should error because of the mocked collection")
}

func (suite *MongoDBAdapterTestSuite) TestSessionSetTTL_UpdateError() {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(mongoClient, "Should instantiate a valid mongo client")

	collection := mongoClient.Database(testDatabase).Collection(testCollection)

	cacheItem := bson.M{
		"key":        testutil.TestKeyForSetTTL,
		"item":       testutil.TestValue,
		"expires_at": time.Now().Add(time.Minute),
	}
	_, err = collection.InsertOne(context.Background(), cacheItem)
	suite.Require().NoError(err, "Must insert the invalid item for the test to work")

	mockedCollection := newMockMongoCollection(collection, false, true, true, true)
	session, err := mongodbcacheadapters.NewSession(mockedCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid session adapter.")

	mockedCollection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, testutil.ErrTestingFailureCheck).Once()

	err = session.SetTTL(testutil.TestKeyForSetTTL, testutil.DummyTTL)
	suite.Require().Error(err, "Should error because of the mocked collection")
}

func (suite *MongoDBAdapterTestSuite) TestSessionSetTTL_DecodeCacheStructError() {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(mongoClient, "Should instantiate a valid mongo client")

	collection := mongoClient.Database(testDatabase).Collection(testCollection)

	session, err := mongodbcacheadapters.NewSession(collection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should create a valid session adapter")

	invalidCacheItem := bson.M{
		"key":        testutil.TestKeyForSetTTL,
		"expires_at": "INVALID",
	}
	_, err = collection.InsertOne(context.Background(), invalidCacheItem)
	suite.Require().NoError(err, "Must insert the invalid item for the test to work")

	err = session.SetTTL(testutil.TestKeyForSetTTL, testutil.DummyTTL)
	suite.Require().Error(err, "Should error on getting a non decodable value")
}

func (suite *MongoDBAdapterTestSuite) TestSessionDelete_DeleteError() {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	suite.Require().NoError(err, "Should not error on creating a valid mongo client")
	suite.Require().NotNil(mongoClient, "Should instantiate a valid mongo client")

	collection := mongoClient.Database(testDatabase).Collection(testCollection)

	mockedCollection := newMockMongoCollection(collection, true, true, true, true)
	session, err := mongodbcacheadapters.NewSession(mockedCollection, testutil.DummyTTL)
	suite.Require().NoError(err, "Should not error on creating a new valid session adapter.")

	mockedCollection.On("DeleteOne", mock.Anything, mock.Anything).Return(nil, testutil.ErrTestingFailureCheck).Once()

	err = session.Delete(testutil.TestKeyForDelete)
	suite.Require().Error(err, "Should error because of the mocked collection")
}

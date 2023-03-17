// Copyright 2023 Tryvium Travels LTD
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
	ctx "context"
	"time"

	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tryvium-travels/memongo"
	"github.com/tryvium-travels/memongo/memongolog"
)

const (
	mongoDBVersion string        = "4.2.1"
	testDatabase   string        = "test_database"
	testCollection string        = "test_collection"
	testDefaultTTL time.Duration = time.Millisecond * 50
)

var (
	testCacheIndexKeyName string = "test_cache_index_key"
	testCacheIndexTTLName string = "test_cache_index_ttl"
	localMongoDBServer    *memongo.Server
	testKeyIndexModel     = mongo.IndexModel{
		Keys: map[string]string{
			"key": "hashed",
		},
		Options: &options.IndexOptions{
			Name: &testCacheIndexKeyName,
		},
	}
	testTTLIndexModel = mongo.IndexModel{
		Keys: map[string]byte{
			"expires_at": 1,
		},
		Options: &options.IndexOptions{
			Name:               &testCacheIndexTTLName,
			ExpireAfterSeconds: new(int32), // default value 0, because expires after 0 seconds after the expiration time.
		},
	}
	testMongoOptions *memongo.Options = &memongo.Options{
		MongoVersion:   mongoDBVersion,
		StartupTimeout: 10 * time.Second,
		LogLevel:       memongolog.LogLevelSilent,
	}
)

func startLocalMongoDBServer() {
	var err error
	localMongoDBServer, err = memongo.StartWithOptions(testMongoOptions)
	if err != nil {
		panic(err)
	}

	createTestCollectionAndIndex()
}

func createTestCollectionAndIndex() {
	client, err := mongo.Connect(ctx.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	if err != nil {
		panic(err)
	}

	err = client.Database(testDatabase).CreateCollection(ctx.Background(), testCollection)
	if err != nil {
		panic(err)
	}

	indexes := client.Database(testDatabase).Collection(testCollection).Indexes()

	_, err = indexes.CreateOne(ctx.Background(), testKeyIndexModel)
	if err != nil {
		panic(err)
	}

	_, err = indexes.CreateOne(ctx.Background(), testTTLIndexModel)
	if err != nil {
		panic(err)
	}
}

func stopLocalMongoDBServer() {
	if localMongoDBServer != nil {
		localMongoDBServer.Stop()
		localMongoDBServer = nil
	}
}

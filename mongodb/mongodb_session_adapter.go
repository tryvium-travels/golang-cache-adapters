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

package mongodbcacheadapters

import (
	"context"
	"fmt"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	"go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBSessionAdapter struct {
	session    *mongo.Session    // The MongoDB session connected to this adapter.
	collection *mongo.Collection // The used MongoDB collection.
	defaultTTL time.Duration     // The defaultTTL of the Set operations.
}

type cacheItem struct {
	Key       string    `bson:"key"`        // The string key that identifies the item in cache
	Item      bson.Raw  `bson:"item"`       // The actual item in cache.
	ExpiresAt time.Time `bson:"expires_at"` // The expiration time of the item in cache.
}

// NesSession create a new MongoDB Session adapter
func NewSession(session *mongo.Session, collection *mongo.Collection, defaultTTL time.Duration) (cacheadapters.CacheSessionAdapter, error) {
	if session == nil {
		return nil, ErrNilSession
	}

	if collection == nil {
		return nil, ErrNilCollection
	}

	if defaultTTL < 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &MongoDBSessionAdapter{
		session:    session,
		collection: collection,
		defaultTTL: defaultTTL,
	}, nil
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (msa *MongoDBSessionAdapter) Close() error {
	return fmt.Errorf("TODO: CLOSE not yet implemented in MongoDB Session Adapter")
}

func (msa *MongoDBSessionAdapter) Get(key string, objectRef interface{}) error {
	result := msa.collection.FindOne(context.Background(), bson.M{"key": key})

	if result == nil {
		return cacheadapters.ErrNotFound
	}

	var valueFromMemory cacheItem

	err := result.Decode(valueFromMemory)
	if err != nil {
		return err
	}
	now := time.Now()
	if valueFromMemory.ExpiresAt.UnixNano() < now.UnixNano() {
		msa.Delete(key)
		return cacheadapters.ErrNotFound
	}

	err = bson.Unmarshal(valueFromMemory.Item, objectRef)
	if err != nil {
		return err
	}

	return nil
}

// Set sets a value represented by the object parameter into the cache,
// with the specified key.
func (msa *MongoDBSessionAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	marshalledObj, err := bson.Marshal(&object)
	if err != nil {
		return err
	}

	now := time.Now()
	expiresAt := now.Add(*TTL)

	optionsUpdate := options.Update().SetUpsert(true)
	filter := bson.M{"key": key}
	update := bson.M{
		"key":        key,
		"item":       marshalledObj,
		"expires_at": expiresAt,
	}

	_, err = msa.collection.UpdateOne(context.Background(), filter, update, optionsUpdate)
	if err != nil {
		return err
	}

	return nil
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (msa *MongoDBSessionAdapter) SetTTL(key string, newTTL time.Duration) error {
	if newTTL <= cacheadapters.TTLExpired {
		msa.Delete(key)
		return cacheadapters.ErrInvalidTTL
	}

	mongoResult := msa.collection.FindOne(context.Background(), bson.M{"key": key})
	if mongoResult == nil {
		return cacheadapters.ErrNotFound
	}

	var result cacheItem
	err := mongoResult.Decode(result)
	if err != nil {
		return err
	}

	now := time.Now()
	if result.ExpiresAt.UnixNano() < now.UnixNano() {
		msa.Delete(key)
		return cacheadapters.ErrNotFound
	}

	result.ExpiresAt = now.Add(newTTL)
	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"expires_at": result.ExpiresAt,
		},
	}
	_, err = msa.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a key from the cache.
func (msa *MongoDBSessionAdapter) Delete(key string) error {
	_, err := msa.collection.DeleteOne(context.Background(), bson.M{"key": key})
	if err != nil {
		return err
	}

	return nil
}
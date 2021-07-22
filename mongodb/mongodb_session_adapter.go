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

package mongodbcacheadapters

import (
	"context"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBSessionAdapter struct {
	collection MongoCollection // The used MongoDB collection.
	defaultTTL time.Duration   // The defaultTTL of the Set operations.
}

type cacheItem struct {
	Key       string    `bson:"key"`        // The string key that identifies the item in cache
	Item      bson.Raw  `bson:"item"`       // The actual item in cache.
	ExpiresAt time.Time `bson:"expires_at"` // The expiration time of the item in cache.
}

// NesSession create a new MongoDB Session adapter
func NewSession(collection MongoCollection, defaultTTL time.Duration) (cacheadapters.CacheSessionAdapter, error) {
	if collection == nil {
		return nil, ErrNilCollection
	}

	if defaultTTL <= 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &MongoDBSessionAdapter{
		collection: collection,
		defaultTTL: defaultTTL,
	}, nil
}

// Close closes the Cache Session.
func (msa *MongoDBSessionAdapter) Close() error {
	return nil
}

func (msa *MongoDBSessionAdapter) Get(key string, objectRef interface{}) error {
	if objectRef == nil {
		return cacheadapters.ErrGetRequiresObjectReference
	}

	result := msa.collection.FindOne(context.Background(), bson.M{"key": key})
	if result == nil || result.Err() != nil {
		return cacheadapters.ErrNotFound
	}

	var valueFromDB cacheItem

	err := result.Decode(&valueFromDB)
	if err != nil {
		return err
	}
	now := time.Now()
	if valueFromDB.ExpiresAt.UnixNano() < now.UnixNano() {
		msa.Delete(key)
		return cacheadapters.ErrNotFound
	}

	err = bson.Unmarshal(valueFromDB.Item, objectRef)
	if err != nil {
		return err
	}

	return nil
}

// Set sets a value represented by the object parameter into the cache,
// with the specified key.
func (msa *MongoDBSessionAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	if TTL == nil {
		TTL = &msa.defaultTTL
	}

	if *TTL <= 0 {
		return cacheadapters.ErrInvalidTTL
	}

	marshalledObj, err := bson.Marshal(&object)
	if err != nil {
		return err
	}

	now := time.Now()
	expiresAt := now.Add(*TTL)

	optionsUpdate := options.Update().SetUpsert(true)
	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"key":        key,
			"item":       bson.Raw(marshalledObj),
			"expires_at": expiresAt,
		},
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
		return nil
	}

	mongoResult := msa.collection.FindOne(context.Background(), bson.M{"key": key})
	if mongoResult == nil || mongoResult.Err() != nil {
		return cacheadapters.ErrNotFound
	}

	var result cacheItem
	err := mongoResult.Decode(&result)
	if err != nil {
		return err
	}

	now := time.Now()
	if result.ExpiresAt.UnixNano() < now.UnixNano() {
		msa.Delete(key)
		return nil
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

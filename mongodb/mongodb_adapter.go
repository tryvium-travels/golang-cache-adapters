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
	"strings"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

type MongoDBAdapter struct {
	client         MongoClient   // The MongoDB client interface to interact with the collection.
	databaseName   string        // The name of the database used in MongoDB to cache data.
	collectionName string        // The name of the collection used in MongoDB to cache data.
	defaultTTL     time.Duration // The defaultTTL of the Set operations.
}

// NesSession create a new MongoDB Cache adapter from an existing
// MongoDB client and the name of the database and the collection,
// with a given default TTL.
func New(client MongoClient, databaseName string, collectionName string, defaultTTL time.Duration) (cacheadapters.CacheAdapter, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	databaseName = strings.TrimSpace(databaseName)
	if databaseName == "" {
		return nil, ErrInvalidDatabaseName
	}

	collectionName = strings.TrimSpace(collectionName)
	if collectionName == "" {
		return nil, ErrInvalidCollectionName
	}

	if defaultTTL <= 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &MongoDBAdapter{
		client:         client,
		databaseName:   databaseName,
		collectionName: collectionName,
		defaultTTL:     defaultTTL,
	}, nil
}

func (ma *MongoDBAdapter) OpenSession() (cacheadapters.CacheSessionAdapter, error) {
	mongoSession, err := ma.client.StartSession()
	if err != nil {
		return nil, err
	}

	mongoDatabase := ma.client.Database(ma.databaseName)
	mongoCollection := mongoDatabase.Collection(ma.collectionName)

	return NewSession(mongoSession, mongoCollection, ma.defaultTTL)
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (ma *MongoDBAdapter) Get(key string, objectRef interface{}) error {
	msa, err := ma.OpenSession()
	if err != nil {
		return err
	}

	defer msa.Close()

	return msa.Get(key, objectRef)
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (ma *MongoDBAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	rsa, err := ma.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Set(key, object, TTL)
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (ma *MongoDBAdapter) SetTTL(key string, newTTL time.Duration) error {
	rsa, err := ma.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.SetTTL(key, newTTL)
}

// Delete deletes a key from the cache.
func (ma *MongoDBAdapter) Delete(key string) error {
	rsa, err := ma.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Delete(key)
}

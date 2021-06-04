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

package rediscacheadapters

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// RedisSessionAdapter is the CacheSessionAdapter implementation
// for Redis.
type RedisSessionAdapter struct {
	conn          redis.Conn    // The redis connection used to connect.
	defaultTTL    time.Duration // The defaultTTL of the Set operations.
	inTransaction bool          // True if inside a transaction.
	mutex         *sync.Mutex   // mutex to handle transactions.
}

// NewSession creates a new Redis Cache Session adapter from
// an existing Redis connection.
func NewSession(conn redis.Conn, defaultTTL time.Duration) (cacheadapters.CacheSessionAdapter, error) {
	if conn == nil {
		return nil, cacheadapters.ErrInvalidConnection
	}

	if defaultTTL < 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &RedisSessionAdapter{
		conn:          conn,
		defaultTTL:    defaultTTL,
		inTransaction: false,
		mutex:         &sync.Mutex{},
	}, nil
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (rsa *RedisSessionAdapter) Get(key string, objectRef interface{}) error {
	if rsa.inTransaction {
		return rsa.conn.Send("GET", key)
	}

	resultContent, err := redis.Bytes(rsa.conn.Do("GET", key))
	if err == redis.ErrNil {
		return cacheadapters.ErrNotFound
	}

	if err != nil {
		return err
	}

	if objectRef == nil {
		return cacheadapters.ErrGetRequiresObjectReference
	}

	err = json.Unmarshal(resultContent, objectRef)
	if err != nil {
		return err
	}

	return nil
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (rsa *RedisSessionAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	if TTL == nil {
		TTL = new(time.Duration)
		*TTL = rsa.defaultTTL
	} else if *TTL <= 0 {
		return cacheadapters.ErrInvalidTTL
	}

	objectContent, err := json.Marshal(object)
	if err != nil {
		return err
	}

	if rsa.inTransaction {
		return rsa.conn.Send("SETEX", key, (*TTL).Seconds(), objectContent)
	}

	_, err = rsa.conn.Do("SETEX", key, (*TTL).Seconds(), objectContent)
	if err != nil {
		return err
	}

	return nil
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (rsa *RedisSessionAdapter) SetTTL(key string, newTTL time.Duration) error {
	var err error

	if newTTL > cacheadapters.TTLExpired {
		_, err = rsa.conn.Do("EXPIRE", key, newTTL.Seconds())
	} else {
		return rsa.Delete(key)
	}

	return err
}

// Delete deletes a key from the cache.
func (rsa *RedisSessionAdapter) Delete(key string) error {
	_, err := rsa.conn.Do("DEL", key)

	return err
}

// InTransaction allows to execute multiple Cache Sets and Gets in a Transaction, then tries to
// Unmarshal the array of results into the specified array of object references.
func (rsa *RedisSessionAdapter) InTransaction(inTransactionFunc cacheadapters.InTransactionFunc, objectRefs []interface{}) error {
	if rsa.inTransaction {
		return cacheadapters.ErrNoNestedTransactions
	}

	rsa.mutex.Lock()
	rsa.inTransaction = true
	defer func() {
		rsa.mutex.Lock()
		rsa.inTransaction = false
		rsa.mutex.Unlock()
	}()
	rsa.mutex.Unlock()

	if inTransactionFunc == nil {
		return nil
	}

	err := rsa.conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = inTransactionFunc(rsa)
	if err != nil {
		rsa.conn.Do("DISCARD")
		return err
	}

	transactionResults, err := redis.Values(rsa.conn.Do("EXEC"))
	if err != nil {
		return err
	}

	if objectRefs == nil {
		return cacheadapters.ErrGetRequiresObjectReference
	}

	if len(objectRefs) != len(transactionResults) {
		return cacheadapters.ErrInTransactionObjectReferencesLengthMismatch
	}

	for i, transactionResult := range transactionResults {
		if transactionResult == nil {
			if objectRefs[i] == nil {
				continue
			}
			return cacheadapters.ErrInTransactionMarshalValue
		}

		result, ok := transactionResult.([]byte)
		if !ok {
			resultString, ok := transactionResult.(string)
			isOkString := resultString == "OK"
			if ok {
				if isOkString {
					if objectRefs[i] == nil {
						continue
					} else {
						return cacheadapters.ErrInTransactionMarshalValue
					}
				} else {
					result = []byte(resultString)
				}
			} else {
				return cacheadapters.ErrInTransactionMarshalValue
			}
		}

		err := json.Unmarshal([]byte(result), objectRefs[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the Cache Session.
func (rsa *RedisSessionAdapter) Close() error {
	return rsa.conn.Close()
}

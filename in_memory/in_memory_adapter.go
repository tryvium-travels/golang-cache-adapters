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

package inmemorycacheadapters

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// cacheItem is the internal struct
// handling the mechanism of cache expiration.
type cacheItem struct {
	item      json.RawMessage // The actual item in cache.
	expiresAt time.Time       // The expiration time of the item in cache.
}

// cacheData is the container of all the in-memory
// cache used by the adapter.
type cacheData map[string]cacheItem

// InMemoryAdapter is the cache adapter which uses internal memory
// of the process.
type InMemoryAdapter struct {
	defaultTTL time.Duration // The defaultTTL of the Set operations.
	data       cacheData     // The data being stored in the in-memory cache.
	mutex      sync.Mutex    // The mutex locking the operations.
}

// stackNode is a node of a linked stack
type stackNode struct {
	element interface{}
	next    *stackNode
}

// simpleStack is a simple linked-stack of interfaces
type simpleStack struct {
	mutex *sync.Mutex // The mutex locking the operations.
	size  int
	top   *stackNode
}

// inTransactionWrapperIMA is a wrapper of InMemoryAdapter
// which implements the InTransaction operation by observing
// every action and collecting the results.
type inTransactionWrapperIMA struct {
	ima             *InMemoryAdapter // delegator of all Adapter actions
	processedValues *simpleStack     // stack of processed values by actions
}

// New creates a new InMemoryAdapter from an default TTL.
func New(defaultTTL time.Duration) (*InMemoryAdapter, error) {
	if defaultTTL <= 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &InMemoryAdapter{
		defaultTTL: defaultTTL,
		data:       make(cacheData),
	}, nil
}

// OpenSession opens a new Cache Session.
// Returns the same adapter because the
// session with the memory is always open.
func (ima *InMemoryAdapter) OpenSession() (*InMemoryAdapter, error) {
	return ima, nil
}

// Close closes the Cache Session.
// Returns nil because the session with
// the memory is always on and does not
// need to be closed.
func (ima *InMemoryAdapter) Close() error {
	return nil
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (ima *InMemoryAdapter) Get(key string, resultRef interface{}) error {
	ima.mutex.Lock()
	valueFromMemory, exists := ima.data[key]
	ima.mutex.Unlock()
	if !exists {
		return cacheadapters.ErrNotFound
	}

	now := time.Now()
	if valueFromMemory.expiresAt.Unix() < now.Unix() {
		ima.Delete(key)
		return cacheadapters.ErrNotFound
	}

	err := json.Unmarshal(valueFromMemory.item, resultRef)
	if err != nil {
		return err
	}

	return nil
}

// Set sets a value represented by the object parameter into the cache,
// with the specified key.
func (ima *InMemoryAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	if TTL == nil {
		TTL = new(time.Duration)
		*TTL = ima.defaultTTL
	} else if *TTL <= 0 {
		return cacheadapters.ErrInvalidTTL
	}

	now := time.Now()
	expiresAt := now.Add(*TTL)

	content, err := json.Marshal(object)
	if err != nil {
		return err
	}

	ima.mutex.Lock()
	ima.data[key] = cacheItem{
		item:      content,
		expiresAt: expiresAt,
	}
	ima.mutex.Unlock()

	return nil
}

// SetTTL marks the specified key new expiration, deletes it via using
// cacheadapters.TTLExpired or negative duration.
func (ima *InMemoryAdapter) SetTTL(key string, newTTL time.Duration) error {
	if newTTL < 0 {
		return cacheadapters.ErrInvalidTTL
	}
	if newTTL == cacheadapters.TTLExpired {
		return ima.Delete(key)
	}

	ima.mutex.Lock()
	valueFromMemory, exists := ima.data[key]
	ima.mutex.Unlock()
	if !exists {
		return cacheadapters.ErrNotFound
	}

	now := time.Now()
	newExpiresAt := now.Add(newTTL)

	if valueFromMemory.expiresAt.Unix() < now.Unix() {
		ima.Delete(key)
	}

	valueFromMemory.expiresAt = newExpiresAt

	ima.mutex.Lock()
	ima.data[key] = valueFromMemory
	ima.mutex.Unlock()

	return nil
}

// Delete deletes a key from the cache.
func (ima *InMemoryAdapter) Delete(key string) error {
	ima.mutex.Lock()
	delete(ima.data, key)
	ima.mutex.Unlock()
	return nil
}

// InTransaction allows to execute multiple Cache Sets and Gets in a Transaction, then tries to
// Unmarshal the array of results into the specified array of object references.
func (ima *InMemoryAdapter) InTransaction(inTransactionFunc cacheadapters.InTransactionFunc, objectRefs []interface{}) error {
	wrapper := &inTransactionWrapperIMA{
		ima:             ima,
		processedValues: newStack(),
	}
	err := inTransactionFunc(wrapper)
	processedValues := wrapper.processedValues.toArray()
	// traspose the processed values into objectReferences
	minLength := 0
	if len(processedValues) < len(objectRefs) {
		minLength = len(processedValues)
	} else {
		minLength = len(objectRefs)
	}
	i := 0
	for ; i < minLength; i++ {
		objectRefs[i] = processedValues
	}
	minLength = len(objectRefs)
	for ; i < minLength; i++ {
		objectRefs[i] = nil
	}

	return err
}

func newStack() *simpleStack {
	return &simpleStack{
		mutex: &sync.Mutex{},
		size:  0,
		top:   nil,
	}
}

func (stack *simpleStack) push(obj interface{}) {
	newNode := &stackNode{
		element: obj,
		next:    nil,
	}
	fmt.Printf("adding: %+v\n", obj)
	stack.mutex.Lock()
	stack.size++
	newNode.next = stack.top
	stack.top = newNode
	stack.mutex.Unlock()
}

func (stack *simpleStack) toArray() []interface{} {
	stack.mutex.Lock()
	collectedArguments := make([]interface{}, stack.size)
	var node **stackNode = &stack.top
	i := 0
	for (*node) != nil { // equivalent to: "i < stack.size"
		collectedArguments[(stack.size-i)-1] = *((*node).element.(*interface{}))
		node = &((*node).next)
		i++
	}
	stack.mutex.Unlock()
	return collectedArguments
}

func (wrapper *inTransactionWrapperIMA) OpenSession() (*InMemoryAdapter, error) {
	return wrapper.ima.OpenSession()
}
func (wrapper *inTransactionWrapperIMA) Close() error {
	return wrapper.ima.Close()
}

func (wrapper *inTransactionWrapperIMA) Get(key string, resultRef interface{}) error {
	err := wrapper.ima.Get(key, resultRef)
	if err == nil {
		wrapper.processedValues.push(resultRef)
	}
	return err
}

func (wrapper *inTransactionWrapperIMA) Set(key string, object interface{}, TTL *time.Duration) error {
	err := wrapper.ima.Set(key, object, TTL)
	if err == nil {
		wrapper.processedValues.push(&object)
	}
	return err
}

func (wrapper *inTransactionWrapperIMA) SetTTL(key string, newTTL time.Duration) error {
	return wrapper.ima.SetTTL(key, newTTL)
}

func (wrapper *inTransactionWrapperIMA) Delete(key string) error {
	return wrapper.ima.Delete(key)
}

func (wrapper *inTransactionWrapperIMA) InTransaction(inTransactionFunc cacheadapters.InTransactionFunc, objectRefs []interface{}) error {
	return inTransactionFunc(wrapper)
}

package cacheadapters

import (
	"time"
)

// CacheAdapter represents a Cache Mechanism abstraction.
type CacheAdapter interface {
	// OpenSession opens a new Cache Session.
	OpenSession() CacheSessionAdapter

	CacheSessionAdapter
}

// CacheSessionAdapter represents a Cache Session Mechanism abstraction.
type CacheSessionAdapter interface {
	// Get obtains a value from the cache using a key, then tries to unmarshal
	// it into the object reference passed as parameter.
	Get(key string, objectRef interface{}) error

	// Set sets a value represented by the object parameter into the cache, with the specified key.
	Set(key string, object interface{}, TTL *time.Duration) error

	// InTransaction allows to execute multiple Cache Sets and Gets in a Transaction, then tries to
	// Unmarshal the array of results into the specified array of object references.
	InTransaction(inTransactionFunc func(adapter CacheSessionAdapter) error, objectRefs []interface{}) error

	// Close closes the Cache Session.
	Close() error
}

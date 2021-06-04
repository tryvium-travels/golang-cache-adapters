package rediscacheadapters

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// RedisAdapter is the CacheAdapter implementation for Redis.
type RedisAdapter struct {
	pool       *redis.Pool   // The Redis pool used to create connections.
	defaultTTL time.Duration // The defaultTTL of the Set operations.
}

// New creates a new RedisAdapter from an initialized Redis pool.
func New(pool *redis.Pool, defaultTTL time.Duration) (*RedisAdapter, error) {
	if pool == nil {
		return nil, fmt.Errorf("the Redis Pool cannot be nil")
	}

	if defaultTTL <= 0 {
		return nil, cacheadapters.ErrInvalidTTL
	}

	return &RedisAdapter{
		pool:       pool,
		defaultTTL: defaultTTL,
	}, nil
}

// OpenSession opens a new Cache Session.
func (ra *RedisAdapter) OpenSession() (cacheadapters.CacheSessionAdapter, error) {
	conn, err := ra.pool.Dial()
	if err != nil {
		return nil, err
	}

	return &RedisSessionAdapter{
		conn:          conn,
		defaultTTL:    ra.defaultTTL,
		inTransaction: false,
	}, nil
}

// Get obtains a value from the cache using a key, then tries to unmarshal
// it into the object reference passed as parameter.
func (ra *RedisAdapter) Get(key string, objectRef interface{}) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Get(key, objectRef)
}

// Set sets a value represented by the object parameter into the cache, with the specified key.
func (ra *RedisAdapter) Set(key string, object interface{}, TTL *time.Duration) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.Set(key, object, TTL)
}

// InTransaction allows to execute multiple Cache Sets and Gets in a Transaction, then tries to
// Unmarshal the array of results into the specified array of object references.
func (ra *RedisAdapter) InTransaction(inTransactionFunc func(adapter cacheadapters.CacheSessionAdapter) error, objectRefs []interface{}) error {
	rsa, err := ra.OpenSession()
	if err != nil {
		return err
	}

	defer rsa.Close()

	return rsa.InTransaction(inTransactionFunc, objectRefs)
}

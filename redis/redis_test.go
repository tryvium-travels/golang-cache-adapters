package rediscacheadapters_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// testStruct is just an example struct to check if the json
// marchalling and unmarshalling are correct in all tests.
type testStruct struct {
	Value string `json:"value"`
}

var (
	localRedisServer *miniredis.Miniredis      // The local in-memory redis instance
	testRedisPool    *redis.Pool               // The pool used in all the tests, except for the "InvalidPool" ones.
	invalidRedisPool *redis.Pool               // The pool used when in need to test invalid connection behaviours.
	testKeyForGet    = "test:key:for-get:1234" // The test key used to test the Get operations
	testKeyForSet    = "test:key:for-set:1234" // The test key used to test the Set operations
	testValue        = testStruct{"1"}         // The test value being Set
)

// startLocalRedisServer starts a local, in-memory redis instance for the tests.
func startLocalRedisServer() {
	var err error

	invalidRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return nil, errors.New("I fail perché si")
		},
	}

	localRedisServer, err = miniredis.Run()
	if err != nil {
		log.Fatalf("Cannot start local redis server: %s", err)
	}

	localRedisServer.Select(0)

	testValueContent, err := json.Marshal(testValue)
	if err != nil {
		log.Fatalf("Cannot set initial testKeyForGet on local redis: %s", err)
	}

	// set initial value for testKeyForGet
	err = localRedisServer.Set(testKeyForGet, string(testValueContent))
	if err != nil {
		log.Fatalf("Cannot set initial testKeyForGet on local redis: %s", err)
	}

	testRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", localRedisServer.Addr())
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}
}

// stopLocalRedisServer stops the previously started, local, in-memory Redis instance.
func stopLocalRedisServer() {
	localRedisServer.Close()
}

func inTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	err := session.Get(testKeyForGet, nil)
	if err != nil {
		return err
	}

	testValue2 := testStruct{
		Value: "2",
	}

	err = session.Set(testKeyForSet, testValue2, nil)
	if err != nil {
		return err
	}

	return nil
}

func erroringInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	return fmt.Errorf("THIS IS AN ERROR TO PROVE ERRORING FUNCTIONS ALSO WORK IN TESTS")
}

// TestMain adds Global test setups and teardowns.
func TestMain(m *testing.M) {
	startLocalRedisServer()
	defer stopLocalRedisServer()

	code := m.Run()

	os.Exit(code)
}

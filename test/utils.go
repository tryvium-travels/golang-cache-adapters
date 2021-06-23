package testutil

import "time"

var (
	TestKeyForGet    = "test:key:for-get:1234"     // The test key used to test the Get operations
	TestKeyForSet    = "test:key:for-set:1234"     // The test key used to test the Set operations
	TestKeyForSetTTL = "test:key:for-set-ttl:1234" // The test key used to test the SetTTL operations
	TestKeyForDelete = "test:key:for-delete:1234"  // The test key used to test the Delete operations
	TestValue        = TestStruct{"1"}             // The test value being Set
	TestValueJSON    = []byte(`{"value":"1"}`)     // The Test value as JSON string
)

// TestStruct is just an example struct to check if the json
// marchalling and unmarshalling are correct in all tests.
type TestStruct struct {
	Value string `json:"value"`
}

var (
	// DummyTTL represents a dummy TTL to put in
	// test operations when needed.
	DummyTTL time.Duration = time.Second
	// InvalidTTL represents a dummy TTL to put in
	// test operations when needed.
	InvalidTTL time.Duration = -time.Second
	// ZeroTTL represents a dummy TTL initialized to
	// a zero value to put in test operations when
	// needed.
	ZeroTTL time.Duration = 0
)

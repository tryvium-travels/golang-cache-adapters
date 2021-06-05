package testutil

import (
	"fmt"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

// InTransactionFunc is a well formed inTransaction func for tests.
func InTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	err := session.Get(TestKeyForGet, nil)
	if err != nil {
		return err
	}

	testValue2 := TestStruct{
		Value: "2",
	}

	err = session.Set(TestKeyForSet, testValue2, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetDelInTransactionFunc is a well formed inTransaction func for tests which does
// Get and Delete.
func GetDelInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	err := session.Get(TestKeyForGet, nil)
	if err != nil {
		return err
	}

	err = session.Delete(TestKeyForSet)
	if err != nil {
		return err
	}

	return nil
}

// NestedInTransactionFunc is a well formed inTransaction func for tests involving the
// nesting of transactions.
func NestedInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	return session.InTransaction(func(session cacheadapters.CacheSessionAdapter) error {
		return nil
	}, nil)
}

func ErroringInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	return fmt.Errorf("THIS IS AN ERROR TO PROVE ERRORING FUNCTIONS ALSO WORK IN TESTS")
}

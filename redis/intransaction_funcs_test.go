package rediscacheadapters_test

import (
	"fmt"
	"time"

	cacheadapters "github.com/tryvium-travels/golang-cache-adapters"
)

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

func getDelInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	err := session.Get(testKeyForGet, nil)
	if err != nil {
		return err
	}

	err = session.Delete(testKeyForSet)
	if err != nil {
		return err
	}

	return nil
}

func nestedInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	return session.InTransaction(func(session cacheadapters.CacheSessionAdapter) error {
		return nil
	}, nil)
}

func setGetExFloat64InTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	err := session.Set(testKeyForSet, 2.5, nil)
	if err != nil {
		return err
	}

	err = session.Get(testKeyForSet, nil)
	if err != nil {
		return err
	}

	err = session.SetTTL(testKeyForSet, time.Second)
	if err != nil {
		return err
	}

	return nil
}

func erroringInTransactionFunc(session cacheadapters.CacheSessionAdapter) error {
	return fmt.Errorf("THIS IS AN ERROR TO PROVE ERRORING FUNCTIONS ALSO WORK IN TESTS")
}

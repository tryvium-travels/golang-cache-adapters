package cacheadapters

import "fmt"

var (
	//ErrInvalidConnection will come out if you try to use an invalid connection in a session.
	ErrInvalidConnection = fmt.Errorf("cannot use an invalid connection")
	// ErrNotFound will come out if a key is not found in the cache.
	ErrNotFound = fmt.Errorf("the value tried to get has not been found, check if it may be expired")
	// ErrGetRequiresObjectReference will come out if a nil object
	// reference is passed in a Get operation.
	ErrGetRequiresObjectReference = fmt.Errorf("in Get operations it is mandatory to provide a non-nil object reference to store the result in, nil found")
	// ErrInTransactionObjectReferencesLengthMismatch will come out
	// if there is a mismatch in number of commands in the transaction
	// and the length of the object references array.
	ErrInTransactionObjectReferencesLengthMismatch = fmt.Errorf("in InTransactions you must provide an array of reference objects with length equal to the number of commands you call in the transaction")
	// ErrInvalidTTL will come out if you try to set a zero-or-negative
	// TTL in a Set operation.
	ErrInvalidTTL = fmt.Errorf("cannot provide a negative TTL to Set operations")
)

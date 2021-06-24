package rediscacheadapters

import "fmt"

var (
	//ErrInvalidConnection will come out if you try to use an invalid connection in a session.
	ErrInvalidConnection = fmt.Errorf("cannot use an invalid connection")
)

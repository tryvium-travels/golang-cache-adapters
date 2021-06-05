package testutil

import "fmt"

var (
	// ErrTestingFailureCheck will come out if returned by a test mock which
	// errors on purpose.
	ErrTestingFailureCheck = fmt.Errorf("TEST ERROR: THIS FAILS BECAUSE IT IS INTENDED TO")
)

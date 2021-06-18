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
	// ErrInvalidTTL will come out if you try to set a zero-or-negative
	// TTL in a Set operation.
	ErrInvalidTTL = fmt.Errorf("cannot provide a negative TTL to Set operations")
)

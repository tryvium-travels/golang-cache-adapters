// Copyright 2023 Tryvium Travels LTD
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

package multicacheadapters

import "fmt"

var (
	// ErrInvalidSubAdapters will come out if you try to pass one or more nil
	// sub-adapters when creating a new MultiCacheAdapter or a session.
	ErrInvalidSubAdapters = fmt.Errorf("you must pass at least one valid sub-adapter")
	// ErrMultiCacheWarning will come out paired with other errors in case
	// an non-fatal error occurs during a multicache operation.
	//
	// This includes for example when a GET operation fails on the first
	// adapter but is successful in the second adapter.
	ErrMultiCacheWarning = fmt.Errorf("warning when performing an operation with a multicache adapter")
)

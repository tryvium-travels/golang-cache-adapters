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

package mongodbcacheadapters

import "fmt"

var (
	//ErrDBConnectionNotCreated will come out if you try to create a database helper
	// over a database connection that is not yet created.
	ErrDBConnectionNotCreated = fmt.Errorf("cannot create a database helper instance over a never-created database connection")

	//ErrNilClient will come out if you try to create a CacheAdapter instance
	// when providing a nil MongoDB Client instance
	ErrNilClient = fmt.Errorf("cannot create the adapter with nil MongoDB Client instance")

	//ErrInvalidDatabaseName will come out if you try to create a CacheAdapter instance
	// when providing an invalid MongoDB Database name
	ErrInvalidDatabaseName = fmt.Errorf("cannot create the adapter with invalid MongoDB Database name")

	//ErrInvalidCollectionName will come out if you try to create a CacheAdapter instance
	// when providing an invalid MongoDB Collection name
	ErrInvalidCollectionName = fmt.Errorf("cannot create the adapter with invalid MongoDB Collection name")

	//ErrNilDatabase will come out if you try to create a CacheAdapter instance
	// when providing a nil MongoDB Database instance
	ErrNilDatabase = fmt.Errorf("cannot create the adapter with nil MongoDB Database instance")

	//ErrNilCollection will come out if you try to create a CacheAdapter instance
	// when providing a nil MongoDB Collection instance
	ErrNilCollection = fmt.Errorf("cannot create the adapter with nil MongoDB Collection instance")

	//ErrNilSession will come out if you try to create a CacheSessionAdapter instance
	// when providing a nil MongoDB Session instance
	ErrNilSession = fmt.Errorf("cannot create the session adapter with nil MongoDB Session instance")

	//ErrSessionClosed will come out if you try to do operation on an already
	// closed session
	ErrSessionClosed = fmt.Errorf("cannot use a closed connection")
)

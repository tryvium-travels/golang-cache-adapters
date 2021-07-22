// Copyright 2021 Tryvium Travels LTD
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

// Package cacheadapters contains the generic entities used to access
// the cache.
//
// In this package you will find the CacheAdapter and CacheSessionAdapter
// interfaces, used by the specific implementations to access the cache,
// along with a MultiCacheProvider struct you can use to access multiple
// cache adapters with an index-based priority mechanism.
//
//    If you want to see the specific implementations go to the folder with the
//    name of the implementation you are searching (e.g. "redis").
package cacheadapters

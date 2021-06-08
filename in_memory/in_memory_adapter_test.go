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

package inmemorycacheadapters_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	inmemorycacheadapters "github.com/tryvium-travels/golang-cache-adapters/in_memory"
	testutil "github.com/tryvium-travels/golang-cache-adapters/test"
)

// TestMain adds Global test setups and teardowns.
func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestNewInMemoryCacheAdapterOK(t *testing.T) {
	_, err := inmemorycacheadapters.New(testutil.DummyTTL)
	require.NoError(t, err, "Should not give error on valid New")
}

func TestNewInMemoryCacheAdapterWithInvalidTTL(t *testing.T) {
	_, err := inmemorycacheadapters.New(testutil.InvalidTTL)
	require.Error(t, err, "Should give error on valid New with invalid TTL")
}

func TestNewInMemoryCacheAdapterWithInvalidTTLZero(t *testing.T) {
	_, err := inmemorycacheadapters.New(testutil.ZeroTTL)
	require.Error(t, err, "Should give error on valid New with zero TTL")
}

func TestOpenCloseSessionOK(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)
	session, _ := adapter.OpenSession()

	err := session.Close()
	require.NoError(t, err, "Should not error on closing an open session")
}

func TestSetGetOK(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)

	err := adapter.Set(testutil.TestKeyForSet, testutil.TestValue, &testutil.DummyTTL)
	require.NoError(t, err, "Should not error on valid set")

	var actual testutil.TestStruct

	err = adapter.Get(testutil.TestKeyForSet, &actual)
	require.NoError(t, err, "Should not error on valid get")

	require.EqualValues(t, testutil.TestValue, actual, "The value obatined with get should be equal to the test value set before")
}

func TestSetWithNonUnmarshalableReference(t *testing.T) {
	adapter, _ := inmemorycacheadapters.New(testutil.DummyTTL)

	var actual complex128

	err := adapter.Get(testutil.TestKeyForSet, &actual)
	require.Error(t, err, "Should error on non unmarshalable get")
}

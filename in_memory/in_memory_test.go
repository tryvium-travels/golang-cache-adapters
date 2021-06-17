package inmemorycacheadapters_test

import (
	"os"
	"testing"
)

// TestMain adds Global test setups and teardowns.
func TestMain(m *testing.M) {

	code := m.Run()

	os.Exit(code)
}

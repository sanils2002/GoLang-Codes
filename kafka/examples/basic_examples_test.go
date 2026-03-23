package examples

import (
	"testing"
)

func TestBasicRetryExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BasicRetryExample panicked: %v", r)
		}
	}()

	// Run the example
	BasicRetryExample()
}

func TestHighPriorityRetryExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("HighPriorityRetryExample panicked: %v", r)
		}
	}()

	// Run the example
	HighPriorityRetryExample()
}

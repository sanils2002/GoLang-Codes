package examples

import (
	"testing"
)

func TestConnectionRetryExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ConnectionRetryExample panicked: %v", r)
		}
	}()

	// Run the example
	ConnectionRetryExample()
}

func TestDifferentUseCasesExample(t *testing.T) {
	// Test that the function runs without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DifferentUseCasesExample panicked: %v", r)
		}
	}()

	// Run the example
	DifferentUseCasesExample()
}

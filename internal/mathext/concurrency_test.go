package mathext

import (
	"testing"
)

func TestTopLevelConcurrency(t *testing.T) {
	tests := []struct {
		maxConcurrency int
		numberOfChecks int
		expectedResult int
	}{
		{10, 5, 5},
		{5, 10, 5},
		{0, 0, 0},
	}

	for _, test := range tests {
		result := TopLevelConcurrency(test.maxConcurrency, test.numberOfChecks)
		if result != test.expectedResult {
			t.Errorf("TopLevelConcurrency(%d, %d) = %d; want %d", test.maxConcurrency, test.numberOfChecks, result, test.expectedResult)
		}
	}
}

func TestInternalLevelConcurrency(t *testing.T) {
	tests := []struct {
		maxConcurrency         int
		minInternalConcurrency int
		numberOfChecks         int
		expectedResult         int
	}{
		{10, 2, 5, 4},
		{10, 5, 5, 10},
		{5, 10, 10, 10},
		{0, 0, 0, 0},
	}

	for _, test := range tests {
		result := InternalLevelConcurrency(test.maxConcurrency, test.minInternalConcurrency, test.numberOfChecks)
		if result != test.expectedResult {
			t.Errorf("InternalLevelConcurrency(%d, %d, %d) = %d; want %d",
				test.maxConcurrency,
				test.minInternalConcurrency,
				test.numberOfChecks,
				result,
				test.expectedResult)
		}
	}
}

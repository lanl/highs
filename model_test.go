// This file tests the high package's high-level model API.

package highs

import (
	"testing"
)

// TestMakeSparseMatrix tests the conversion of a slice of Nonzeros to start,
// index, and value slices.
func TestMakeSparseMatrix(t *testing.T) {
	// Construct a sparse matrix.
	var model Model
	model.CoeffMatrix = []Nonzero{
		{0, 1, 1.0},
		{1, 0, 1.0},
		{1, 1, 2.0},
		{2, 0, 3.0},
		{2, 1, 2.0},
	}
	start, index, value, err := nonzerosToCSR(model.CoeffMatrix, false)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the three slices.
	compSlices(t, "start", start, []int{0, 1, 3})
	compSlices(t, "index", index, []int{1, 0, 1, 0, 1})
	compSlices(t, "value", value, []float64{1.0, 1.0, 2.0, 3.0, 2.0})
}

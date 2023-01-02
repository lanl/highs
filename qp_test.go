// This file tests the high package's high-level API with quadratic-programming
// models.

package highs

import (
	"math"
	"testing"
)

// TestMinimalAPIQPMin mimics the first test in HiGHS's minimal_api_qp function
// from check/TestCAPI.c:
//
//	minimize -x_2 - 3x_3 + (1/2)(2x_1^2 - 2x_1x_3 + 0.2x_2^2 + 2x_3^2)
//
//	subject to x_1 + x_3 <= 2; x>=0
//
// Like TestCAPI.c, we don't actually enforce the x>=0 column constraints.
func TestMinimalAPIQPMin(t *testing.T) {
	// Prepare the model.
	var model Model
	model.ColCosts = []float64{0.0, -1.0, -3.0}
	model.AddDenseRow(-1e30, []float64{1.0, 0.0, 1.0}, 2.0)
	model.HessianMatrix = []Nonzero{
		{0, 0, 2.0},
		{0, 2, -1.0},
		{1, 1, 0.2},
		{2, 2, 2.0},
	}

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatalf("Solve failed (%s)", err)
	}
	if soln.Status != Optimal {
		t.Fatalf("Solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	primal := roundFloats(0.001, soln.ColumnPrimal)
	compSlices(t, "ColumnPrimal", primal, []float64{0.5, 5.0, 1.5})

	// Validate the objective value.
	if math.Round(soln.Objective/0.001)*0.001 != -5.25 {
		t.Fatalf("objective value was %.2f but should have been -5.25", soln.Objective)
	}
}

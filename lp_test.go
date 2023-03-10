// This file tests the high package's high-level API with linear-programming
// models.

package highs

import (
	"testing"
)

// TestMinimalAPIMin mimics the first test in HiGHS's minimal_api function from
// examples/call_highs_from_c.c:
//
//	Min    f  =  x_0 +  x_1 + 3
//	s.t.                x_1 <= 7
//	       5 <=  x_0 + 2x_1 <= 15
//	       6 <= 3x_0 + 2x_1
//	0 <= x_0 <= 4; 1 <= x_1
func TestMinimalAPIMin(t *testing.T) {
	// Prepare the model.
	var model Model
	model.Offset = 3.0
	model.ColCosts = []float64{1.0, 1.0}
	model.ColLower = []float64{0.0, 1.0}
	model.ColUpper = []float64{4.0, 1.0e30}
	model.RowLower = []float64{-1.0e30, 5.0, 6.0}
	model.RowUpper = []float64{7.0, 15.0, 1.0e30}
	model.ConstMatrix = []Nonzero{
		{0, 1, 1.0},
		{1, 0, 1.0},
		{1, 1, 2.0},
		{2, 0, 3.0},
		{2, 1, 2.0},
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
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{0.5, 2.25})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{2.25, 5.0, 6.0})
	compSlices(t, "ColumnDual", soln.ColumnDual, []float64{0.0, 0.0})
	compSlices(t, "RowDual", soln.RowDual, []float64{0.0, 0.25, 0.25})
	compSlices(t, "ColumnBasis", soln.ColumnBasis, []BasisStatus{Basic, Basic})
	compSlices(t, "RowBasis", soln.RowBasis, []BasisStatus{Basic, Lower, Lower})

	// Validate the objective value.
	if soln.Objective != 5.75 {
		t.Fatalf("objective value was %.2f but should have been 5.75", soln.Objective)
	}
}

// TestMinimalAPIMax mimics the second test in HiGHS's minimal_api function from
// examples/call_highs_from_c.c:
//
//	Max    f  =  x_0 +  x_1 + 3
//	s.t.                x_1 <= 7
//	       5 <=  x_0 + 2x_1 <= 15
//	       6 <= 3x_0 + 2x_1
//	0 <= x_0 <= 4; 1 <= x_1
func TestMinimalAPIMax(t *testing.T) {
	// Prepare the model.
	var model Model
	model.Maximize = true
	model.Offset = 3.0
	model.ColCosts = []float64{1.0, 1.0}
	model.ColLower = []float64{0.0, 1.0}
	model.ColUpper = []float64{4.0, 1.0e30}
	model.RowLower = []float64{-1.0e30, 5.0, 6.0}
	model.RowUpper = []float64{7.0, 15.0, 1.0e30}
	model.ConstMatrix = []Nonzero{
		{0, 1, 1.0},
		{1, 0, 1.0},
		{1, 1, 2.0},
		{2, 0, 3.0},
		{2, 1, 2.0},
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
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{4.0, 5.5})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{5.5, 15.0, 23.0})
	compSlices(t, "ColumnDual", soln.ColumnDual, []float64{0.5, 0.0})
	compSlices(t, "RowDual", soln.RowDual, []float64{0.0, 0.5, 0.0})
	compSlices(t, "ColumnBasis", soln.ColumnBasis, []BasisStatus{Upper, Basic})
	compSlices(t, "RowBasis", soln.RowBasis, []BasisStatus{Basic, Upper, Basic})

	// Validate the objective value.
	if soln.Objective != 12.5 {
		t.Fatalf("objective value was %.2f but should have been 12.5", soln.Objective)
	}
}

// TestAddDenseRow repeats the test in TestMinimalAPIMin but using the
// AddDenseRow convenience method.
func TestAddDenseRow(t *testing.T) {
	// Prepare the model.
	var model Model
	model.Maximize = false // Unnecessary but included for testing
	model.Offset = 3.0
	model.ColCosts = []float64{1.0, 1.0}
	model.ColLower = []float64{0.0, 1.0}
	model.ColUpper = []float64{4.0, 1.0e30}
	model.AddDenseRow(-1.0e30, []float64{0.0, 1.0}, 7.0)
	model.AddDenseRow(5.0, []float64{1.0, 2.0}, 15.0)
	model.AddDenseRow(6.0, []float64{3.0, 2.0}, 1.0e30)

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatalf("Solve failed (%s)", err)
	}
	if soln.Status != Optimal {
		t.Fatalf("Solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{0.5, 2.25})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{2.25, 5.0, 6.0})
	compSlices(t, "ColumnDual", soln.ColumnDual, []float64{0.0, 0.0})
	compSlices(t, "RowDual", soln.RowDual, []float64{0.0, 0.25, 0.25})
	compSlices(t, "ColumnBasis", soln.ColumnBasis, []BasisStatus{Basic, Basic})
	compSlices(t, "RowBasis", soln.RowBasis, []BasisStatus{Basic, Lower, Lower})

	// Validate the objective value.
	if soln.Objective != 5.75 {
		t.Fatalf("objective value was %.2f but should have been 5.75", soln.Objective)
	}
}

// TestImplicitColumnBounds tests that column bounds of (-inf, +inf) can be
// left unspecified.  It solves the following problem:
//
//	Satisfy 23 <= x_0 + x_1 <= 23
//	        17 <= x_0 - x_1 <= 17
func TestImplicitColumnBounds(t *testing.T) {
	// Prepare the model.
	var model Model
	model.AddDenseRow(23.0, []float64{1.0, 1.0}, 23.0)
	model.AddDenseRow(17.0, []float64{1.0, -1.0}, 17.0)

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatalf("Solve failed (%s)", err)
	}
	if soln.Status != Optimal {
		t.Fatalf("Solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{20.0, 3.0})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{23.0, 17.0})
}

// TestLPModelToRawModel sets up an LP model, converts it to a RawModel, and
// solves it.  It formulates the following test problem:
//
//	Satisfy 1 <= x_0 - x_1 <= 1
//	        5 <= x_0 + x_1 <= 5
func TestLPModelToRawModel(t *testing.T) {
	// Prepare the model.
	var model Model
	model.AddDenseRow(1.0, []float64{1.0, -1.0}, 1.0)
	model.AddDenseRow(5.0, []float64{1.0, 1.0}, 5.0)

	// Convert the Model to a RawModel.
	raw, err := model.ToRawModel()
	if err != nil {
		t.Fatal(err)
	}
	checkErr(t, raw.SetBoolOption("output_flag", false))

	// Solve the model.
	soln, err := raw.Solve()
	if err != nil {
		t.Fatal(err)
	}
	if soln.Status != Optimal {
		t.Fatalf("Solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{3.0, 2.0})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{1.0, 5.0})
}

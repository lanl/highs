// This file tests the high package's mixed-integer programming wrappers.

package highs

import "testing"

// TestMinimalAPIMaxMIP mimics the third test in HiGHS's minimal_api function
// from examples/call_highs_from_c.c:
//
//	Max    f  =  x_0 +  x_1 + 3
//	s.t.                x_1 <= 7
//	       5 <=  x_0 + 2x_1 <= 15
//	       6 <= 3x_0 + 2x_1
//	0 <= x_0 <= 4; 1 <= x_1
func TestMinimalAPIMaxMIP(t *testing.T) {
	// Prepare the model.
	var model MIPModel
	model.Maximize = true
	model.Offset = 3.0
	model.ColCosts = []float64{1.0, 1.0}
	model.ColLower = []float64{0.0, 1.0}
	model.ColUpper = []float64{4.0, 1.0e30}
	model.RowLower = []float64{-1.0e30, 5.0, 6.0}
	model.RowUpper = []float64{7.0, 15.0, 1.0e30}
	model.CoeffMatrix = []Nonzero{
		{0, 1, 1.0},
		{1, 0, 1.0},
		{1, 1, 2.0},
		{2, 0, 3.0},
		{2, 1, 2.0},
	}
	model.VarTypes = []VariableType{IntegerType, IntegerType}

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatalf("solve failed (%s)", err)
	}
	if soln.Status != Optimal {
		t.Fatalf("solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{4.0, 5.0})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{5, 14.0, 22.0})

	// Validate the objective value.
	if int(soln.Objective) != 12 {
		t.Fatalf("objective value was %d but should have been 12", int(soln.Objective))
	}
}

// TestMIPModelToRawModel sets up an MIPModel, converts it to a RawModel, and
// solves it.  We use the following test problem:
//
//	Satisfy 1 <= x_0 - x_1 <= 1
//	        5 <= x_0 + x_1 <= 5
func TestMIPModelToRawModel(t *testing.T) {
	// Prepare the model.
	var model MIPModel
	model.AddDenseRow(1.0, []float64{1.0, -1.0}, 1.0)
	model.AddDenseRow(5.0, []float64{1.0, 1.0}, 5.0)
	model.VarTypes = []VariableType{IntegerType, IntegerType}

	// Convert the MIPModel to a RawModel.
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
		t.Fatalf("solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{3.0, 2.0})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{1.0, 5.0})
}

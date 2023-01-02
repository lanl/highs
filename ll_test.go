// This file tests the high package's low-level API wrappers.

package highs

import (
	"math"
	"os"
	"testing"
)

// TestFullAPIMin mimics the first test in HiGHS's full_api function from
// examples/call_highs_from_c.c:
//
//	Min    f  =  x_0 +  x_1 + 3
//	s.t.                x_1 <= 7
//	       5 <=  x_0 + 2x_1 <= 15
//	       6 <= 3x_0 + 2x_1
//	0 <= x_0 <= 4; 1 <= x_1
func TestFullAPIMin(t *testing.T) {
	// Prepare the model.
	model := NewRawModel()
	checkErr(t, model.SetBoolOption("output_flag", false))
	checkErr(t, model.SetMaximization(false)) // Unnecessary but included for testing
	checkErr(t, model.SetOffset(3.0))
	checkErr(t, model.AddColumnBounds([]float64{0.0, 1.0},
		[]float64{4.0, 1.0e30}))
	checkErr(t, model.SetColumnCosts([]float64{1.0, 1.0}))
	checkErr(t, model.AddCompSparseRows([]float64{-1.0e30, 5.0, 6.0},
		[]int{0, 1, 3}, []int{1, 0, 1, 0, 1}, []float64{1.0, 1.0, 2.0, 3.0, 2.0},
		[]float64{7.0, 15.0, 1.0e30}))

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatal(err)
	}
	if soln.Status != Optimal {
		t.Fatalf("solve returned %s instead of Optimal", soln.Status)
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

// TestFullAPIInfeasible verifies that infeasible models are handled properly.
// It defines the following model:
//
//	Satisfy 4 <= x_0 <= 4
//	        5 <= x_0 <= 5
//	subject to 0 <= x_0 <= 10
func TestFullAPIInfeasible(t *testing.T) {
	// Prepare the model.
	model := NewRawModel()
	checkErr(t, model.SetBoolOption("output_flag", false))
	checkErr(t, model.AddColumnBounds([]float64{0.0, 0.0}, []float64{10.0, 10.0}))
	checkErr(t, model.AddDenseRow(4.0, []float64{1.0}, 4.0))
	checkErr(t, model.AddDenseRow(5.0, []float64{1.0}, 5.0))

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatal(err)
	}
	if soln.Status != Infeasible {
		t.Fatalf("solve returned %s instead of Infeasible", soln.Status)
	}
}

// TestFullAPIMaxMIP mimics with the "full" API the third test in HiGHS's
// minimal_api function from examples/call_highs_from_c.c:
//
//	Max    f  =  x_0 +  x_1 + 3
//	s.t.                x_1 <= 7
//	       5 <=  x_0 + 2x_1 <= 15
//	       6 <= 3x_0 + 2x_1
//	0 <= x_0 <= 4; 1 <= x_1
func TestFullAPIMaxMIP(t *testing.T) {
	// Prepare the model.
	model := NewRawModel()
	checkErr(t, model.SetBoolOption("output_flag", false))
	checkErr(t, model.SetOffset(3.0))
	checkErr(t, model.AddColumnBounds([]float64{0.0, 1.0},
		[]float64{4.0, 1.0e30}))
	checkErr(t, model.SetColumnCosts([]float64{1.0, 1.0}))
	checkErr(t, model.AddCompSparseRows([]float64{-1.0e30, 5.0, 6.0},
		[]int{0, 1, 3}, []int{1, 0, 1, 0, 1}, []float64{1.0, 1.0, 2.0, 3.0, 2.0},
		[]float64{7.0, 15.0, 1.0e30}))
	checkErr(t, model.SetIntegrality([]VariableType{IntegerType, IntegerType}))
	checkErr(t, model.SetMaximization(true))

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatal(err)
	}
	if soln.Status != Optimal {
		t.Fatalf("solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{4.0, 5.0})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{5.0, 14.0, 22.0})

	// Validate the objective value.
	if int(soln.Objective) != 12 {
		t.Fatalf("objective value was %d but should have been 12", int(soln.Objective))
	}
}

var mpsFile *os.File // MPS file to write and read

// TestWriteModel creates a model and writes it to a throwaway file.  The model
// represented is as follows:
//
//	Min. 2*x_0 + x_1
//	s.t. 10 <= x_0 + x_1 <= 10
//	      4 <= x_0 - x_1 <=  4
//	1 <= x_1 <= 25, 1 <= x_2 <= 25
func TestWriteModel(t *testing.T) {
	// Prepare the model.
	model := NewRawModel()
	checkErr(t, model.SetBoolOption("output_flag", false))
	checkErr(t, model.AddColumnBounds([]float64{1.0, 1.0},
		[]float64{25.0, 25.0}))
	checkErr(t, model.SetColumnCosts([]float64{2.0, 1.0}))
	checkErr(t, model.AddDenseRow(10.0, []float64{1.0, 1.0}, 10.0))
	checkErr(t, model.AddDenseRow(4.0, []float64{1.0, -1.0}, 4.0))
	checkErr(t, model.SetIntegrality([]VariableType{IntegerType, IntegerType}))

	// Write the model to a temporary file.  Remember the name of the
	// temporary file for use in TestReadModel.
	var err error
	mpsFile, err = os.CreateTemp("", "highs-*.mps")
	if err != nil {
		t.Fatal(err)
	}
	defer mpsFile.Close()
	t.Logf("Writing MPS to %s", mpsFile.Name())
	checkErr(t, model.WriteModelToFile(mpsFile.Name()))
}

// TestReadModel reads the model previously written by TestWriteModel and
// solves it.  If successful, it deletes the throwaway model file.
func TestReadModel(t *testing.T) {
	// Read the model.
	fname := mpsFile.Name()
	t.Logf("Reading MPS from %s", fname)
	model := NewRawModel()
	checkErr(t, model.SetBoolOption("output_flag", false))
	err := model.ReadModelFromFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatal(err)
	}
	if soln.Status != Optimal {
		t.Fatalf("solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	compSlices(t, "ColumnPrimal", soln.ColumnPrimal, []float64{7.0, 3.0})
	compSlices(t, "RowPrimal", soln.RowPrimal, []float64{10.0, 4.0})

	// Validate the objective value.
	if int(soln.Objective) != 17 {
		t.Fatalf("objective value was %d but should have been 17", int(soln.Objective))
	}

	// Remove the throwaway MPS file.
	err = os.Remove(fname)
	if err != nil {
		t.Fatal(err)
	}
}

// TestFullAPIQPMin performs the same test as TestMinimalAPIQPMin but
// using the low-level API:
//
//	minimize -x_2 - 3x_3 + (1/2)(2x_1^2 - 2x_1x_3 + 0.2x_2^2 + 2x_3^2)
//
//	subject to x_1 + x_3 <= 2; x>=0
func TestFullAPIQPMin(t *testing.T) {
	// Prepare the model.
	model := NewRawModel()
	checkErr(t, model.SetBoolOption("output_flag", false))
	pInf := math.Inf(1)
	checkErr(t, model.AddColumnBounds([]float64{0.0, 0.0, 0.0},
		[]float64{pInf, pInf, pInf}))
	checkErr(t, model.SetColumnCosts([]float64{0.0, -1.0, -3.0}))
	checkErr(t, model.AddCompSparseRows([]float64{math.Inf(-1)},
		[]int{0}, []int{0, 2}, []float64{1.0, 1.0},
		[]float64{2.0}))
	checkErr(t, model.AddCompSparseHessian([]int{0, 2, 3},
		[]int{0, 2, 1, 2},
		[]float64{2.0, -1.0, 0.2, 2.0}))

	// Solve the model.
	soln, err := model.Solve()
	if err != nil {
		t.Fatalf("solve failed (%s)", err)
	}
	if soln.Status != Optimal {
		t.Fatalf("solve returned %s instead of Optimal", soln.Status)
	}

	// Confirm that each field is as expected.
	primal := roundFloats(0.001, soln.ColumnPrimal)
	compSlices(t, "ColumnPrimal", primal, []float64{0.5, 5.0, 1.5})

	// Validate the objective value.
	if math.Round(soln.Objective/0.001)*0.001 != -5.25 {
		t.Fatalf("objective value was %.2f but should have been -5.25", soln.Objective)
	}
}

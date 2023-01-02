// This file tests the high package's high-level model API.  Although solvers
// are invoked, the focus of this file is on features other than solving.
// Other test files stress the solvers.

package highs

import (
	"bytes"
	"os"
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

var mpsFile *os.File // MPS file to write and read

// TestWriteModelToFile creates a model and writes it to a throwaway file.  The
// model represented is as follows:
//
//	Min. 2*x_0 + x_1
//	s.t. 10 <= x_0 + x_1 <= 10
//	      4 <= x_0 - x_1 <=  4
//	1 <= x_1 <= 25, 1 <= x_2 <= 25
func TestWriteModelToFile(t *testing.T) {
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

// TestReadModelFromFile reads the model previously written by TestWriteModel
// and solves it.  If successful, it deletes the throwaway model file.
func TestReadModelFromFile(t *testing.T) {
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
		t.Fatalf("Solve returned %s instead of Optimal", soln.Status)
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

// TestReadWriteModel tests writing a model to a buffer then reading it back in
// and solving it.  It uses the same model as
// TestWriteModelToFile/TestReadModelFromFile.
func TestReadWriteModel(t *testing.T) {
	// Prepare the model.
	m1 := NewRawModel()
	checkErr(t, m1.SetBoolOption("output_flag", false))
	checkErr(t, m1.AddColumnBounds([]float64{1.0, 1.0},
		[]float64{25.0, 25.0}))
	checkErr(t, m1.SetColumnCosts([]float64{2.0, 1.0}))
	checkErr(t, m1.AddDenseRow(10.0, []float64{1.0, 1.0}, 10.0))
	checkErr(t, m1.AddDenseRow(4.0, []float64{1.0, -1.0}, 4.0))
	checkErr(t, m1.SetIntegrality([]VariableType{IntegerType, IntegerType}))

	// Write the model to a buffer.
	var buf bytes.Buffer
	checkErr(t, m1.WriteModel(&buf))

	// Read from the buffer into a second model.
	m2 := NewRawModel()
	checkErr(t, m2.SetBoolOption("output_flag", false))
	checkErr(t, m2.ReadModel(&buf))
}

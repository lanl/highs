// This file tests the high package's RawSolution wrapper.

package highs

import (
	"bytes"
	"testing"
)

// modelAndSolve is a helper function that constructs and solves a simple MIP
// model.  The model is as follows:
//
//	Min.  3*x_0 + 2*x_1 + 1*x_2
//	s.t.  1 <= x_0 - x_1
//	      1 <= x_1 - x_2
//	      10 <= x_0 + x_1 + x_2
//	with  0 <= x_0, x_1, x_2
func modelAndSolve() (*RawSolution, error) {
	// Prepare the model.
	var model Model
	model.ColCosts = []float64{3.0, 2.0, 1.0}
	model.ColLower = []float64{0.0, 0.0, 0.0}
	model.RowLower = []float64{1.0, 1.0, 10.0}
	model.CoeffMatrix = []Nonzero{
		{0, 0, 1.0},
		{0, 1, -1.0},
		{1, 1, 1.0},
		{1, 2, -1.0},
		{2, 0, 1.0},
		{2, 1, 1.0},
		{2, 2, 1.0},
	}
	model.VarTypes = []VariableType{IntegerType, IntegerType, IntegerType}

	// Convert the Model to a RawModel and solve it.
	raw, err := model.ToRawModel()
	if err != nil {
		return nil, err
	}
	err = raw.SetBoolOption("output_flag", false)
	if err != nil {
		return nil, err
	}
	return raw.Solve() // Solution and error code
}

// TestWriteSolution tests the writing of a solution in a textual format.
func TestWriteSolution(t *testing.T) {
	// Produce a solution.
	soln, err := modelAndSolve()
	if err != nil {
		t.Fatal(err)
	}

	// Write the solution to a buffer.
	var buf bytes.Buffer
	checkErr(t, soln.WriteSolution(&buf, false))

	// Compare to the expected contents.
	exp := `Model status
Optimal

# Primal solution values
Feasible
Objective 23
# Columns 3
C0 5
C1 3
C2 2
# Rows 3
R0 2
R1 1
R2 10

# Dual solution values
None

# Basis
HiGHS v1
None
`
	if buf.String() != exp {
		t.Logf("Expected: %q", exp)
		t.Logf("Actual:   %q", buf.String())
		t.Fatal("textual solution was not as expected")
	}
}

// TestWriteSolutionPretty tests the writing of a solution in a human-friendly
// textual format.
func TestWriteSolutionPretty(t *testing.T) {
	// Produce a solution.
	soln, err := modelAndSolve()
	if err != nil {
		t.Fatal(err)
	}

	// Write the solution to a buffer.
	var buf bytes.Buffer
	checkErr(t, soln.WriteSolution(&buf, true))

	// Compare to the expected contents.
	exp := `Columns
    Index Status        Lower        Upper       Primal         Dual  Type      
        0                   0          inf            5               Integer   
        1                   0          inf            3               Integer   
        2                   0          inf            2               Integer   
Rows
    Index Status        Lower        Upper       Primal         Dual
        0                   1          inf            2             
        1                   1          inf            1             
        2                  10          inf           10             

Model status: Optimal

Objective value: 23
`
	if buf.String() != exp {
		t.Logf("Expected: %q", exp)
		t.Logf("Actual:   %q", buf.String())
		t.Fatal("textual solution was not as expected")
	}
}

// This file tests the high package's RawSolution wrapper.

package highs

import (
	"bytes"
	"errors"
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

// TestGetIntInfo tests that GetIntInfo works.
func TestGetIntInfo(t *testing.T) {
	// Produce a solution.
	soln, err := modelAndSolve()
	if err != nil {
		t.Fatal(err)
	}

	// Query it for information.
	dss, err := soln.GetIntInfo("dual_solution_status")
	if err != nil {
		t.Fatal(err)
	}
	if dss != 0 {
		t.Fatalf("expected a dual solution status of 0 but saw %d", dss)
	}
}

// TestGetInt64Info tests that GetInt64Info works.
func TestGetInt64Info(t *testing.T) {
	// Produce a solution.
	soln, err := modelAndSolve()
	if err != nil {
		t.Fatal(err)
	}

	// Query it for information.
	mnc, err := soln.GetInt64Info("mip_node_count")
	if err != nil {
		t.Fatal(err)
	}
	if mnc != 1 {
		t.Fatalf("expected a MIP node count of 1 but saw %d", mnc)
	}
}

// TestGetInt64InfoBad tests that GetInt64Info returns an error for a
// nonexistent key.
func TestGetInt64InfoBad(t *testing.T) {
	// Produce a solution.
	soln, err := modelAndSolve()
	if err != nil {
		t.Fatal(err)
	}

	// Query it for nonexistent information.
	_, err = soln.GetInt64Info("bogus info key")
	var cs CallStatus
	if !errors.As(err, &cs) {
		t.Fatalf("expected a failure code but received %v", err)
	}
}

// TestGetFloat64Info tests that GetFloat64Info works.
func TestGetFloat64Info(t *testing.T) {
	// Produce a solution.
	soln, err := modelAndSolve()
	if err != nil {
		t.Fatal(err)
	}

	// Query it for information.
	miv, err := soln.GetFloat64Info("max_integrality_violation")
	if err != nil {
		t.Fatal(err)
	}
	if miv != 0.0 {
		t.Fatalf("expected a maximum integrality violation of 0 but saw %v", miv)
	}
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

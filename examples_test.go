// This file provides a few highs usage examples.

package highs_test

import (
	"fmt"
	"math"

	"github.com/lanl/highs"
)

// This example demonstrates how to convert the sparse matrix
// [ 1.1 0.0 0.0 0.0 2.2 ; 0.0 3.3 0.0 4.4 0.0 ; 0.0 0.0 5.5 0.0 0.0 ]
// (MATLAB syntax) from a list of nonzero elements expressed as {row, column,
// value} tuples to compressed sparse row format.
func ExampleNonzerosToCSR() {
	// Define a matrix in terms of its nonzero elements.
	nz := []highs.Nonzero{
		{0, 0, 1.1},
		{0, 4, 2.2},
		{1, 1, 3.3},
		{1, 3, 4.4},
		{2, 2, 5.5},
	}
	fmt.Println("nz:", nz)

	// Convert nonzeros to compressed sparse row form.
	start, index, value, err := highs.NonzerosToCSR(nz, false)
	if err != nil {
		panic("failed to convert nonzeros to CSR")
	}
	fmt.Println("start:", start)
	fmt.Println("index:", index)
	fmt.Println("value:", value)
	// Output:
	// nz: [{0 0 1.1} {0 4 2.2} {1 1 3.3} {1 3 4.4} {2 2 5.5}]
	// start: [0 2 4]
	// index: [0 4 1 3 2]
	// value: [1.1 2.2 3.3 4.4 5.5]
}

// The following code adds three rows to a model, each with a lower and upper
// row bound.
func ExampleModel_AddDenseRow() {
	var m highs.Model
	m.AddDenseRow(1.0, []float64{1.0, -1.0, 0.0, 0.0}, 2.0) // 1.0 ≤ A − B ≤ 2.0
	m.AddDenseRow(2.0, []float64{0.0, 1.0, -1.0, 0.0}, 4.0) // 2.0 ≤ B − C ≤ 4.0
	m.AddDenseRow(3.0, []float64{0.0, 0.0, 1.0, -1.0}, 8.0) // 3.0 ≤ C − D ≤ 8.0
	fmt.Println("RowLower:", m.RowLower)
	fmt.Println("RowUpper:", m.RowUpper)
	fmt.Println("ConstMatrix:", m.ConstMatrix)
	// Output:
	// RowLower: [1 2 3]
	// RowUpper: [2 4 8]
	// ConstMatrix: [{0 0 1} {0 1 -1} {1 1 1} {1 2 -1} {2 2 1} {2 3 -1}]
}

// Low-level models default to writing verbose status messages.  SetBoolOption
// can be used to disable these.
func ExampleRawModel_SetBoolOption() {
	m := highs.NewRawModel()
	m.SetBoolOption("output_flag", false)
}

// Here is a complete example of using the highs package's high-level interface
// to set up a model, solve it, and report the solution.  The problem we choose
// to solve is as follows: What is the maximum total face value of three
// six-sided dice A, B, and C such that the difference in face value between A
// and B is exactly twice the difference in face value between B and C, where B
// is strictly greater than C?
func ExampleModel_Solve() {
	// Prepare a mixed-integer programming model.
	var m highs.Model
	m.VarTypes = []highs.VariableType{
		highs.IntegerType, // A is an integer.
		highs.IntegerType, // B is an integer.
		highs.IntegerType, // C is an integer.
	}
	m.ColCosts = []float64{1.0, 1.0, 1.0}                      // Objective function is A + B + C.
	m.Maximize = true                                          // Maximize the objective function.
	m.ColLower = []float64{1.0, 1.0, 1.0}                      // A ≥ 1, B ≥ 1, C ≥ 1.
	m.ColUpper = []float64{6.0, 6.0, 6.0}                      // A ≤ 6, B ≤ 6, C ≤ 6.
	m.AddDenseRow(0.0, []float64{1.0, -3.0, 2.0}, 0.0)         // A − B = 2(B − C), expressed as 0 ≤ A − 3B + 2C ≤ 0.
	m.AddDenseRow(1.0, []float64{0.0, 1.0, -1.0}, math.Inf(1)) // B > C, expressed as 1 ≤ B − C ≤ ∞.

	// Find an optimal solution.
	soln, err := m.Solve()
	if err != nil {
		panic(err)
	}

	// Output the A, B, and C that maximize the objective function and the objective value itself.
	fmt.Println("A:", soln.ColumnPrimal[0])
	fmt.Println("B:", soln.ColumnPrimal[1])
	fmt.Println("C:", soln.ColumnPrimal[2])
	fmt.Println("Total face value:", soln.Objective)
	// Output:
	// A: 6
	// B: 4
	// C: 3
	// Total face value: 13
}

// Here is a complete example of using the highs package's low-level interface
// to set up a model, solve it, and report the solution.  The problem we choose
// to solve is as follows: What is the maximum total face value of three
// six-sided dice A, B, and C such that the difference in face value between A
// and B is exactly twice the difference in face value between B and C, where B
// is strictly greater than C?
//
// Remove the SetBoolOption line to view HiGHS status output.
func ExampleRawModel_Solve() {
	// Define a function that panics on error.
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// Prepare a mixed-integer programming model.  AddColumnBounds adds
	// variables (columns) so it must be called before any of the functions
	// that set column properties.
	m := highs.NewRawModel()
	checkErr(m.SetBoolOption("output_flag", false))
	checkErr(m.AddColumnBounds(
		[]float64{1.0, 1.0, 1.0},  // A ≥ 1, B ≥ 1, C ≥ 1.
		[]float64{6.0, 6.0, 6.0})) // A ≤ 6, B ≤ 6, C ≤ 6.
	checkErr(m.SetIntegrality([]highs.VariableType{
		highs.IntegerType, // A is an integer.
		highs.IntegerType, // B is an integer.
		highs.IntegerType, // C is an integer.
	}))
	checkErr(m.SetColumnCosts([]float64{1.0, 1.0, 1.0})) // Objective function is A + B + C.
	checkErr(m.SetMaximization(true))                    // Maximize the objective function.

	// Represent A − B = 2(B − C) and B > C as the sparse matrix inequality,
	//     ⌈0⌉   ⌈ 1 −3  2 ⌉ ⌈A⌉   ⌈0⌉
	//     | | ≤ |         | |B| ≤ | |
	//     ⌊1⌋   ⌊ 0  1 −1 ⌋ ⌊C⌋   ⌊∞⌋
	checkErr(m.AddCompSparseRows(
		[]float64{0.0, 1.0},
		[]int{0, 3},
		[]int{0, 1, 2, 1, 2},
		[]float64{1.0, -3.0, 2.0, 1.0, -1.0},
		[]float64{0.0, math.Inf(1)}))

	// Find an optimal solution.
	soln, err := m.Solve()
	if err != nil {
		panic(err)
	}

	// Output the A, B, and C that maximize the objective function and the objective value itself.
	fmt.Println("A:", soln.ColumnPrimal[0])
	fmt.Println("B:", soln.ColumnPrimal[1])
	fmt.Println("C:", soln.ColumnPrimal[2])
	fmt.Println("Total face value:", soln.Objective)
	// Output:
	// A: 6
	// B: 4
	// C: 3
	// Total face value: 13
}

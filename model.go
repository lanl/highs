// This file provides fields and methods that are common across multiple HiGHS
// models.

package highs

// A NonZero represents a nonzero entry in a sparse matrix.  Rows and columns
// are indexed from zero.
type NonZero struct {
	Row   int
	Col   int
	Value float64
}

// An commonModel represents fields common to many HiGHS models.
type commonModel struct {
	minimize    bool      // true=minimize; false=maximize
	colCost     []float64 // Column costs (objective function)
	offset      float64   // Objective-function offset
	colLower    []float64 // Column lower bounds
	colUpper    []float64 // Column upper bounds
	rowLower    []float64 // Row lower bounds
	rowUpper    []float64 // Row upper bounds
	coeffMatrix []NonZero // Sparse "A" matrix
}

// Set a model's lower and upper row bounds.
func (m *commonModel) SetRowBounds(lb, ub []float64) {
	m.rowLower = lb
	m.rowUpper = ub
}

// Set a model's lower and upper column bounds.
func (m *commonModel) SetColumnBounds(lb, ub []float64) {
	m.colLower = lb
	m.colUpper = ub
}

// Set a model's coefficient matrix.
func (m *commonModel) SetCoefficients(nz []NonZero) {
	m.coeffMatrix = nz
}

// highs provides a Go interface to the HiGHS optimizer.
package highs

// A NonZero represents a nonzero entry in a sparse matrix.
type NonZero struct {
	Row   int
	Col   int
	Value float64
}

// An LPModel represents a HiGHS linear-programming model.
type LPModel struct {
	minimize    bool      // true=minimize; false=maximize
	colCost     []float64 // Column costs (objective function)
	offset      float     // Objective-function offset
	colLower    []float64 // Column lower bounds
	colUpper    []float64 // Column upper bounds
	rowLower    []float64 // Row lower bounds
	rowUpper    []float64 // Row upper bounds
	coeffMatrix []NonZero // Sparse "A" matrix
}

// NewLPModel allocates and returns an empty linear-programming model.
func NewLPModel() *LPModel {
	return &LPModel{}
}

// Set a model's lower and upper row bounds.
func (m *LPModel) SetRowBounds(lb, ub []float64) {
	m.rowLower = lb
	m.rowUpper = ub
}

// Set a model's lower and upper column bounds.
func (m *LPModel) SetColumnBounds(lb, ub []float64) {
	m.colLower = lb
	m.colUpper = ub
}

// Set a model's coefficient matrix.
func (m *LPModel) SetCoefficients(nz []NonZero) {
	m.coeffMatrix = nz
}

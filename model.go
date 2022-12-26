// This file provides fields and methods that are common across multiple HiGHS
// models.

package highs

import "math"

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

// prepareBounds replaces nil column or row bounds with infinities.
func (m *commonModel) prepareBounds(lb, ub []float64) ([]float64, []float64) {
	switch {
	case lb == nil && ub == nil:
		// No bounds were provided.
	case lb == nil:
		// Replace nil lower bounds with minus infinity.
		mInf := math.Inf(-1)
		lb = make([]float64, len(ub))
		for i := range lb {
			lb[i] = mInf
		}
	case ub == nil:
		// Replace nil upper bounds with plus infinity.
		pInf := math.Inf(1)
		ub = make([]float64, len(lb))
		for i := range ub {
			ub[i] = pInf
		}
	case len(lb) != len(ub):
		panic("different numbers of lower and upper bounds were provided")
	}
	return lb, ub
}

// SetRowBounds sets a model's lower and upper row bounds.
func (m *commonModel) SetRowBounds(lb, ub []float64) {
	m.rowLower, m.rowUpper = m.prepareBounds(lb, ub)
}

// SetColumnBounds sets a model's lower and upper column bounds.
func (m *commonModel) SetColumnBounds(lb, ub []float64) {
	m.colLower, m.colUpper = m.prepareBounds(lb, ub)
}

// SetCoefficients sets a model's coefficient matrix.
func (m *commonModel) SetCoefficients(nz []NonZero) {
	m.coeffMatrix = nz
}

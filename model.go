// This file provides fields and methods that are common across multiple HiGHS
// models.

package highs

import (
	"fmt"
	"math"
	"sort"
)

// #include <interfaces/highs_c_api.h>
import "C"

// A NonZero represents a nonzero entry in a sparse matrix.  Rows and columns
// are indexed from zero.
type NonZero struct {
	Row   int
	Col   int
	Value float64
}

// A commonModel represents fields common to many HiGHS models.
type commonModel struct {
	maximize    bool      // true=maximize; false=minimize
	colCosts    []float64 // Column costs (objective function)
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

// SetRowBounds specifies a model's lower and upper row bounds.  If the
// lower-bound argument is nil it is replaced with a slice of negative
// infinities.  If the upper-bound argument is nil, it is replaced with a slice
// of positive infinities.
func (m *commonModel) SetRowBounds(lb, ub []float64) {
	m.rowLower, m.rowUpper = m.prepareBounds(lb, ub)
}

// SetColumnBounds specifies a model's lower and upper column bounds.  If the
// lower-bound argument is nil it is replaced with a slice of negative
// infinities.  If the upper-bound argument is nil, it is replaced with a slice
// of positive infinities.
func (m *commonModel) SetColumnBounds(lb, ub []float64) {
	m.colLower, m.colUpper = m.prepareBounds(lb, ub)
}

// SetCoefficients specifies a model's coefficient matrix in terms of its
// nonzero entries.
func (m *commonModel) SetCoefficients(nz []NonZero) {
	// Complain about negative indices.
	for _, v := range nz {
		if v.Row < 0 || v.Col < 0 {
			panic(fmt.Sprintf("(%d, %d) is not a valid coordinate for a matrix coefficient",
				v.Row, v.Col))
		}
	}

	// Make a copy of the nonzeroes and sort the copy in place.
	sorted := make([]NonZero, len(nz))
	copy(sorted, nz)
	sort.SliceStable(sorted, func(i, j int) bool {
		nz0 := sorted[i]
		nz1 := sorted[j]
		switch {
		case nz0.Row < nz1.Row:
			return true
		case nz0.Row > nz1.Row:
			return false
		case nz0.Col < nz1.Col:
			return true
		case nz0.Col > nz1.Col:
			return false
		default:
			return false // Equal coordinates
		}
	})

	// Elide duplicate entries, keeping the latest value.
	m.coeffMatrix = make([]NonZero, 0, len(sorted))
	for _, v := range sorted {
		i := len(m.coeffMatrix)
		switch {
		case i == 0:
			// First element: always include.
			m.coeffMatrix = append(m.coeffMatrix, v)
		case v.Row == m.coeffMatrix[i-1].Row && v.Col == m.coeffMatrix[i-1].Col:
			// Duplicate coordinate: retain the later value.
			m.coeffMatrix[i-1].Value = v.Value
		default:
			// New coordinate.
			m.coeffMatrix = append(m.coeffMatrix, v)
		}
	}
}

// SetMaximization tells a model to maximize (true) or minimize (false) its
// objective function.
func (m *commonModel) SetMaximization(max bool) {
	m.maximize = max
}

// SetColumnCosts specifies a model's column costs (i.e., its objective
// function).
func (m *commonModel) SetColumnCosts(cs []float64) {
	m.colCosts = cs
}

// SetOffset specifies a constant offset for the objective function.
func (m *commonModel) SetOffset(o float64) {
	m.offset = o
}

// makeSparseMatrix converts coeffMatrix to a row-wise sparse-matrix
// representation in the form of a set of C vectors accepted by the HiGHS APIs.
func (m *commonModel) makeSparseMatrix() (start, index []C.HighsInt, value []C.double) {
	// Allocate memory for all of our return vectors.
	start = make([]C.HighsInt, 0, len(m.coeffMatrix))
	index = make([]C.HighsInt, 0, len(m.coeffMatrix))
	value = make([]C.double, 0, len(m.coeffMatrix))

	// Construct slices of C types.
	prevRow := -1
	for _, nz := range m.coeffMatrix {
		if nz.Row > prevRow {
			start = append(start, C.HighsInt(len(value)))
			index = append(index, C.HighsInt(nz.Col))
			prevRow = nz.Row
		}
		value = append(value, C.double(nz.Value))
	}
	return start, index, value
}

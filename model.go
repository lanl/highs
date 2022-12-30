// This file provides fields and methods that are common across multiple HiGHS
// models.

package highs

import (
	"fmt"
	"math"
	"sort"
)

// #include "highs-externs.h"
import "C"

// A Nonzero represents a nonzero entry in a sparse matrix.  Rows and columns
// are indexed from zero.
type Nonzero struct {
	Row   int
	Col   int
	Value float64
}

// A BasisStatus represents the basis status of a row or column.
type BasisStatus int

// These are the values a BasisStatus accepts.
const (
	UnknownBasisStatus BasisStatus = iota
	Lower
	Basic
	Upper
	Zero
	NonBasic
)

// convertHighsBasisStatus converts a kHighsBasisStatus to a BasisStatus.
func convertHighsBasisStatus(hbs C.HighsInt) BasisStatus {
	switch hbs {
	case C.kHighsBasisStatusLower:
		return Lower
	case C.kHighsBasisStatusBasic:
		return Basic
	case C.kHighsBasisStatusUpper:
		return Upper
	case C.kHighsBasisStatusZero:
		return Zero
	case C.kHighsBasisStatusNonbasic:
		return NonBasic
	default:
		return UnknownBasisStatus
	}
}

// A ModelStatus represents the status of an attempt to solve a model.
type ModelStatus int

// These are the values a ModelStatus accepts.
const (
	UnknownModelStatus ModelStatus = iota
	NotSet
	LoadError
	ModelError
	PresolveError
	SolveError
	PostsolveError
	ModelEmpty
	Optimal
	Infeasible
	UnboundedOrInfeasible
	Unbounded
	ObjectiveBound
	ObjectiveTarget
	TimeLimit
	IterationLimit
)

// convertHighsModelStatus converts a kHighsModelStatus to a ModelStatus.
func convertHighsModelStatus(hms C.HighsInt) ModelStatus {
	switch hms {
	case C.kHighsModelStatusNotset:
		return NotSet
	case C.kHighsModelStatusLoadError:
		return LoadError
	case C.kHighsModelStatusModelError:
		return ModelError
	case C.kHighsModelStatusPresolveError:
		return PresolveError
	case C.kHighsModelStatusSolveError:
		return SolveError
	case C.kHighsModelStatusPostsolveError:
		return PostsolveError
	case C.kHighsModelStatusModelEmpty:
		return ModelEmpty
	case C.kHighsModelStatusOptimal:
		return Optimal
	case C.kHighsModelStatusInfeasible:
		return Infeasible
	case C.kHighsModelStatusUnboundedOrInfeasible:
		return UnboundedOrInfeasible
	case C.kHighsModelStatusUnbounded:
		return Unbounded
	case C.kHighsModelStatusObjectiveBound:
		return ObjectiveBound
	case C.kHighsModelStatusObjectiveTarget:
		return ObjectiveTarget
	case C.kHighsModelStatusTimeLimit:
		return TimeLimit
	case C.kHighsModelStatusIterationLimit:
		return IterationLimit
	default:
		return UnknownModelStatus
	}
}

// convertHighsStatusToError converts a kHighsStatus to a Go error.
func convertHighsStatusToError(st C.HighsInt, caller string) error {
	switch st {
	case C.kHighsStatusOk:
		return nil
	case C.kHighsStatusError:
		return fmt.Errorf("%s failed with an error", caller)
	case C.kHighsStatusWarning:
		return fmt.Errorf("%s failed with a warning", caller)
	default:
		return fmt.Errorf("%s failed with unknown status", caller)
	}
}

// A commonModel represents fields common to many HiGHS models.
type commonModel struct {
	maximize    bool      // true=maximize; false=minimize
	colCosts    []float64 // Column costs (i.e., the objective function)
	offset      float64   // Objective-function offset
	colLower    []float64 // Column lower bounds
	colUpper    []float64 // Column upper bounds
	rowLower    []float64 // Row lower bounds
	rowUpper    []float64 // Row upper bounds
	coeffMatrix []Nonzero // Sparse "A" matrix
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
func (m *commonModel) SetCoefficients(nz []Nonzero) {
	// Complain about negative indices.
	for _, v := range nz {
		if v.Row < 0 || v.Col < 0 {
			panic(fmt.Sprintf("(%d, %d) is not a valid coordinate for a matrix coefficient",
				v.Row, v.Col))
		}
	}

	// Make a copy of the nonzeroes and sort the copy in place.
	sorted := make([]Nonzero, len(nz))
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
	m.coeffMatrix = make([]Nonzero, 0, len(sorted))
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
			prevRow = nz.Row
		}
		index = append(index, C.HighsInt(nz.Col))
		value = append(value, C.double(nz.Value))
	}
	return start, index, value
}

// replaceNilSlices infers the number of rows and columns in a model and
// replaces nil slices with default-valued slices of the appropriate size.  It
// returns the number of rows, the number of columns, and a success flag.  The
// success flag is false if the number of columns is inconsistent across
// non-nil fields.
func (m *commonModel) replaceNilSlices() (int, int, bool) {
	// Infer the number of rows and columns in the model.
	nc, nr := 0, 0
	for _, nz := range m.coeffMatrix {
		if nz.Row >= nr {
			nr = nz.Row + 1
		}
		if nz.Col >= nc {
			nc = nz.Col + 1
		}
	}
	if len(m.colCosts) > nc {
		nc = len(m.colCosts)
	}
	if len(m.colLower) > nc {
		nc = len(m.colLower)
	}
	if len(m.colUpper) > nc {
		nc = len(m.colUpper)
	}
	if len(m.rowLower) > nr {
		nr = len(m.rowLower)
	}
	if len(m.rowUpper) > nr {
		nr = len(m.rowUpper)
	}

	// Replace nil slices with slices of the appropriate size.
	if m.colCosts == nil {
		m.colCosts = make([]float64, nc)
	}
	mInf, pInf := math.Inf(-1), math.Inf(1)
	if m.colLower == nil {
		m.colLower = make([]float64, nc)
		for i := range m.colLower {
			m.colLower[i] = mInf
		}
	}
	if m.colUpper == nil {
		m.colUpper = make([]float64, nc)
		for i := range m.colUpper {
			m.colUpper[i] = pInf
		}
	}
	if m.rowLower == nil {
		m.rowLower = make([]float64, nr)
		for i := range m.rowLower {
			m.rowLower[i] = mInf
		}
	}
	if m.rowUpper == nil {
		m.rowUpper = make([]float64, nr)
		for i := range m.rowUpper {
			m.rowUpper[i] = pInf
		}
	}

	// Complain if any slice is the wrong size.
	switch {
	case len(m.colLower) != nc:
		return 0, 0, false
	case len(m.colLower) != nc, len(m.colUpper) != nc:
		return 0, 0, false
	case len(m.rowLower) != nr, len(m.rowUpper) != nr:
		return 0, 0, false
	}

	// Return the row and column sizes.
	return nr, nc, true
}

// AddRow is a convenience function that lets the caller add to the
// model a single row's lower bound, matrix coefficients (specified
// densely, but stored sparsely), and upper bound.
func (m *commonModel) AddRow(lb float64, coeffs []float64, ub float64) {
	r := len(m.rowLower)
	m.rowLower = append(m.rowLower, lb)
	m.rowUpper = append(m.rowUpper, ub)
	for c, v := range coeffs {
		if v == 0.0 {
			continue
		}
		nz := Nonzero{
			Row:   r,
			Col:   c,
			Value: v,
		}
		m.coeffMatrix = append(m.coeffMatrix, nz)
	}
}

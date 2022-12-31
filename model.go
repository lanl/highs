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

// These are the values a BasisStatus accepts:
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

// These are the values a ModelStatus accepts:
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

// A Model encapsulates fields common to many HiGHS models.  It should not be
// used directly.
type Model struct {
	Maximize    bool      // true=maximize; false=minimize
	ColCosts    []float64 // Column costs (i.e., the objective function itself)
	Offset      float64   // Objective-function constant offset
	ColLower    []float64 // Column lower bounds
	ColUpper    []float64 // Column upper bounds
	RowLower    []float64 // Row lower bounds
	RowUpper    []float64 // Row upper bounds
	CoeffMatrix []Nonzero // Sparse matrix of per-row variable coefficients
}

// filterNonzeros sorts a list of Nonzero elements and removes duplicates.  It
// serves as a helper function for nonzerosToCSR.
func filterNonzeros(nz []Nonzero) ([]Nonzero, error) {
	// Complain about negative indices.
	for _, v := range nz {
		if v.Row < 0 || v.Col < 0 {
			err := fmt.Errorf("(%d, %d) is not a valid coordinate for a matrix coefficient",
				v.Row, v.Col)
			return nil, err
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
	noDups := make([]Nonzero, 0, len(sorted))
	for _, v := range sorted {
		i := len(noDups)
		switch {
		case i == 0:
			// First element: always include.
			noDups = append(noDups, v)
		case v.Row == noDups[i-1].Row && v.Col == noDups[i-1].Col:
			// Duplicate coordinate: retain the later value.
			noDups[i-1].Value = v.Value
		default:
			// New coordinate.
			noDups = append(noDups, v)
		}
	}
	return noDups, nil
}

// nonzerosToCSR converts a list of Nonzero elements to a compressed sparse row
// representation in the form of a set of C vectors accepted by the HiGHS APIs.
func nonzerosToCSR(nz []Nonzero) (start, index []C.HighsInt, value []C.double, err error) {
	// Allocate memory for all of our return vectors.
	var nonzeros []Nonzero
	nonzeros, err = filterNonzeros(nz)
	if err != nil {
		return nil, nil, nil, err
	}
	start = make([]C.HighsInt, 0, len(nonzeros))
	index = make([]C.HighsInt, 0, len(nonzeros))
	value = make([]C.double, 0, len(nonzeros))

	// Construct slices of C types.
	prevRow := -1
	for _, nz := range nonzeros {
		if nz.Row > prevRow {
			start = append(start, C.HighsInt(len(value)))
			prevRow = nz.Row
		}
		index = append(index, C.HighsInt(nz.Col))
		value = append(value, C.double(nz.Value))
	}
	return start, index, value, nil
}

// replaceNilSlices infers the number of rows and columns in a model and
// replaces nil slices with default-valued slices of the appropriate size.  It
// returns the number of rows, the number of columns, a modified copy of the
// Model, and a success flag.  The success flag is false if the number of
// columns is inconsistent across non-nil fields.
func (m *Model) replaceNilSlices() (int, int, *Model, bool) {
	// Infer the number of rows and columns in the model.
	nc, nr := 0, 0
	for _, nz := range m.CoeffMatrix {
		if nz.Row >= nr {
			nr = nz.Row + 1
		}
		if nz.Col >= nc {
			nc = nz.Col + 1
		}
	}
	if len(m.ColCosts) > nc {
		nc = len(m.ColCosts)
	}
	if len(m.ColLower) > nc {
		nc = len(m.ColLower)
	}
	if len(m.ColUpper) > nc {
		nc = len(m.ColUpper)
	}
	if len(m.RowLower) > nr {
		nr = len(m.RowLower)
	}
	if len(m.RowUpper) > nr {
		nr = len(m.RowUpper)
	}

	// Make a copy of the model so we don't modify the original.
	mc := *m

	// Replace nil slices with slices of the appropriate size.
	if mc.ColCosts == nil {
		mc.ColCosts = make([]float64, nc)
	}
	mInf, pInf := math.Inf(-1), math.Inf(1)
	if mc.ColLower == nil {
		mc.ColLower = make([]float64, nc)
		for i := range mc.ColLower {
			mc.ColLower[i] = mInf
		}
	}
	if mc.ColUpper == nil {
		mc.ColUpper = make([]float64, nc)
		for i := range mc.ColUpper {
			mc.ColUpper[i] = pInf
		}
	}
	if mc.RowLower == nil {
		mc.RowLower = make([]float64, nr)
		for i := range mc.RowLower {
			mc.RowLower[i] = mInf
		}
	}
	if mc.RowUpper == nil {
		mc.RowUpper = make([]float64, nr)
		for i := range mc.RowUpper {
			mc.RowUpper[i] = pInf
		}
	}

	// Complain if any slice is the wrong size.
	switch {
	case len(mc.ColLower) != nc:
		return 0, 0, nil, false
	case len(mc.ColLower) != nc, len(mc.ColUpper) != nc:
		return 0, 0, nil, false
	case len(mc.RowLower) != nr, len(mc.RowUpper) != nr:
		return 0, 0, nil, false
	}

	// Return the row and column sizes.
	return nr, nc, &mc, true
}

// AddDenseRow is a convenience function that lets the caller add to the model
// a single row's lower bound, matrix coefficients (specified densely, but
// stored sparsely), and upper bound.
func (m *Model) AddDenseRow(lb float64, coeffs []float64, ub float64) {
	r := len(m.RowLower)
	m.RowLower = append(m.RowLower, lb)
	m.RowUpper = append(m.RowUpper, ub)
	for c, v := range coeffs {
		if v == 0.0 {
			continue
		}
		nz := Nonzero{
			Row:   r,
			Col:   c,
			Value: v,
		}
		m.CoeffMatrix = append(m.CoeffMatrix, nz)
	}
}

// This file provides the high-level Model type and associated functions and
// methods.

package highs

import (
	"errors"
	"fmt"
	"math"
)

// #include "highs-externs.h"
import "C"

// A Model encapsulates all the data needed to express linear-programming
// models, mixed-integer models, and quadratic-programming models.
type Model struct {
	Maximize      bool           // true=maximize; false=minimize
	ColCosts      []float64      // Column costs (i.e., the objective function itself)
	Offset        float64        // Objective-function constant offset
	ColLower      []float64      // Column lower bounds
	ColUpper      []float64      // Column upper bounds
	RowLower      []float64      // Row lower bounds
	RowUpper      []float64      // Row upper bounds
	ConstMatrix   []Nonzero      // Sparse constraint matrix (per-row variable coefficients)
	HessianMatrix []Nonzero      // Sparse, upper-triangular matrix of second partial derivatives of quadratic constraints
	VarTypes      []VariableType // Type of each model variable
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
			Row: r,
			Col: c,
			Val: v,
		}
		m.ConstMatrix = append(m.ConstMatrix, nz)
	}
}

// modelSize returns the number of rows and columns in a model.  It works by
// taking the maximum encountered in any of the fields representing rows or
// columns.
func (m *Model) modelSize() (int, int) {
	nr, nc := 0, 0
	for _, nz := range m.ConstMatrix {
		if nz.Row >= nr {
			nr = nz.Row + 1
		}
		if nz.Col >= nc {
			nc = nz.Col + 1
		}
	}
	for _, nz := range m.HessianMatrix {
		// A Hessian matrix is nc by nc.
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
	if len(m.VarTypes) > nc {
		nc = len(m.VarTypes)
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
	return nr, nc
}

// ToRawModel converts a high-level model to a low-level model.
func (m *Model) ToRawModel() (*RawModel, error) {
	// Construct an empty raw model.  Turn off output, which is out of
	// place in a method like ToRawModel.
	raw := NewRawModel()
	outFlag, err := raw.GetBoolOption("output_flag") // Presumably "true"
	if err != nil {
		return &RawModel{}, err
	}
	err = raw.SetBoolOption("output_flag", false)
	if err != nil {
		return &RawModel{}, err
	}

	// Convert ConstMatrix and HessianMatrix to CSR format.
	aStart, aIndex, aValue, err := nonzerosToCSR(m.ConstMatrix, false)
	if err != nil {
		return &RawModel{}, err
	}
	qStart, qIndex, qValue, err := nonzerosToCSR(m.HessianMatrix, true)
	if err != nil {
		return &RawModel{}, err
	}

	// Convert Go values to C values.
	nr, nc := m.modelSize()
	numCol := C.HighsInt(nc)
	numRow := C.HighsInt(nr)
	numNZ := C.HighsInt(len(aValue))
	qNumNZ := C.HighsInt(len(qValue))
	aFormat := C.kHighsMatrixFormatRowwise
	qFormat := C.kHighsHessianFormatTriangular
	sense := C.kHighsObjSenseMinimize
	if m.Maximize {
		sense = C.kHighsObjSenseMaximize
	}
	offset := C.double(m.Offset)
	colCost := convertSlice[C.double, float64](m.ColCosts)
	colLower := convertSlice[C.double, float64](m.ColLower)
	colUpper := convertSlice[C.double, float64](m.ColUpper)
	rowLower := convertSlice[C.double, float64](m.RowLower)
	rowUpper := convertSlice[C.double, float64](m.RowUpper)
	integrality := make([]C.HighsInt, len(m.VarTypes))
	for i, vt := range m.VarTypes {
		integrality[i] = variableTypeToHighs[vt]
	}

	// Ensure that all slices have consistent lengths.
	var ok bool
	if colCost, ok = expandToLen(nc, colCost, 1.0); !ok {
		return &RawModel{}, fmt.Errorf("inconsistent column counts")
	}
	mInf, pInf := C.double(math.Inf(-1)), C.double(math.Inf(1))
	if colLower, ok = expandToLen(nc, colLower, mInf); !ok {
		return &RawModel{}, fmt.Errorf("inconsistent column counts")
	}
	if colUpper, ok = expandToLen(nc, colUpper, pInf); !ok {
		return &RawModel{}, fmt.Errorf("inconsistent column counts")
	}
	if rowLower, ok = expandToLen(nr, rowLower, mInf); !ok {
		return &RawModel{}, fmt.Errorf("inconsistent row counts")
	}
	if rowUpper, ok = expandToLen(nr, rowUpper, pInf); !ok {
		return &RawModel{}, fmt.Errorf("inconsistent row counts")
	}
	if integrality, ok = expandToLen(nc, integrality, C.kHighsVarTypeContinuous); !ok {
		return &RawModel{}, fmt.Errorf("inconsistent column counts")
	}

	// Construct a low-level model.
	status := C.Highs_passModel(raw.obj, numCol, numRow,
		numNZ, qNumNZ,
		aFormat, qFormat, sense,
		offset, sliceToPointer(colCost),
		sliceToPointer(colLower), sliceToPointer(colUpper),
		sliceToPointer(rowLower), sliceToPointer(rowUpper),
		sliceToPointer(aStart), sliceToPointer(aIndex), sliceToPointer(aValue),
		sliceToPointer(qStart), sliceToPointer(qIndex), sliceToPointer(qValue),
		sliceToPointer(integrality))
	err = newCallStatus(status, "Highs_passModel", "ToRawModel")
	if err != nil {
		return &RawModel{}, err
	}

	// Restore the previous value of output_flag.
	err = raw.SetBoolOption("output_flag", outFlag)
	if err != nil {
		return &RawModel{}, err
	}
	return raw, nil
}

// A Solution encapsulates all the values returned by any of HiGHS's solvers.
// Not all fields will be meaningful when returned by any given solver.
type Solution struct {
	Status       ModelStatus   // Status of the LP solve
	ColumnPrimal []float64     // Primal column solution
	RowPrimal    []float64     // Primal row solution
	ColumnDual   []float64     // Dual column solution
	RowDual      []float64     // Dual row solution
	ColumnBasis  []BasisStatus // Basis status of each column
	RowBasis     []BasisStatus // Basis status of each row
	Objective    float64       // Objective value
}

// Solve solves the model as either an LP, MIP, or QP problem, depending on
// which fields are non-nil.
func (m *Model) Solve() (Solution, error) {
	// Convert the Model to a RawModel.
	var cs CallStatus
	raw, err := m.ToRawModel()
	if err != nil {
		if errors.As(err, &cs) {
			// Hide the fact that ToRawModel was invoked internally.
			cs.GoName = "Solve"
		}
		return Solution{}, err
	}

	// Disable status output.
	err = raw.SetBoolOption("output_flag", false)
	if err != nil {
		if errors.As(err, &cs) {
			// Hide the fact that SetBoolOption was invoked
			// internally.
			cs.GoName = "Solve"
		}
		return Solution{}, err
	}

	// Solve the raw model.
	soln, err := raw.Solve()
	if err != nil {
		return Solution{}, err
	}
	return soln.Solution, nil
}

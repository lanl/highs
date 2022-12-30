// This file provides support for constructing and solving linear-programming
// models.

package highs

import "fmt"

// #include "highs-externs.h"
import "C"

// An LPModel represents a HiGHS linear-programming model.
type LPModel struct {
	commonModel
}

// NewLPModel allocates and returns an empty linear-programming model.
func NewLPModel() *LPModel {
	return &LPModel{}
}

// An LPSolution encapsulates all the values returned by HiGHS's
// linear-programming solver.
type LPSolution struct {
	Status       ModelStatus   // Status of the LP solve
	ColumnPrimal []float64     // Primal column solution
	RowPrimal    []float64     // Primal row solution
	ColumnDual   []float64     // Dual column solution
	RowDual      []float64     // Dual row solution
	ColumnBasis  []BasisStatus // Basis status of each column
	RowBasis     []BasisStatus // Basis status of each row
	Objective    float64       // Objective value
}

// Solve solves a linear-programming model.
func (m *LPModel) Solve() (LPSolution, error) {
	var soln LPSolution

	// Find the model's dimensions.
	nr, nc, ok := m.replaceNilSlices()
	if !ok {
		return soln, fmt.Errorf("model has inconsistent dimensions")
	}

	// Convert the coefficient matrix to HiGHS format.
	aStart, aIndex, aValue := m.makeSparseMatrix()

	// Convert other model parameters from Go data types to C data types.
	numCol := C.HighsInt(nc)
	numRow := C.HighsInt(nr)
	numNZ := C.HighsInt(len(aValue))
	aFormat := C.kHighsMatrixFormatRowwise // Column-wise is not currently supported.
	sense := C.kHighsObjSenseMinimize
	if m.maximize {
		sense = C.kHighsObjSenseMaximize
	}
	offset := C.double(m.offset)
	colCost := convertSlice[C.double, float64](m.colCosts)
	colLower := convertSlice[C.double, float64](m.colLower)
	colUpper := convertSlice[C.double, float64](m.colUpper)
	rowLower := convertSlice[C.double, float64](m.rowLower)
	rowUpper := convertSlice[C.double, float64](m.rowUpper)

	// Allocate storage for return values.
	colValue := make([]C.double, nc)
	colDual := make([]C.double, nc)
	rowValue := make([]C.double, nr)
	rowDual := make([]C.double, nr)
	colBasisStatus := make([]C.HighsInt, nc)
	rowBasisStatus := make([]C.HighsInt, nr)
	var modelStatus C.HighsInt

	// We finally can invoke Highs_lpCall!
	status := C.Highs_lpCall(numCol, numRow, numNZ,
		aFormat, sense, offset,
		&colCost[0], &colLower[0], &colUpper[0],
		&rowLower[0], &rowUpper[0],
		&aStart[0], &aIndex[0], &aValue[0],
		&colValue[0], &colDual[0],
		&rowValue[0], &rowDual[0],
		&colBasisStatus[0], &rowBasisStatus[0],
		&modelStatus)
	err := convertHighsStatusToError(status, "Solve")
	if err != nil {
		return soln, err
	}

	// Convert C return types to Go types.
	soln.Status = convertHighsModelStatus(modelStatus)
	soln.ColumnPrimal = convertSlice[float64, C.double](colValue)
	soln.RowPrimal = convertSlice[float64, C.double](rowValue)
	soln.ColumnDual = convertSlice[float64, C.double](colDual)
	soln.RowDual = convertSlice[float64, C.double](rowDual)
	soln.ColumnBasis = make([]BasisStatus, nc)
	for i, cbs := range colBasisStatus {
		soln.ColumnBasis[i] = convertHighsBasisStatus(cbs)
	}
	soln.RowBasis = make([]BasisStatus, nr)
	for i, rbs := range rowBasisStatus {
		soln.RowBasis[i] = convertHighsBasisStatus(rbs)
	}

	// Compute the objective value as a convenience for the user.
	soln.Objective = m.offset
	for i, cp := range soln.ColumnPrimal {
		soln.Objective += cp * m.colCosts[i]
	}
	return soln, nil
}

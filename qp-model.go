// This file provides support for constructing and solving quadratic-programming
// models.

package highs

import "fmt"

// #include "highs-externs.h"
import "C"

// A QPModel represents a HiGHS quadratic-programming model.
type QPModel struct {
	Model
	HessianMatrix []Nonzero // Sparse, upper-triangular matrix of second partial derivatives of quadratic constraints
}

// An QPSolution encapsulates all the values returned by HiGHS's
// linear-programming solver.
type QPSolution struct {
	Status       ModelStatus   // Status of the QP solve
	ColumnPrimal []float64     // Primal column solution
	RowPrimal    []float64     // Primal row solution
	ColumnDual   []float64     // Dual column solution
	RowDual      []float64     // Dual row solution
	ColumnBasis  []BasisStatus // Basis status of each column
	RowBasis     []BasisStatus // Basis status of each row
	Objective    float64       // Objective value
}

// Solve solves a linear-programming model.
func (m *QPModel) Solve() (QPSolution, error) {
	var soln QPSolution

	// Find the model's dimensions.
	nr, nc, cm, ok := m.replaceNilSlices()
	if !ok {
		return soln, fmt.Errorf("model has inconsistent dimensions")
	}

	// Convert the coefficient matrix to HiGHS format.
	aStart, aIndex, aValue, err := nonzerosToCSR(m.CoeffMatrix, false)
	if err != nil {
		return soln, err
	}

	// Convert the Hessian matrix to HiGHS format.
	hessian, err := filterNonzeros(m.HessianMatrix, true) // Needed below
	if err != nil {
		return soln, err
	}
	qStart, qIndex, qValue, err := nonzerosToCSR(hessian, true)
	if err != nil {
		return soln, err
	}

	// Convert other model parameters from Go data types to C data types.
	numCol := C.HighsInt(nc)
	numRow := C.HighsInt(nr)
	numNZ := C.HighsInt(len(aValue))
	qNumNZ := C.HighsInt(len(qValue))
	aFormat := C.kHighsMatrixFormatRowwise     // Column-wise is not currently supported.
	qFormat := C.kHighsHessianFormatTriangular // Square is not currently supported.
	sense := C.kHighsObjSenseMinimize
	if cm.Maximize {
		sense = C.kHighsObjSenseMaximize
	}
	offset := C.double(cm.Offset)
	colCost := convertSlice[C.double, float64](cm.ColCosts)
	colLower := convertSlice[C.double, float64](cm.ColLower)
	colUpper := convertSlice[C.double, float64](cm.ColUpper)
	rowLower := convertSlice[C.double, float64](cm.RowLower)
	rowUpper := convertSlice[C.double, float64](cm.RowUpper)

	// Allocate storage for return values.
	colValue := make([]C.double, nc)
	colDual := make([]C.double, nc)
	rowValue := make([]C.double, nr)
	rowDual := make([]C.double, nr)
	colBasisStatus := make([]C.HighsInt, nc)
	rowBasisStatus := make([]C.HighsInt, nr)
	var modelStatus C.HighsInt

	// We finally can invoke Highs_qpCall!
	status := C.Highs_qpCall(numCol, numRow, numNZ, qNumNZ,
		aFormat, qFormat, sense, offset,
		&colCost[0], &colLower[0], &colUpper[0],
		&rowLower[0], &rowUpper[0],
		&aStart[0], &aIndex[0], &aValue[0],
		&qStart[0], &qIndex[0], &qValue[0],
		&colValue[0], &colDual[0],
		&rowValue[0], &rowDual[0],
		&colBasisStatus[0], &rowBasisStatus[0],
		&modelStatus)
	err = newCallStatus(status, "Highs_qpCall", "Solve")
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
	soln.Objective = cm.Offset // Constant term
	for i, cp := range soln.ColumnPrimal {
		// Linear terms
		soln.Objective += cp * cm.ColCosts[i]
	}
	for _, nz := range hessian {
		// Quadratic terms
		r, c := nz.Row, nz.Col
		v := soln.ColumnPrimal[r] * soln.ColumnPrimal[c] * nz.Value / 2.0
		if r != c {
			v *= 2.0
		}
		soln.Objective += v
	}
	return soln, nil
}

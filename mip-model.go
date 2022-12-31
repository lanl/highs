// This file provides support for constructing and solving mixed-integer
// programming models.

package highs

import "fmt"

// #include "highs-externs.h"
import "C"

// A MIPModel represents a HiGHS mixed-integer programming model.
type MIPModel struct {
	Model
	VarTypes []VariableType // Type of each model variable
}

// A VariableType indicates the type of a model variable.
type VariableType int

// These are the values a VariableType accepts:
const (
	ContinuousType VariableType = iota
	IntegerType
	SemiContinuousType
	SemiIntegerType
	ImplicitIntegerType
)

// variableTypeToHighs maps a VariableType to a kHighsVarType.  This slice must
// be kept up to date with the VariableType constants.
var variableTypeToHighs = []C.HighsInt{
	C.kHighsVarTypeContinuous,
	C.kHighsVarTypeInteger,
	C.kHighsVarTypeSemiContinuous,
	C.kHighsVarTypeSemiInteger,
	C.kHighsVarTypeImplicitInteger,
}

// A MIPSolution encapsulates all the values returned by HiGHS's mixed-integer
// programming solver.
type MIPSolution struct {
	Status       ModelStatus // Status of the LP solve
	ColumnPrimal []float64   // Primal column solution
	RowPrimal    []float64   // Primal row solution
	Objective    float64     // Objective value
}

// Solve solves a mixed-integer programming model.
func (m *MIPModel) Solve() (MIPSolution, error) {
	var soln MIPSolution

	// Find the model's dimensions.
	nr, nc, cm, ok := m.replaceNilSlices()
	if !ok {
		return soln, fmt.Errorf("model has inconsistent dimensions")
	}

	// Convert the coefficient matrix to HiGHS format.
	aStart, aIndex, aValue, err := nonzerosToCSR(m.CoeffMatrix)
	if err != nil {
		return soln, err
	}

	// Convert other model parameters from Go data types to C data types.
	numCol := C.HighsInt(nc)
	numRow := C.HighsInt(nr)
	numNZ := C.HighsInt(len(aValue))
	aFormat := C.kHighsMatrixFormatRowwise // Column-wise is not currently supported.
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
	integrality := make([]C.HighsInt, nc) // Defaults to ContinuousType
	for i, vt := range m.VarTypes {
		integrality[i] = variableTypeToHighs[vt]
	}

	// Allocate storage for return values.
	colValue := make([]C.double, nc)
	rowValue := make([]C.double, nr)
	var modelStatus C.HighsInt

	// We finally can invoke Highs_mipCall!
	status := C.Highs_mipCall(numCol, numRow, numNZ,
		aFormat, sense, offset,
		&colCost[0], &colLower[0], &colUpper[0],
		&rowLower[0], &rowUpper[0],
		&aStart[0], &aIndex[0], &aValue[0],
		&integrality[0],
		&colValue[0], &rowValue[0],
		&modelStatus)
	err = newHighsStatus(status, "Highs_mipCall", "Solve")
	if err != nil {
		return soln, err
	}

	// Convert C return types to Go types.
	soln.Status = convertHighsModelStatus(modelStatus)
	soln.ColumnPrimal = convertSlice[float64, C.double](colValue)
	soln.RowPrimal = convertSlice[float64, C.double](rowValue)

	// Compute the objective value as a convenience for the user.
	soln.Objective = cm.Offset
	for i, cp := range soln.ColumnPrimal {
		soln.Objective += cp * cm.ColCosts[i]
	}
	return soln, nil
}

// This file provides support for constructing and solving models using HiGHS's
// "full" (low-level) API.

package highs

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"unsafe"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <interfaces/highs_c_api.h>
import "C"

// A RawModel represents a HiGHS low-level model.
type RawModel struct {
	obj unsafe.Pointer
}

// NewRawModel allocates and returns an empty raw model.
func NewRawModel() *RawModel {
	model := &RawModel{}
	model.obj = C.Highs_create()
	runtime.SetFinalizer(model, func(m *RawModel) {
		C.Highs_destroy(m.obj)
	})
	return model
}

// ReadModelFromFile overwrites the model with a model read in MPS format from
// a named file.
func (m *RawModel) ReadModelFromFile(fn string) error {
	// Convert the filename argument from Go to C.
	fName := C.CString(fn)
	defer C.free(unsafe.Pointer(fName))

	// Read into the model.
	status := C.Highs_readModel(m.obj, fName)
	return newCallStatus(status, "Highs_readModel", "ReadModelFromFile")
}

// ReadModel overwrites the model with a model read in MPS format from an
// io.Reader.
func (m *RawModel) ReadModel(r io.Reader) error {
	// Copy from the reader to a throwaway file.
	tFile, err := os.CreateTemp("", "highs-*.mps")
	if err != nil {
		return err
	}
	fName := tFile.Name()
	defer os.Remove(fName)
	_, err = io.Copy(tFile, r)
	if err != nil {
		return err
	}
	err = tFile.Close()
	if err != nil {
		return err
	}

	// Convert the throwaway filename from Go to C.
	cFName := C.CString(fName)
	defer C.free(unsafe.Pointer(cFName))

	// Read into the model.
	status := C.Highs_readModel(m.obj, cFName)
	return newCallStatus(status, "Highs_readModel", "ReadModel")
}

// WriteModelToFile writes a model in MPS format to a named file.
func (m *RawModel) WriteModelToFile(fn string) error {
	// Convert the filename argument from Go to C.
	cFName := C.CString(fn)
	defer C.free(unsafe.Pointer(cFName))

	// Write the model.
	status := C.Highs_writeModel(m.obj, cFName)
	return newCallStatus(status, "Highs_writeModel", "WriteModelToFile")
}

// WriteModel writes a model in MPS format to an io.Writer.
func (m *RawModel) WriteModel(w io.Writer) error {
	// Create a throwaway file to use as a staging area.
	tFile, err := os.CreateTemp("", "highs-*.mps")
	if err != nil {
		return err
	}
	fName := tFile.Name()
	defer os.Remove(fName)
	err = tFile.Close()
	if err != nil {
		return err
	}

	// Convert the throwaway filename from Go to C.
	cFName := C.CString(fName)
	defer C.free(unsafe.Pointer(cFName))

	// Write the model to the throwaway file.
	status := C.Highs_writeModel(m.obj, cFName)
	err = newCallStatus(status, "Highs_writeModel", "WriteModel")

	// Ignore warnings (common for Highs_writeModel).
	var cs CallStatus
	if errors.As(err, &cs) {
		if !cs.IsWarning() {
			return err
		}
	} else if err != nil {
		return err
	}

	// Copy the contents of the throwaway file to the io.Writer.
	tFile, err = os.Open(fName)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, tFile)
	if err != nil {
		return err
	}
	err = tFile.Close()
	if err != nil {
		return err
	}
	return cs // Propagate any warnings.
}

// SetBoolOption assigns a Boolean value to a named option.
func (m *RawModel) SetBoolOption(opt string, v bool) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	var val C.HighsInt
	if v {
		val = 1
	}

	// Set the option.
	status := C.Highs_setBoolOptionValue(m.obj, str, val)
	return newCallStatus(status, "Highs_setBoolOptionValue", "SetBoolOption")
}

// SetIntOption assigns an integer value to a named option.
func (m *RawModel) SetIntOption(opt string, v int) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	val := C.HighsInt(v)

	// Set the option.
	status := C.Highs_setIntOptionValue(m.obj, str, val)
	return newCallStatus(status, "Highs_setIntOptionValue", "SetIntOption")
}

// SetFloat64Option assigns a floating-point value to a named option.
func (m *RawModel) SetFloat64Option(opt string, v float64) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	val := C.double(v)

	// Set the option.
	status := C.Highs_setDoubleOptionValue(m.obj, str, val)
	return newCallStatus(status, "Highs_setDoubleOptionValue", "SetFloat64Option")
}

// SetStringOption assigns a string value to a named option.
func (m *RawModel) SetStringOption(opt string, v string) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	val := C.CString(v)
	defer C.free(unsafe.Pointer(val))

	// Set the option.
	status := C.Highs_setStringOptionValue(m.obj, str, val)
	return newCallStatus(status, "Highs_setStringOptionValue", "SetStringOption")
}

// GetBoolOption returns the Boolean value of a named option.
func (m *RawModel) GetBoolOption(opt string) (bool, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.HighsInt
	status := C.Highs_getBoolOptionValue(m.obj, str, &val)
	err := newCallStatus(status, "Highs_getBoolOptionValue", "GetBoolOption")
	if err != nil {
		return false, err
	}
	var v bool
	if val != 0 {
		v = true
	}
	return v, nil
}

// GetIntOption returns the Integer value of a named option.
func (m *RawModel) GetIntOption(opt string) (int, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.HighsInt
	status := C.Highs_getIntOptionValue(m.obj, str, &val)
	err := newCallStatus(status, "Highs_getIntOptionValue", "GetIntOption")
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

// GetFloat64Option returns the floating-point value of a named option.
func (m *RawModel) GetFloat64Option(opt string) (float64, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.double
	status := C.Highs_getDoubleOptionValue(m.obj, str, &val)
	err := newCallStatus(status, "Highs_getDoubleOptionValue", "GetFloat64Option")
	if err != nil {
		return 0.0, err
	}
	return float64(val), nil
}

// GetStringOption returns the string value of a named option.  Do not invoke
// this method in security-sensitive applications because it runs a risk of
// buffer overflow.
func (m *RawModel) GetStringOption(opt string) (string, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// The value could potentially be of any size.  Allocate "enough"
	// memory and hope for the best.
	val := (*C.char)(C.calloc(65536, 1))
	defer C.free(unsafe.Pointer(val))

	// Get the value.
	status := C.Highs_getStringOptionValue(m.obj, str, val)
	err := newCallStatus(status, "Highs_getStringOptionValue", "GetStringOption")
	if err != nil {
		return "", err
	}
	return C.GoString(val), nil
}

// SetMaximization tells a model to maximize (true) or minimize (false) its
// objective function.
func (m *RawModel) SetMaximization(max bool) error {
	var sense C.HighsInt = C.kHighsObjSenseMinimize
	if max {
		sense = C.kHighsObjSenseMaximize
	}
	status := C.Highs_changeObjectiveSense(m.obj, sense)
	return newCallStatus(status, "Highs_changeObjectiveSense", "SetMaximization")
}

// SetColumnCosts specifies a model's column costs (i.e., its objective
// function).
func (m *RawModel) SetColumnCosts(cs []float64) error {
	cost := convertSlice[C.double, float64](cs)
	status := C.Highs_changeColsCostByRange(m.obj,
		0, C.HighsInt(len(cs)-1),
		&cost[0])
	return newCallStatus(status, "Highs_changeColsCostByRange", "SetColumnCosts")
}

// SetOffset specifies a constant offset for the objective function.
func (m *RawModel) SetOffset(o float64) error {
	status := C.Highs_changeObjectiveOffset(m.obj, C.double(o))
	return newCallStatus(status, "Highs_changeObjectiveOffset", "SetOffset")
}

// prepareBounds replaces nil column or row bounds with infinities.
func prepareBounds(lb, ub []float64) ([]float64, []float64, error) {
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
		return nil, nil, fmt.Errorf("different numbers of lower and upper bounds were provided (%d vs. %d)", len(lb), len(ub))
	}
	return lb, ub, nil
}

// AddColumnBounds appends to a model's lower and upper column bounds.  If the
// lower-bound argument is nil it is replaced with a slice of negative
// infinities.  If the upper-bound argument is nil, it is replaced with a slice
// of positive infinities.
func (m *RawModel) AddColumnBounds(lb, ub []float64) error {
	colLower, colUpper, err := prepareBounds(lb, ub)
	if err != nil {
		return err
	}
	lower := convertSlice[C.double, float64](colLower)
	upper := convertSlice[C.double, float64](colUpper)
	status := C.Highs_addVars(m.obj, C.HighsInt(len(lower)),
		&lower[0], &upper[0])
	return newCallStatus(status, "Highs_addVars", "SetColumnBounds")
}

// AddCompSparseRows appends compressed sparse rows to the model.
func (m *RawModel) AddCompSparseRows(lb []float64, start []int, index []int, value []float64, ub []float64) error {
	// Check for simple errors.
	if len(lb) != len(ub) {
		return fmt.Errorf("lb and ub must be the same length (%d vs. %d)",
			len(lb), len(ub))
	}
	if len(index) != len(value) {
		return fmt.Errorf("index and value must be the same length (%d vs. %d)",
			len(index), len(value))
	}

	// Invoke the HiGHS API.
	hLower := convertSlice[C.double, float64](lb)
	hUpper := convertSlice[C.double, float64](ub)
	hStart := convertSlice[C.HighsInt, int](start)
	hIndex := convertSlice[C.HighsInt, int](index)
	hValue := convertSlice[C.double, float64](value)
	status := C.Highs_addRows(m.obj, C.HighsInt(len(lb)),
		&hLower[0], &hUpper[0],
		C.HighsInt(len(value)), &hStart[0], &hIndex[0], &hValue[0])
	return newCallStatus(status, "Highs_addRows", "AddCompSparseRows")
}

// AddDenseRow is a convenience function that lets the caller add to the model
// a single row's lower bound, matrix coefficients (specified densely, but
// stored sparsely), and upper bound.
func (m *RawModel) AddDenseRow(lb float64, coeffs []float64, ub float64) error {
	// Convert dense to sparse.
	var numNewNz C.HighsInt
	index := make([]C.HighsInt, 0, len(coeffs))
	value := make([]C.double, 0, len(coeffs))
	for i, v := range coeffs {
		if v == 0.0 {
			continue
		}
		index = append(index, C.HighsInt(i))
		value = append(value, C.double(v))
		numNewNz++
	}

	// Add the row.
	status := C.Highs_addRow(m.obj, C.double(lb), C.double(ub),
		numNewNz, &index[0], &value[0])
	return newCallStatus(status, "Highs_addRow", "AddDenseRow")
}

// SetIntegrality specifies the type of each column (variable) in the model.
func (m *RawModel) SetIntegrality(ts []VariableType) error {
	integrality := make([]C.HighsInt, len(ts))
	for i, t := range ts {
		integrality[i] = variableTypeToHighs[t]
	}
	status := C.Highs_changeColsIntegralityByRange(m.obj,
		0, C.HighsInt(len(integrality)-1),
		&integrality[0])
	return newCallStatus(status, "Highs_changeColsIntegralityByRange", "SetIntegrality")
}

// AddCompSparseHessian assigns a Hessian in compressed sparse row form to the
// model.  This is used to formulate quadratic constraints in a
// quadratic-programming model.
func (m *RawModel) AddCompSparseHessian(start []int, index []int, value []float64) error {
	// Check for simple errors.
	if len(index) != len(value) {
		return fmt.Errorf("index and value must be the same length (%d vs. %d)",
			len(index), len(value))
	}

	// Invoke the HiGHS API.
	hStart := convertSlice[C.HighsInt, int](start)
	hIndex := convertSlice[C.HighsInt, int](index)
	hValue := convertSlice[C.double, float64](value)
	status := C.Highs_passHessian(m.obj, C.HighsInt(len(start)),
		C.HighsInt(len(value)), C.kHighsHessianFormatTriangular,
		&hStart[0], &hIndex[0], &hValue[0])
	return newCallStatus(status, "Highs_passHessian", "AddCompSparseHessian")
}

// Solve solves a model.
func (m *RawModel) Solve() (*RawSolution, error) {
	// Solve the model.  We assume the user has already set up all the
	// required parameters.
	status := C.Highs_run(m.obj)
	err := newCallStatus(status, "Highs_run", "Solve")
	if err != nil {
		return &RawSolution{}, err
	}

	// Extract the solution as Go data.
	var soln RawSolution
	soln.obj = m.obj
	soln.Status = convertHighsModelStatus(C.Highs_getModelStatus(soln.obj))
	nc := int(C.Highs_getNumCol(soln.obj))
	nr := int(C.Highs_getNumRow(soln.obj))
	colValue := make([]C.double, nc)
	colDual := make([]C.double, nc)
	rowValue := make([]C.double, nr)
	rowDual := make([]C.double, nr)
	status = C.Highs_getSolution(soln.obj, &colValue[0], &colDual[0],
		&rowValue[0], &rowDual[0])
	err = newCallStatus(status, "Highs_getSolution", "Solve")
	if err != nil {
		return &RawSolution{}, err
	}
	soln.ColumnPrimal = convertSlice[float64, C.double](colValue)
	soln.RowPrimal = convertSlice[float64, C.double](rowValue)
	soln.Objective, err = soln.GetFloat64Info("objective_function_value")
	if err != nil {
		return &RawSolution{}, err
	}

	// Assign dual slices only if the dual-solution status is "feasible".
	dss, err := soln.GetIntInfo("dual_solution_status")
	if err != nil {
		return &RawSolution{}, err
	}
	if dss == int(C.kHighsSolutionStatusFeasible) {
		soln.ColumnDual = convertSlice[float64, C.double](colDual)
		soln.RowDual = convertSlice[float64, C.double](rowDual)
	}

	// If basis data are available, convert them from C to Go.
	bValid, err := soln.GetIntInfo("basis_validity")
	if err == nil && bValid == int(C.kHighsBasisValidityValid) {
		colBasisStatus := make([]C.HighsInt, nc)
		rowBasisStatus := make([]C.HighsInt, nr)
		status = C.Highs_getBasis(soln.obj, &colBasisStatus[0], &rowBasisStatus[0])
		err = newCallStatus(status, "Highs_getBasis", "Solve")
		if err != nil {
			return &RawSolution{}, err
		}
		soln.ColumnBasis = make([]BasisStatus, nc)
		for i, cbs := range colBasisStatus {
			soln.ColumnBasis[i] = convertHighsBasisStatus(cbs)
		}
		soln.RowBasis = make([]BasisStatus, nr)
		for i, rbs := range rowBasisStatus {
			soln.RowBasis[i] = convertHighsBasisStatus(rbs)
		}
	}
	return &soln, nil
}

// This file provides support for manipulating solutions returned from a HiGHS
// solver using the "full" (low-level) API.

package highs

import (
	"unsafe"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include "highs-externs.h"
import "C"

// A RawSolution encapsulates all the values returned by various HiGHS solvers
// and provides methods to retrieve additional information.
type RawSolution struct {
	obj      unsafe.Pointer // Underlying opaque highs data type
	Solution                // Values returned by the solver
}

// GetIntInfo returns the integer value of a named piece of information.
func (s *RawSolution) GetIntInfo(info string) (int, error) {
	// Convert the info argument from Go to C.
	str := C.CString(info)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.HighsInt
	status := C.Highs_getIntInfoValue(s.obj, str, &val)
	err := newCallStatus(status, "Highs_getIntInfoValue", "GetIntInfo")
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

// GetInt64Info returns the 64-bit integer value of a named piece of
// information.
func (s *RawSolution) GetInt64Info(info string) (int64, error) {
	// Convert the info argument from Go to C.
	str := C.CString(info)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.int64_t
	status := C.Highs_getInt64InfoValue(s.obj, str, &val)
	err := newCallStatus(status, "Highs_getInt64InfoValue", "GetInt64Info")
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// GetFloat64Info returns the floating-point value of a named piece of
// information.
func (s *RawSolution) GetFloat64Info(info string) (float64, error) {
	// Convert the info argument from Go to C.
	str := C.CString(info)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.double
	status := C.Highs_getDoubleInfoValue(s.obj, str, &val)
	err := newCallStatus(status, "Highs_getDoubleInfoValue", "GetFloat64Info")
	if err != nil {
		return 0.0, err
	}
	return float64(val), nil
}

// This file provides support for manipulating solutions returned from a HiGHS
// solver using the "full" (low-level) API.

package highs

import (
	"io"
	"os"
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

// WriteSolutionToFile writes a textual version of the solution to a named
// file.  If the second argument is false, WriteSolutiontoFile will use a more
// computer-friendly format; if true, it will use a more human-friendly format.
func (s *RawSolution) WriteSolutionToFile(fn string, pretty bool) error {
	// Convert the filename argument from Go to C.
	cFName := C.CString(fn)
	defer C.free(unsafe.Pointer(cFName))

	// Write the solution.
	if pretty {
		status := C.Highs_writeSolutionPretty(s.obj, cFName)
		return newCallStatus(status, "Highs_writeSolutionPretty", "WriteSolutionToFile")
	}
	status := C.Highs_writeSolution(s.obj, cFName)
	return newCallStatus(status, "Highs_writeSolution", "WriteSolutionToFile")
}

// WriteSolution writes a textual version of the solution to an io.Writer.  If
// the second argument is false, WriteSolutiontoFile will use a more
// computer-friendly format; if true, it will use a more human-friendly format.
func (s *RawSolution) WriteSolution(w io.Writer, pretty bool) error {
	// Create a throwaway file to use as a staging area.
	tFile, err := os.CreateTemp("", "highs-*.txt")
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

	// Write the solution to the throwaway file.
	if pretty {
		status := C.Highs_writeSolutionPretty(s.obj, cFName)
		err = newCallStatus(status, "Highs_writeSolutionPretty", "WriteSolution")
	} else {
		status := C.Highs_writeSolution(s.obj, cFName)
		err = newCallStatus(status, "Highs_writeSolution", "WriteSolution")
	}
	if err != nil {
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
	return nil
}

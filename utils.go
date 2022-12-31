// This file provides miscellaneous utility functions.

package highs

import (
	"golang.org/x/exp/constraints"
)

// #include "highs-externs.h"
import "C"

// A CallStatus wraps a kHighsStatus returned by a call to HiGHS.  A CallStatus
// may be an error or just a warning.
type CallStatus struct {
	Status int    // kHighsStatus value
	CName  string // Name of the HiGHS function that returned a non-Ok status
	GoName string // Name of the highs package function that called the CName function
}

// Error returns a CallStatus as a string.
func (e CallStatus) Error() string {
	switch e.Status {
	case int(C.kHighsStatusError):
		return "%s failed with an error"
	case int(C.kHighsStatusWarning):
		return "%s failed with a warning"
	default:
		return "%s failed with an unknown status"
	}
}

// IsWarning returns true if the CallStatus is merely a warning.
func (e CallStatus) IsWarning() bool {
	return e.Status == int(C.kHighsStatusWarning)
}

// newCallStatus constructs a CallStatus or returns nil if the status
// is kHighsStatusOk.
func newCallStatus(st C.HighsInt, hName, gName string) error {
	if st == C.kHighsStatusOk {
		return nil
	}
	return CallStatus{
		Status: int(st),
		CName:  hName,
		GoName: gName,
	}
}

// A numeric is any integer or any floating-point type.
type numeric interface {
	constraints.Integer | constraints.Float
}

// convertSlice is a helper function that converts a slice from one type to
// another.
func convertSlice[T, F numeric](from []F) []T {
	to := make([]T, len(from))
	for i, f := range from {
		to[i] = T(f)
	}
	return to
}

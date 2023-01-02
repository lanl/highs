// This file provides miscellaneous utility functions.

package highs

import (
	"fmt"
	"sort"

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
		return fmt.Sprintf("%s failed with an error", e.GoName)
	case int(C.kHighsStatusWarning):
		return fmt.Sprintf("%s completed with a warning", e.GoName)
	default:
		return fmt.Sprintf("%s exited with an unknown status", e.GoName)
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

// filterNonzeros sorts a list of Nonzero elements, removes duplicates, and, if
// tri is true, rejects lower-triangular elements.  filterNonzeros serves as a
// helper function for nonzerosToCSR.
func filterNonzeros(nz []Nonzero, tri bool) ([]Nonzero, error) {
	// Complain about negative indices.
	for _, v := range nz {
		if v.Row < 0 || v.Col < 0 {
			err := fmt.Errorf("(%d, %d) is not a valid coordinate for a matrix coefficient",
				v.Row, v.Col)
			return nil, err
		}
	}

	// Optionally complain about lower-triangular indices.
	if tri {
		for _, v := range nz {
			if v.Row > v.Col {
				err := fmt.Errorf("(%d, %d) is not a valid upper-triangular coordinate for a matrix coefficient",
					v.Row, v.Col)
				return nil, err
			}
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
			noDups[i-1].Val = v.Val
		default:
			// New coordinate.
			noDups = append(noDups, v)
		}
	}
	return noDups, nil
}

// nonzerosToCSR converts a list of Nonzero elements to a compressed sparse row
// representation in the form of a set of C vectors accepted by the HiGHS APIs.
func nonzerosToCSR(nz []Nonzero, tri bool) (start, index []C.HighsInt, value []C.double, err error) {
	// Allocate memory for all of our return vectors.
	var nonzeros []Nonzero
	nonzeros, err = filterNonzeros(nz, tri)
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
		value = append(value, C.double(nz.Val))
	}
	return start, index, value, nil
}

// expandToLen takes a length, a slice, and a value.  If the slice has the
// given length, it returns the slice unmodified.  If the slice has length
// zero, it returns a length-sized slice of value.  If the slice has any other
// length, it returns a failure code.
func expandToLen[T any](n int, xs []T, v T) ([]T, bool) {
	switch {
	case len(xs) == n:
		return xs, true
	case len(xs) == 0:
		ys := make([]T, n)
		for i := range ys {
			ys[i] = v
		}
		return ys, true
	default:
		return nil, false
	}
}

// sliceToPointer returns a pointer to the first element of a slice or nil if
// the slice is empty.
func sliceToPointer[T any](xs []T) *T {
	if len(xs) == 0 {
		return nil
	}
	return &xs[0]
}

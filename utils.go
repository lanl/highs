// This file provides miscellaneous utility functions.

package highs

import (
	"golang.org/x/exp/constraints"
)

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

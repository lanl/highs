// This file provides miscellaneous utility functions useful for unit tests.

package highs

import (
	"errors"
	"math"
	"testing"
)

// compSlices compares two slices for equality.
func compSlices[AType, EType numeric](t *testing.T, name string, act []AType, exp []EType) {
	if len(act) != len(exp) {
		t.Fatalf("%s: expected %v but observed %v", name, exp, act)
	}
	for i, e := range exp {
		if EType(act[i]) != e {
			t.Fatalf("%s: expected %v but observed %v", name, exp, act)
		}
	}
}

// roundFloats rounds a list of floats to a given precision.
func roundFloats(e float64, xs []float64) []float64 {
	rs := make([]float64, len(xs))
	for i, x := range xs {
		rs[i] = math.Round(x/e) * e
	}
	return rs
}

// checkErr calls another function and aborts the test if it returns a non-nil
// error.
func checkErr(t *testing.T, e error) {
	if e == nil {
		return
	}
	var cs CallStatus
	if errors.As(e, &cs) && !cs.IsWarning() {
		t.Fatal(e)
	}
}

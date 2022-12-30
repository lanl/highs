// This file provides miscellaneous utility functions useful for unit tests.

package highs

import "testing"

// compSlices is a helper function for unit tests that compares two slices for
// equality.
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

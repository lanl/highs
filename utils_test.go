// This file provides miscellaneous utility functions useful for unit tests.

package highs

import (
	"errors"
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

// checkErr calls another function and aborts the test if it returns a non-nil
// error.
func checkErr(t *testing.T, e error) {
	if e == nil {
		return
	}
	var hs HighsStatus
	if errors.As(e, &hs) && !hs.IsWarning() {
		t.Fatal(e)
	}
}

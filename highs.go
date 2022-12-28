// Package highs provides a Go interface to the HiGHS optimizer.
package highs

// #cgo pkg-config: highs
// #include <interfaces/highs_c_api.h>
import "C"

//go:generate stringer -type=BasisStatus
//go:generate stringer -type=ModelStatus
//go:generate stringer -type=VariableType

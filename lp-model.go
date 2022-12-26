// This file provides support for constructing and solving
// linear-programming models.

package highs

// An LPModel represents a HiGHS linear-programming model.
type LPModel struct {
	commonModel
}

// NewLPModel allocates and returns an empty linear-programming model.
func NewLPModel() *LPModel {
	return &LPModel{}
}

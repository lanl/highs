// This file provides a set of simple types needed to construct a model.

package highs

// #include "highs-externs.h"
import "C"

// A Nonzero represents a nonzero entry in a sparse matrix.  Rows and columns
// are indexed from zero.
type Nonzero struct {
	Row int
	Col int
	Val float64
}

// A BasisStatus represents the basis status of a row or column.
type BasisStatus int

// These are the values a BasisStatus accepts:
const (
	UnknownBasisStatus BasisStatus = iota
	Lower
	Basic
	Upper
	Zero
	NonBasic
)

// convertHighsBasisStatus converts a kHighsBasisStatus to a BasisStatus.
func convertHighsBasisStatus(hbs C.HighsInt) BasisStatus {
	switch hbs {
	case C.kHighsBasisStatusLower:
		return Lower
	case C.kHighsBasisStatusBasic:
		return Basic
	case C.kHighsBasisStatusUpper:
		return Upper
	case C.kHighsBasisStatusZero:
		return Zero
	case C.kHighsBasisStatusNonbasic:
		return NonBasic
	default:
		return UnknownBasisStatus
	}
}

//go:generate stringer -type=BasisStatus

// A ModelStatus represents the status of an attempt to solve a model.
type ModelStatus int

// These are the values a ModelStatus accepts:
const (
	UnknownModelStatus ModelStatus = iota
	NotSet
	LoadError
	ModelError
	PresolveError
	SolveError
	PostsolveError
	ModelEmpty
	Optimal
	Infeasible
	UnboundedOrInfeasible
	Unbounded
	ObjectiveBound
	ObjectiveTarget
	TimeLimit
	IterationLimit
)

// convertHighsModelStatus converts a kHighsModelStatus to a ModelStatus.
func convertHighsModelStatus(hms C.HighsInt) ModelStatus {
	switch hms {
	case C.kHighsModelStatusNotset:
		return NotSet
	case C.kHighsModelStatusLoadError:
		return LoadError
	case C.kHighsModelStatusModelError:
		return ModelError
	case C.kHighsModelStatusPresolveError:
		return PresolveError
	case C.kHighsModelStatusSolveError:
		return SolveError
	case C.kHighsModelStatusPostsolveError:
		return PostsolveError
	case C.kHighsModelStatusModelEmpty:
		return ModelEmpty
	case C.kHighsModelStatusOptimal:
		return Optimal
	case C.kHighsModelStatusInfeasible:
		return Infeasible
	case C.kHighsModelStatusUnboundedOrInfeasible:
		return UnboundedOrInfeasible
	case C.kHighsModelStatusUnbounded:
		return Unbounded
	case C.kHighsModelStatusObjectiveBound:
		return ObjectiveBound
	case C.kHighsModelStatusObjectiveTarget:
		return ObjectiveTarget
	case C.kHighsModelStatusTimeLimit:
		return TimeLimit
	case C.kHighsModelStatusIterationLimit:
		return IterationLimit
	default:
		return UnknownModelStatus
	}
}

//go:generate stringer -type=ModelStatus

// A VariableType indicates the type of a model variable.
type VariableType int

// These are the values a VariableType accepts:
const (
	ContinuousType VariableType = iota
	IntegerType
	SemiContinuousType
	SemiIntegerType
	ImplicitIntegerType
)

// variableTypeToHighs maps a VariableType to a kHighsVarType.  This slice must
// be kept up to date with the VariableType constants.
var variableTypeToHighs = []C.HighsInt{
	C.kHighsVarTypeContinuous,
	C.kHighsVarTypeInteger,
	C.kHighsVarTypeSemiContinuous,
	C.kHighsVarTypeSemiInteger,
	C.kHighsVarTypeImplicitInteger,
}

//go:generate stringer -type=VariableType

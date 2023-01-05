/*
Package highs provides a Go interface to the [HiGHS] optimizer.  HiGHS—and the
highs package—support linear programming (LP), mixed-integer programming
(MIP) [also known as mixed-integer linear programming (MILP)], and quadratic
programming (QP) models.

The highs package provides both a low-level and a high-level interface to
HiGHS.  The low-level interface, provided by the [RawModel] type, includes a
rich set of methods that work directly with HiGHS's underlying, opaque data
type.  The high-level interface, provided by the [Model] type, is a simple
struct whose fields the programmer can modify directly.  The [Model.ToRawModel]
method converts from a Model to a RawModel, enabling models to be specified
conveniently at a high level then manipulated with the more featureful
low-level API.

In terms of the names of fields in Model, the highs package solves numerous
variants of the following core problem:

    Minimize    ColCosts ⋅ ColumnPrimal
    subject to  RowLower ≤ ConstMatrix ColumnPrimal ≤ RowUpper
    and         ColLower ≤ ColumnPrimal ≤ ColUpper

where all variables in the above are vectors except for ConstMatrix, which is a
matrix.  "ColCosts ⋅ ColumnPrimal" denotes the inner product of those two
vectors, and "ConstMatrix ColumnPrimal" denotes matrix-vector multiplication.

ColumnPrimal is a member of the [Solution] struct and is what the preceding
formulation is solving for.

[HiGHS]: https://highs.dev/
*/
package highs

// #cgo pkg-config: highs
import "C"

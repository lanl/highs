highs
=====

[![Go Reference](https://pkg.go.dev/badge/github.com/lanl/highs.svg)](https://pkg.go.dev/github.com/lanl/highs) [![Go project version](https://badge.fury.io/go/github.com%2Flanl%2Fhighs.svg)](https://badge.fury.io/go/github.com%2Flanl%2Fhighs) [![Go Report Card](https://goreportcard.com/badge/github.com/lanl/highs)](https://goreportcard.com/report/github.com/lanl/highs)

Description
-----------

The `highs` package provides a [Go](https://go.dev/) interface to the [HiGHS](https://highs.dev/) constraint-programming solver.  HiGHS—and the `highs` package—support large-scale sparse [linear programming](https://en.wikipedia.org/wiki/Linear_programming) (LP), [mixed-integer programming](https://en.wikipedia.org/wiki/Linear_programming#Integer_unknowns) (MIP), and [quadratic programming](https://en.wikipedia.org/wiki/Quadratic_programming) (QP) models.  The goal of such solvers is to minimize or maximize an expression subject to a set of constraints expressed as inequalities.  The basic form such an LP problem in HiGHS takes is as follows:

```math
\begin{array}{ll}
  \text{Find vector}    & x \\
  \text{that minimizes} & c^T x + o \\
  \text{subject to}     & b_L \leq A x \leq b_H \\
  \text{and}            & d_L \leq x \leq d_L
\end{array}
```
where "minimizes" can alternatively be "maximizes".  A MIP problem additionally constrains certain elements of $x$ to be integers, and a QP problem additionally includes an $x^T Q x$ term in the objective function.

A [detailed example of formulating a problem with `highs`](https://github.com/lanl/highs/wiki/highs-tutorial) is available on the [`highs` wiki](https://github.com/lanl/highs/wiki).

Installation
------------

`highs` has been tested only on Linux.  The package requires a HiGHS installation to build.  To check if HiGHS is installed, ensure that the following command runs without error:
```bash
pkg-config highs --cflags --libs
```
(It will typically output something like `-I/usr/include/highs -lhighs`.)

Once HiGHS installation is confirmed, the `highs` package can be installed.  From the directory of an application or package that has opted into the [Go module system](https://blog.golang.org/using-go-modules), run
```bash
go install github.com/lanl/highs
```

Documentation
-------------

See the [`highs` wiki](https://github.com/lanl/highs/wiki).

Legal statement
---------------

Copyright © 2021 Triad National Security, LLC.
All rights reserved.

This program was produced under U.S. Government contract 89233218CNA000001 for Los Alamos National Laboratory (LANL), which is operated by Triad National Security, LLC for the U.S.  Department of Energy/National Nuclear Security Administration. All rights in the program are reserved by Triad National Security, LLC, and the U.S. Department of Energy/National Nuclear Security Administration. The Government is granted for itself and others acting on its behalf a nonexclusive, paid-up, irrevocable worldwide license in this material to reproduce, prepare derivative works, distribute copies to the public, perform publicly and display publicly, and to permit others to do so.

This program is open source under the [BSD-3 License](LICENSE.md).  Its LANL-internal identifier is C21038.

Contact
-------

Scott Pakin, *pakin@lanl.gov*

/*
 * This file declares extern all of the HiGHS constants needed by the
 * highs package.  It is required because HiGHS declares all constants
 * with global scope (const) rather than with file scope (static const).
 */

#ifndef _EXTERNS_H_
#define _EXTERNS_H_

#include "util/HighsInt.h"

extern const HighsInt kHighsStatusError;
extern const HighsInt kHighsStatusOk;
extern const HighsInt kHighsStatusWarning;

extern const HighsInt kHighsVarTypeContinuous;
extern const HighsInt kHighsVarTypeInteger;
extern const HighsInt kHighsVarTypeSemiContinuous;
extern const HighsInt kHighsVarTypeSemiInteger;
extern const HighsInt kHighsVarTypeImplicitInteger;

extern const HighsInt kHighsObjSenseMinimize;
extern const HighsInt kHighsObjSenseMaximize;

extern const HighsInt kHighsMatrixFormatColwise;
extern const HighsInt kHighsMatrixFormatRowwise;

extern const HighsInt kHighsModelStatusNotset;
extern const HighsInt kHighsModelStatusLoadError;
extern const HighsInt kHighsModelStatusModelError;
extern const HighsInt kHighsModelStatusPresolveError;
extern const HighsInt kHighsModelStatusSolveError;
extern const HighsInt kHighsModelStatusPostsolveError;
extern const HighsInt kHighsModelStatusModelEmpty;
extern const HighsInt kHighsModelStatusOptimal;
extern const HighsInt kHighsModelStatusInfeasible;
extern const HighsInt kHighsModelStatusUnboundedOrInfeasible;
extern const HighsInt kHighsModelStatusUnbounded;
extern const HighsInt kHighsModelStatusObjectiveBound;
extern const HighsInt kHighsModelStatusObjectiveTarget;
extern const HighsInt kHighsModelStatusTimeLimit;
extern const HighsInt kHighsModelStatusIterationLimit;
extern const HighsInt kHighsModelStatusUnknown;

extern const HighsInt kHighsBasisStatusLower;
extern const HighsInt kHighsBasisStatusBasic;
extern const HighsInt kHighsBasisStatusUpper;
extern const HighsInt kHighsBasisStatusZero;
extern const HighsInt kHighsBasisStatusNonbasic;

extern
HighsInt Highs_lpCall(const HighsInt num_col, const HighsInt num_row,
                      const HighsInt num_nz, const HighsInt a_format,
                      const HighsInt sense, const double offset,
                      const double* col_cost, const double* col_lower,
                      const double* col_upper, const double* row_lower,
                      const double* row_upper, const HighsInt* a_start,
                      const HighsInt* a_index, const double* a_value,
                      double* col_value, double* col_dual, double* row_value,
                      double* row_dual, HighsInt* col_basis_status,
                      HighsInt* row_basis_status, HighsInt* model_status);

extern
HighsInt Highs_mipCall(const HighsInt num_col, const HighsInt num_row,
                       const HighsInt num_nz, const HighsInt a_format,
                       const HighsInt sense, const double offset,
                       const double* col_cost, const double* col_lower,
                       const double* col_upper, const double* row_lower,
                       const double* row_upper, const HighsInt* a_start,
                       const HighsInt* a_index, const double* a_value,
                       const HighsInt* integrality, double* col_value,
                       double* row_value, HighsInt* model_status);
#endif

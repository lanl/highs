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

extern const HighsInt kHighsHessianFormatTriangular;
extern const HighsInt kHighsHessianFormatSquare;

extern const HighsInt kHighsSolutionStatusNone;
extern const HighsInt kHighsSolutionStatusInfeasible;
extern const HighsInt kHighsSolutionStatusFeasible;

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
HighsInt Highs_passModel(void* highs, const HighsInt num_col,
                         const HighsInt num_row, const HighsInt num_nz,
                         const HighsInt q_num_nz, const HighsInt a_format,
                         const HighsInt q_format, const HighsInt sense,
                         const double offset, const double* col_cost,
                         const double* col_lower, const double* col_upper,
                         const double* row_lower, const double* row_upper,
                         const HighsInt* a_start, const HighsInt* a_index,
                         const double* a_value, const HighsInt* q_start,
                         const HighsInt* q_index, const double* q_value,
                         const HighsInt* integrality);

extern
HighsInt Highs_getIntInfoValue(const void* highs, const char* info,
                               HighsInt* value);

extern
HighsInt Highs_getDoubleInfoValue(const void* highs, const char* info,
                                  double* value);

extern
HighsInt Highs_getInt64InfoValue(const void* highs, const char* info,
                                 int64_t* value);

extern
HighsInt Highs_writeSolution(const void* highs, const char* filename);

extern
HighsInt Highs_writeSolutionPretty(const void* highs, const char* filename);

#endif

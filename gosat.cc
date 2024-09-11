#include "gosat.h"
#include "minisat/minisat/core/Solver.h"
#include "minisat/minisat/core/SolverTypes.h"
#include <csetjmp>
#include <csignal>
#include <iostream>
#include <vector>

struct MinisatSolver {
  Minisat::Solver solver;
};

MinisatSolver *minisatgh_new() {
  MinisatSolver *s = new MinisatSolver();
  if (!s) {
    return NULL;
  }
  return s;
}

void minisatgh_delete(MinisatSolver *solver) {
  if (solver) {
    delete solver;
  }
}

int minisatgh_add_cl(MinisatSolver *solver, int *clause, int length) {
  if (!solver)
    return 0;

  Minisat::vec<Minisat::Lit> cl;
  int max_var = -1;

  for (int i = 0; i < length; ++i) {
    int lit = clause[i];
    cl.push((lit > 0) ? Minisat::mkLit(lit - 1, false)
                      : Minisat::mkLit(-lit - 1, true));
    if (abs(lit) > max_var)
      max_var = abs(lit);
  }

  if (max_var > solver->solver.nVars()) {
    for (int i = solver->solver.nVars(); i <= max_var; ++i) {
      solver->solver.newVar();
    }
  }

  return solver->solver.addClause(cl) ? 1 : 0;
}

int minisatgh_solve(MinisatSolver *solver) {
  if (!solver)
    return 0;

  // Solve the problem with the previously added clauses
  bool res = solver->solver.solve();

  return res ? 1 : 0;
}

int minisatgh_model_size(MinisatSolver *solver) {
  if (!solver)
    return 0;
  return solver->solver.model.size();
}

int *minisatgh_model(MinisatSolver *solver) {
  if (!solver)
    return nullptr;

  Minisat::vec<Minisat::lbool> &model = solver->solver.model;

  if (model.size() == 0)
    return nullptr; // No model available

  int *c_model = (int *)malloc(sizeof(int) * model.size());

  // Minisat's l_True is represented by lbool((uint8_t)0)
  Minisat::lbool True = Minisat::lbool((uint8_t)0);

  for (int i = 0; i < model.size(); ++i) {
    c_model[i] = (model[i] == True) ? (i + 1) : -(i + 1);
  }

  return c_model;
}

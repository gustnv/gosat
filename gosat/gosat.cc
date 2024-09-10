#include "gosat.h"
#include <iostream>
#include <minisat/core/Solver.h>
#include <minisat/core/SolverTypes.h>

extern "C" {

// Create a new Minisat Solver
WrapSolver NewSolver(double seed) {
  Minisat::Solver *s = new Minisat::Solver();
  if (s == NULL) {
    return NULL; // Return NULL if solver creation fails
  }
  return (WrapSolver)s;
}

// Declare new variables
void minisat_declare_vars(void *s, const int max_id) {
  Minisat::Solver *solver = (Minisat::Solver *)s;
  while (solver->nVars() < max_id + 1)
    solver->newVar();
}

// Translate iterable to vec<Lit>
bool minisat_iterate(void *obj, void *v, int *max_var) {
  int *lits = (int *)obj;
  Minisat::vec<Minisat::Lit> *vec_v = (Minisat::vec<Minisat::Lit> *)v;

  int len = *max_var; // Assume len is passed as max_var temporarily

  for (int i = 0; i < len; i++) {
    int l = lits[i];
    if (l == 0) {
      return false; // Return false if zero is encountered
    }
    vec_v->push((l > 0) ? Minisat::mkLit(l, false) : Minisat::mkLit(-l, true));
    if (abs(l) > *max_var) {
      *max_var = abs(l);
    }
  }

  return true;
}

// Set start mode
// void minisat_set_start(WrapSolver s_obj, int warm_start) {
//   Minisat::Solver *s = (Minisat::Solver *)s_obj;
//   s->setStartMode((bool)warm_start);
// }

// Add clause to the solver
bool minisat_add_cl(WrapSolver s_obj, void *c_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  Minisat::vec<Minisat::Lit> cl;
  int max_var = -1;

  if (!minisat_iterate(c_obj, cl, &max_var))
    return false;

  if (max_var > 0)
    minisat_declare_vars(s, max_var);

  std::cout << "ok" << std::endl;
  return s->addClause(cl);
}

// Solve the SAT problem
bool minisat_solve(WrapSolver s_obj, void *a_obj, int main_thread) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  Minisat::vec<Minisat::Lit> a;
  int max_var = -1;

  if (!minisat_iterate(a_obj, a, &max_var))
    return false;

  if (max_var > 0)
    minisat_declare_vars(s, max_var);

  return s->solve(a);
}

// Clear interrupt flags
void minisat_clearint(WrapSolver s_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  s->clearInterrupt();
}

// Get the core conflict set
void *minisat_core(WrapSolver s_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  Minisat::LSet *c = &(s->conflict);

  int *core = (int *)malloc(c->size() * sizeof(int));
  for (int i = 0; i < c->size(); ++i) {
    int l = Minisat::var((*c)[i]) * (Minisat::sign((*c)[i]) ? 1 : -1);
    core[i] = l;
  }

  return (void *)core;
}

// Get the model after solving
void *minisat_model(WrapSolver s_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  Minisat::vec<Minisat::lbool> *m = &(s->model);

  int *model = (int *)malloc((m->size() - 1) * sizeof(int));
  for (int i = 1; i < m->size(); ++i) {
    model[i - 1] = i * ((*m)[i] == Minisat::lbool((uint8_t)0) ? 1 : -1);
  }

  return (void *)model;
}

// Get the number of variables
int minisat_nof_vars(WrapSolver s_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  return s->nVars() - 1;
}

// Get the number of clauses
int minisat_nof_cls(WrapSolver s_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  return s->nClauses();
}

// Delete the solver instance
void minisat_del(WrapSolver s_obj) {
  Minisat::Solver *s = (Minisat::Solver *)s_obj;
  delete s;
}
}

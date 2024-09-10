#ifndef GOSAT_H
#define GOSAT_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h> // For boolean types
#include <stdlib.h>  // For NULL

// Define the solver wrapper type
typedef void *WrapSolver;

// Function prototypes

// Create a new Minisat Solver
WrapSolver NewSolver(double seed);

// Declare new variables in the solver
void minisat_declare_vars(void *s, const int max_id);

// Translate an iterable to Minisat literals
bool minisat_iterate(void *obj, void *v, int *max_var);

// Set the start mode of the solver
// void minisat_set_start(WrapSolver s_obj, int warm_start);

// Add a clause to the solver
bool minisat_add_cl(WrapSolver s_obj, void *c_obj);

// Solve the SAT problem
bool minisat_solve(WrapSolver s_obj, void *a_obj, int main_thread);

// Clear interrupt flags in the solver
void minisat_clearint(WrapSolver s_obj);

// Get the core conflict set
void *minisat_core(WrapSolver s_obj);

// Get the model after solving
void *minisat_model(WrapSolver s_obj);

// Get the number of variables in the solver
int minisat_nof_vars(WrapSolver s_obj);

// Get the number of clauses in the solver
int minisat_nof_cls(WrapSolver s_obj);

// Delete the solver instance
void minisat_del(WrapSolver s_obj);

#ifdef __cplusplus
}
#endif

#endif // GOSAT_H

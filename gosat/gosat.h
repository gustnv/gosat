#ifndef GOSAT_H
#define GOSAT_H

#ifdef __cplusplus
extern "C" {
#endif

// Forward declaration of the MinisatSolver struct
typedef struct MinisatSolver MinisatSolver;

// Function declarations using 'struct MinisatSolver'
struct MinisatSolver *minisatgh_new();
void minisatgh_delete(struct MinisatSolver *solver);
int minisatgh_add_cl(struct MinisatSolver *solver, int *clause, int length);
int minisatgh_solve(MinisatSolver *solver);
int minisatgh_model_size(MinisatSolver *solver);
int *minisatgh_model(MinisatSolver *solver);

#ifdef __cplusplus
}
#endif

#endif // GOSAT_H

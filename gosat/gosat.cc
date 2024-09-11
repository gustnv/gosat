#include "gosat.h"
#include "minisat/core/Solver.h"
#include "minisat/core/SolverTypes.h"

struct MinisatSolver {
    Minisat::Solver solver;
};

MinisatSolver* minisatgh_new() {
    MinisatSolver* s = new MinisatSolver();
    if (!s) {
        return NULL;
    }
    return s;
}

void minisatgh_delete(MinisatSolver* solver) {
    if (solver) {
        delete solver;
    }
}

int minisatgh_add_cl(MinisatSolver* solver, int* clause, int length) {
    if (!solver) return 0;

    Minisat::vec<Minisat::Lit> cl;
    int max_var = -1;

    for (int i = 0; i < length; ++i) {
        int lit = clause[i];
        cl.push((lit > 0) ? Minisat::mkLit(lit - 1, false) : Minisat::mkLit(-lit - 1, true));
        if (abs(lit) > max_var) max_var = abs(lit);
    }

    if (max_var > solver->solver.nVars()) {
        for (int i = solver->solver.nVars(); i <= max_var; ++i) {
            solver->solver.newVar();
        }
    }

    return solver->solver.addClause(cl) ? 1 : 0;
}

int minisatgh_solve(MinisatSolver* solver, int* assumptions, int length) {
    if (!solver) return 0;

    Minisat::vec<Minisat::Lit> assumps;
    int max_var = -1;

    for (int i = 0; i < length; ++i) {
        int lit = assumptions[i];
        assumps.push((lit > 0) ? Minisat::mkLit(lit - 1, false) : Minisat::mkLit(-lit - 1, true));
        if (abs(lit) > max_var) max_var = abs(lit);
    }

    if (max_var > solver->solver.nVars()) {
        for (int i = solver->solver.nVars(); i <= max_var; ++i) {
            solver->solver.newVar();
        }
    }

    return solver->solver.solve(assumps) ? 1 : 0;
}

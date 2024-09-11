package gosat

/*
#cgo LDFLAGS: -lminisat
#include "gosat.h" // Header file for interfacing with C/C++ code
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

type MinisatGH struct {
	minisat *C.MinisatSolver
	Status  bool // Exported field
}

func NewMinisatGH(bootstrapWith [][]int, useTimer bool, warmStart bool) (*MinisatGH, error) {
	solver := C.minisatgh_new()
	if solver == nil {
		return nil, errors.New("cannot create a new solver")
	}

	m := &MinisatGH{minisat: solver, Status: true} // Initialize Status as true

	if bootstrapWith != nil {
		for _, clause := range bootstrapWith {
			if err := m.AddClause(clause, true); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func (m *MinisatGH) Delete() {
	if m.minisat != nil {
		C.minisatgh_delete(m.minisat)
		m.minisat = nil
	}
}

func (m *MinisatGH) AddClause(clause []int, noReturn bool) error {
	if m.minisat == nil {
		return errors.New("solver is not initialized")
	}

	cClause := (*C.int)(C.malloc(C.size_t(len(clause)) * C.size_t(unsafe.Sizeof(C.int(0)))))
	defer C.free(unsafe.Pointer(cClause))

	slice := (*[1 << 30]C.int)(unsafe.Pointer(cClause))[:len(clause):len(clause)]
	for i, lit := range clause {
		slice[i] = C.int(lit)
	}

	res := C.minisatgh_add_cl(m.minisat, (*C.int)(cClause), C.int(len(clause)))
	if res == 0 {
		m.Status = false // Update Status on failure
	}

	if !noReturn && res == 0 {
		return errors.New("failed to add clause")
	}

	return nil
}

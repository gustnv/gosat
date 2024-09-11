package gosat

/*
#cgo CXXFLAGS: -I./minisat/minisat
#cgo LDFLAGS: -L./minisat -lminisat
#include <stdlib.h>
#include "gosat.h"
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

// Solve function added to call the underlying C function
func (m *MinisatGH) Solve() (bool, error) {
	if m.minisat == nil {
		return false, errors.New("solver is not initialized")
	}

	// Call the C++ function to solve the problem
	res := C.minisatgh_solve(m.minisat)

	if res == 0 {
		return false, nil
	}

	return true, nil
}

func (m *MinisatGH) GetModel() ([]int, error) {
	if m.minisat == nil {
		return nil, errors.New("solver is not initialized")
	}

	// Call the C function to get the size of the model
	modelSize := C.minisatgh_model_size(m.minisat)
	if modelSize == 0 {
		return nil, errors.New("no model available")
	}

	// Call the C function to get the model
	cModel := C.minisatgh_model(m.minisat)
	if cModel == nil {
		return nil, errors.New("no model available")
	}
	defer C.free(unsafe.Pointer(cModel))

	// Copy the model from the C array to a Go slice
	model := make([]int, int(modelSize))
	slice := (*[1 << 30]C.int)(unsafe.Pointer(cModel))[:modelSize:modelSize]
	for i := range model {
		model[i] = int(slice[i])
	}

	return model, nil
}

func (m *MinisatGH) AppendFormula(formula [][]int) error {
	if m.minisat == nil {
		return errors.New("solver is not initialized")
	}

	for _, clause := range formula {
		err := m.AddClause(clause, true)

		if err != nil {
			m.Status = false // Update the Status field on failure
			return err
		}
	}

	return nil
}

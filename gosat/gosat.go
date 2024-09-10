package gosat

/*
#cgo LDFLAGS: -lminisat
#include "gosat.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

// Minisat is a wrapper around the MiniSat SAT solver.
type Minisat struct {
	minisat  C.WrapSolver
	status   bool
	useTimer bool
	callTime float64
	accuTime float64
}

// NewMinisat creates a new Minisat solver instance.
func NewMinisat(bootstrapWith [][]int, useTimer bool, warmStart bool) (*Minisat, error) {
	solver := &Minisat{}
	solver.minisat = C.NewSolver(0) // Default seed value

	if solver.minisat == nil {
		return nil, errors.New("failed to create Minisat solver instance")
	}

	if len(bootstrapWith) > 0 {
		for _, clause := range bootstrapWith {
			if err := solver.AddClause(clause); err != nil {
				return nil, err
			}
		}
	}

	if warmStart {
		solver.StartMode(true)
	}

	solver.useTimer = useTimer
	solver.callTime = 0.0
	solver.accuTime = 0.0
	return solver, nil
}

// Delete cleans up the Minisat solver instance.
func (s *Minisat) Delete() {
	if s.minisat != nil {
		C.minisat_del(s.minisat)
		s.minisat = nil
	}
}

// Solve attempts to solve the SAT problem.
func (s *Minisat) Solve(assumptions []int) (bool, error) {
	if s.minisat == nil {
		return false, errors.New("solver instance is not initialized")
	}

	if s.useTimer {
		start := time.Now()
		defer func() {
			s.callTime = time.Since(start).Seconds()
			s.accuTime += s.callTime
		}()
	}

	assumptionsC := (*C.int)(C.malloc(C.size_t(len(assumptions)) * C.sizeof_int))
	defer C.free(unsafe.Pointer(assumptionsC))

	for i, a := range assumptions {
		*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(assumptionsC)) + uintptr(i)*C.sizeof_int)) = C.int(a)
	}

	result := C.minisat_solve(s.minisat, unsafe.Pointer(assumptionsC), C.int(1))
	s.status = bool(result)

	return s.status, nil
}

// AddClause adds a clause to the solver.
func (s *Minisat) AddClause(clause []int) error {
	if s.minisat == nil {
		return errors.New("solver instance is not initialized")
	}

	clauseC := (*C.int)(C.malloc(C.size_t(len(clause)) * C.sizeof_int))
	defer C.free(unsafe.Pointer(clauseC))

	for i, l := range clause {
		*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(clauseC)) + uintptr(i)*C.sizeof_int)) = C.int(l)
	}

	res := C.minisat_add_cl(s.minisat, unsafe.Pointer(clauseC))
	if !bool(res) {
		s.status = false
		return errors.New("failed to add clause")
	}
	return nil
}

// StartMode sets the solver's start mode.
func (s *Minisat) StartMode(warm bool) {
	if s.minisat != nil {
		C.minisat_set_start(s.minisat, C.int(boolToInt(warm)))
	}
}

// GetModel returns the model found after solving.
func (s *Minisat) GetModel() ([]int, error) {
	if s.minisat == nil || !s.status {
		return nil, errors.New("no model available")
	}

	modelPtr := C.minisat_model(s.minisat)
	if modelPtr == nil {
		return nil, nil
	}
	defer C.free(unsafe.Pointer(modelPtr))

	numVars := s.NofVars()
	model := make([]int, numVars)

	for i := 0; i < numVars; i++ {
		model[i] = int(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(modelPtr)) + uintptr(i)*C.sizeof_int)))
	}
	return model, nil
}

// GetCore returns the core conflict set.
func (s *Minisat) GetCore() ([]int, error) {
	if s.minisat == nil || s.status {
		return nil, errors.New("no core available")
	}

	corePtr := C.minisat_core(s.minisat)
	if corePtr == nil {
		return nil, nil
	}
	defer C.free(unsafe.Pointer(corePtr))

	numVars := s.NofVars()
	core := make([]int, numVars)

	for i := 0; i < numVars; i++ {
		core[i] = int(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(corePtr)) + uintptr(i)*C.sizeof_int)))
	}
	return core, nil
}

// NofVars returns the number of variables in the solver.
func (s *Minisat) NofVars() int {
	if s.minisat != nil {
		return int(C.minisat_nof_vars(s.minisat))
	}
	return 0
}

// NofClauses returns the number of clauses in the solver.
func (s *Minisat) NofClauses() int {
	if s.minisat != nil {
		return int(C.minisat_nof_cls(s.minisat))
	}
	return 0
}

// ClearInterrupt clears any pending interrupts.
func (s *Minisat) ClearInterrupt() {
	if s.minisat != nil {
		C.minisat_clearint(s.minisat)
	}
}

// Interrupt sends an interrupt signal to the solver.
func (s *Minisat) Interrupt() {
	if s.minisat != nil {
		C.minisat_clearint(s.minisat)
	}
}

// Utility function to convert a Go boolean to a C int.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

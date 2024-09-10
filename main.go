package main

import (
	"log"

	"github.com/gustnv/gosat/gosat"
)

func main() {
	// Create a new Minisat solver instance
	solver, err := gosat.NewMinisat(nil, false, false)
	if err != nil {
		log.Fatalf("failed to create Minisat solver instance: %v", err)
	}
	defer solver.Delete()

	// // Add a clause to the solver
	// clause := []int{1, 2, 3}
	// if err := solver.AddClause(clause); err != nil {
	// 	log.Fatalf("failed to add clause to solver: %v", err)
	// }
	//
	// // Solve the SAT problem
	// status, err := solver.Solve(nil)
	// if err != nil {
	// 	log.Fatalf("failed to solve SAT problem: %v", err)
	// }
	//
	// // Print the result
	// if status {
	// 	log.Println("SAT")
	// } else {
	// 	log.Println("UNSAT")
	// }
}

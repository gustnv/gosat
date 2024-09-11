package main

import (
	"fmt"
	"github.com/gustnv/gosat/gosat"
)

func main() {
	// Create a new MinisatGH instance
	minisatInstance, err := gosat.NewMinisatGH(nil, false, false)
	if err != nil {
		fmt.Println("Failed to create MinisatGH instance:", err)
		return
	}
	defer minisatInstance.Delete()

	// Attempt to add a clause
	clause := []int{1, 2, 3}
	err = minisatInstance.AddClause(clause, false)
	if err != nil {
		fmt.Println("Failed to add clause:", err)
	} else {
		fmt.Println("Clause added successfully.")
	}

	clause = []int{-1, -2, -3}
	err = minisatInstance.AddClause(clause, false)
	if err != nil {
		fmt.Println("Failed to add clause:", err)
	} else {
		fmt.Println("Clause added successfully.")
	}

	// Check the status of the solver
	if minisatInstance.Status {
		fmt.Println("Solver status: OK")
	} else {
		fmt.Println("Solver status: Failed to add clause.")
	}
}

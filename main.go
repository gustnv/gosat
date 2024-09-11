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
	clause := []int{1, -2}
	err = minisatInstance.AddClause(clause, false)
	clause = []int{2, -3}
	err = minisatInstance.AddClause(clause, false)
	clause = []int{3, -1}
	err = minisatInstance.AddClause(clause, false)
	clause = []int{1, 2, 3}
	err = minisatInstance.AddClause(clause, false)
	clause = []int{-1, -2, -3}
	err = minisatInstance.AddClause(clause, false)

	// Solve the formula
	assumptions := []int{}
	result, err := minisatInstance.Solve(assumptions)
	if err != nil {
		fmt.Println("Failed to solve formula:", err)
	} else {
		fmt.Println("Formula solved successfully. Result:", result)
	}

}

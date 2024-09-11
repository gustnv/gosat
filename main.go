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

	minisatInstance.AddClause([]int{1, -2}, false)
	minisatInstance.AddClause([]int{2, -3}, false)
	minisatInstance.AddClause([]int{3, -1}, false)
	minisatInstance.AddClause([]int{1, 2, 3}, false)
	// minisatInstance.AddClause([]int{-1, -2, -3}, false)

	result, err := minisatInstance.Solve()
	if err != nil {
		fmt.Println("Failed to solve formula:", err)
	} else {
		fmt.Println("Formula solved successfully. Result:", result)
	}

	r, err := minisatInstance.GetModel()
	if err != nil {
		fmt.Println("Failed to get model:", err)
	} else {
		fmt.Println("Model:", r)
	}

}

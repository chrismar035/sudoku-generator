package main

import (
	"fmt"

	"github.com/chrismar035/solver"
)

func main() {
	var grid solver.Grid
	puzzle := solver.Puzzle{Initial: grid}
	solver := solver.NewRandBacktrackingSolver()
	puzzle.Solution = solver.Solve(grid)
	fmt.Println(puzzle)
}

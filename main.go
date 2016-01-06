package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chrismar035/sudoku-solver"
)

type removedSquare struct {
	index int
	value int
}

func main() {
	var grid solver.Grid
	randomizer := solver.NewRandBacktrackingSolver()

	puzzle := randomizer.Solve(grid)
	solution := puzzle
	fmt.Println("Solution:")
	fmt.Println(solution)

	indexes := randomizeIndexes()
	var removed []removedSquare

	solver := solver.NewMultiBacktrackingSolver()

	for _, index := range indexes {
		fmt.Println(index)
		removed = append(removed, removedSquare{index: index, value: puzzle[index]})
		puzzle[index] = 0

		if len(solver.Solve(puzzle)) > 1 {
			last := removed[len(removed)-1]
			puzzle[last.index] = last.value

			fmt.Println("\nPuzzle:")
			fmt.Println(puzzle)
			return
		}
	}
	fmt.Println("Couldn't find puzzle")
}

func randomizeIndexes() []int {
	rand.Seed(time.Now().UTC().UnixNano())

	ints := []int{}
	for i := 0; i < 81; i++ {
		ints = append(ints, i)
	}

	mixed := []int{}
	for len(ints) > 0 {
		i := rand.Int() % len(ints)
		mixed = append(mixed, ints[i])
		ints = append(ints[0:i], ints[i+1:]...)
	}

	return mixed
}

package main

import (
	"errors"
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

	for i := 0; i < 4; i++ {
		solution := getShuffledSolution()
		puzzle, err := puzzleFromSolution(solution)
		if err != nil {
			fmt.Println(err)
			fmt.Println(solution)
		} else {
			fmt.Println(puzzle)
			fmt.Println(solution)
		}
		fmt.Println("-----------")
	}
}

func puzzleFromSolution(solution solver.Grid) (solver.Grid, error) {
	puzzle := solution
	indexes := randomizeIndexes()
	var removed []removedSquare

	multiSolver := solver.NewMultiBacktrackingSolver()

	for _, index := range indexes {
		removed = append(removed, removedSquare{index: index, value: puzzle[index]})
		puzzle[index] = 0

		if len(multiSolver.Solve(puzzle)) > 1 {
			last := removed[len(removed)-1]
			puzzle[last.index] = last.value

			return puzzle, nil
		}
	}
	return solver.Grid{}, errors.New("Couldn't find puzzle")
}

func getShuffledSolution() solver.Grid {
	var grid solver.Grid
	randomizer := solver.NewRandBacktrackingSolver()

	return randomizer.Solve(grid)
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

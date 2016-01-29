package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chrismar035/sudoku-solver"
)

var logger = log.New(os.Stdout,
	"Generator: ",
	log.Ldate|log.Ltime|log.Lshortfile)

type removedSquare struct {
	index int
	value int
}

type postParams struct {
	Puzzle   solver.Grid `json:"puzzle"`
	Solution solver.Grid `json:"solution"`
}

type Sudoku struct {
	Id       string  `json:"id"`
	Puzzle   [81]int `json:"puzzle"`
	Solution [81]int `json:"solution"`
	Name     string  `json:"name"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

func main() {
	url := os.Getenv("API_ROOT")

	logger.Println("Starting loop")
	for {
		solution := getShuffledSolution()
		puzzle, err := puzzleFromSolution(solution)
		if err != nil {
			logger.Println("Error generating puzzle", solution, err)
		} else {
			params := postParams{Puzzle: puzzle, Solution: solution}
			jsonStr, err := json.Marshal(params)
			if err != nil {
				logger.Println("Unable to marshal puzzle", puzzle, solution)
				continue
			}

			req, err := http.NewRequest("POST", url+"/puzzle", bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Println("Unable to submit puzzle", err)
				errorResponse := errorFromBody(resp)
				postToSlack("{ \"text\": \"Found Duplicate: " + errorResponse.Id + "\" }")
				continue
			} else {
				sudoku := sudokuFromBody(resp)
				logger.Println("Added:", sudoku.Id)
				// postToSlack("{ \"text\": \"Added: " + sudoku.Id + "\" }")
			}
		}
		logger.Println("Iteration")
	}
	logger.Println("Out of loop. Ending.")
}

func postToSlack(message string) {
	url := "https://hooks.slack.com/services/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"
	req, _ := http.NewRequest("POST", url, strings.NewReader(message))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println("Failed to post to slack:", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logger.Println("Error reading slack response:", err)
		return
	}

	logger.Println("Slack response:", string(body))
}

func puzzleFromSolution(solution solver.Grid) (solver.Grid, error) {
	puzzle := solution
	indexes := randomizeIndexes()
	var removed []removedSquare

	singleSolver := solver.NewSingleBacktrackingSolver()

	for _, index := range indexes {
		removed = append(removed, removedSquare{index: index, value: puzzle[index]})
		puzzle[index] = 0

		_, err := singleSolver.Solve(puzzle)
		if err != nil && err.Error() == "Multiple solutions found" {
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

	random, _ := randomizer.Solve(grid)
	return random
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

func sudokuFromBody(r *http.Response) Sudoku {
	var sudoku Sudoku
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Println("Error reading sudoku:", err)
		return Sudoku{}
	}
	if err := r.Body.Close(); err != nil {
		logger.Println("Error closing body:", err)
		return Sudoku{}
	}
	logger.Println("sudoku body:", body)
	if err := json.Unmarshal(body, &sudoku); err != nil {
		logger.Println("Error unmarshaling sudoku:", err)
		return Sudoku{}
	}
	return sudoku
}

func errorFromBody(r *http.Response) ErrorResponse {
	var response ErrorResponse
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Println("Error reading error:", err)
		return ErrorResponse{}
	}
	if err := r.Body.Close(); err != nil {
		logger.Println("Error closing body:", err)
		return ErrorResponse{}
	}
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Println("Error unmarshaling error:", err)
		return ErrorResponse{}
	}
	return response
}

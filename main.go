package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	count = 0
)

type Board [][]int

func NewBoard(input []string) (*Board, error) {
	if len(input) != 9 {
		return nil, fmt.Errorf("invalid input, please input a 9x9 sudoku Board")
	}

	var board Board = make([][]int, 9)
	for i := 0; i < 9; i++ {
		line := input[i]
		if len([]rune(line)) != 9 {
			return nil, fmt.Errorf("invalid length at line %v, please input a 9x9 sudoku Board", i+1)
		}

		board[i] = make([]int, 9)
		for j, r := range []rune(line) {
			if r == '_' {
				board[i][j] = 0
			} else {
				num, err := strconv.Atoi(string(r))

				if err != nil {
					panic(err)
				}
				if num < 1 || 9 < num {
					return nil, fmt.Errorf("invalid input, please input a number between 1 ~ 9 or '_' for blank")
				}

				board[i][j] = num
			}
		}
	}

	return &board, nil
}

// Prints the board, printing '_' for blanks.
func (b *Board) String() string {
	ret := ""
	for i := 0; i < 9; i++ {
		line := ""
		for j := 0; j < 9; j++ {
			num := (*b)[i][j]
			if num == 0 {
				line += "_"
			} else {
				line += strconv.Itoa(num)
			}
		}
		ret += line + "\n"
	}
	return ret
}

// Solves the sudoku, returns true if successfully solved the game.
func (b *Board) Solve(depth int) bool {
	possibilitiesMap := make([][][]int, 9)

	for i := 0; i < 9; i++ {
		possibilitiesMap[i] = make([][]int, 9)
		for j := 0; j < 9; j++ {
			if (*b)[i][j] != 0 {
				continue
			}

			// check possible numbers for each unfilled square
			possibilitiesMap[i][j] = b.possibleNumbersAt(i, j)
		}
	}

	count++
	if count%100_000 == 0 {
		fmt.Printf("Checked %v-th time, depth %v\n", count, depth)
		fmt.Println(b.String())
	}

	// check if any square is not fillable
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if (*b)[i][j] != 0 {
				continue
			}

			possibilities := possibilitiesMap[i][j]
			if len(possibilities) == 0 {
				return false
			}
		}
	}

	// fill if there's only one possibility
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if (*b)[i][j] != 0 {
				continue
			}

			possibilities := possibilitiesMap[i][j]
			if len(possibilities) == 1 {
				(*b)[i][j] = possibilities[0]
				// fill the number and check
				if b.Solve(depth + 1) {
					return true
				}
				(*b)[i][j] = 0
			}
		}
	}

	// if multiple numbers are possible, check each of them one by one
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if (*b)[i][j] != 0 {
				continue
			}

			possibilities := possibilitiesMap[i][j]
			if len(possibilities) > 1 {
				for _, num := range possibilities {
					(*b)[i][j] = num
					// assume the number and check
					if b.Solve(depth + 1) {
						return true
					}
				}
				// all combinations failed
				(*b)[i][j] = 0
				return false
			}
		}
	}

	// if none of above code returned true of false, then all numbers must have been filled
	return b.isSolved()
}

// Checks the entire board is filled and is valid
func (b *Board) isSolved() bool {
	// check if all is filled
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if (*b)[i][j] == 0 {
				return false
			}
		}
	}
	// check blocks
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !b.checkBlockValidity(i, j) {
				return false
			}
		}
	}
	// check rows
	for i := 0; i < 9; i++ {
		if !b.checkRowValidity(i) {
			return false
		}
	}
	// check columns
	for j := 0; j < 9; j++ {
		if !b.checkColumnValidity(j) {
			return false
		}
	}
	return true
}

func (b *Board) possibleNumbersAt(i, j int) []int {
	old := (*b)[i][j]

	taken := make([]bool, 9)
	// row
	for _, num := range (*b)[i] {
		if num == 0 {
			continue
		}
		taken[num-1] = true
	}
	// column
	for i := 0; i < 9; i++ {
		num := (*b)[i][j]
		if num == 0 {
			continue
		}

		taken[num-1] = true
	}
	// block
	blockI := i / 3
	blockJ := j / 3

	for i := blockI * 3; i < (blockI+1)*3; i++ {
		for j := blockJ * 3; j < (blockJ+1)*3; j++ {
			num := (*b)[i][j]
			if num == 0 {
				continue
			}

			taken[num-1] = true
		}
	}

	(*b)[i][j] = old

	possible := make([]int, 0)
	for i, b := range taken {
		if !b {
			possible = append(possible, i+1)
		}
	}

	return possible
}

// Checks if the board is valid regarding number at i, j.
func (b *Board) isValidAt(i, j int) bool {
	return b.checkRowValidity(i) && b.checkColumnValidity(j) && b.checkBlockValidity(i, j)
}

// Checks row validity
func (b *Board) checkRowValidity(i int) bool {
	checked := make([]bool, 9)
	for _, num := range (*b)[i] {
		if num == 0 {
			continue
		}

		if checked[num-1] {
			return false
		}
		checked[num-1] = true
	}
	return true
}

// Checks column validity
func (b *Board) checkColumnValidity(j int) bool {
	checked := make([]bool, 9)
	for i := 0; i < 9; i++ {
		num := (*b)[i][j]
		if num == 0 {
			continue
		}

		if checked[num-1] {
			return false
		}
		checked[num-1] = true
	}
	return true
}

// Checks block validity of the given coords (NOT block-wise coords)
func (b *Board) checkBlockValidity(i, j int) bool {
	blockI := i / 3
	blockJ := j / 3

	checked := make([]bool, 9)
	for i := blockI * 3; i < (blockI+1)*3; i++ {
		for j := blockJ * 3; j < (blockJ+1)*3; j++ {
			num := (*b)[i][j]
			if num == 0 {
				continue
			}

			if checked[num-1] {
				return false
			}
			checked[num-1] = true
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input a 9x9 sudoku Board:")
	input := make([]string, 0, 9)
	for i := 0; i < 9; i++ {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			panic(err)
		}
		input = append(input, line)
	}

	board, err := NewBoard(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Solving...")
	start := time.Now().UnixNano()
	if board.Solve(0) {
		fmt.Println("Solved!")
	} else {
		fmt.Println("Not solvable...")
	}
	end := time.Now().UnixNano()
	fmt.Println(board.String())
	fmt.Printf("Took %v ms to solve.\n", float64(end-start)/float64(1_000_000))
}

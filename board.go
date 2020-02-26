package main

import (
	"fmt"
	"strconv"
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

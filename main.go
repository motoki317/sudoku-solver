package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

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
	if solvedBoard := board.SolveConcurrently(); solvedBoard != nil {
		fmt.Println("Solved!")
		fmt.Println(solvedBoard.String())
	} else {
		fmt.Println("Not solvable...")
	}
	end := time.Now().UnixNano()
	fmt.Printf("Took %v ms to solve.\n", float64(end-start)/float64(1_000_000))
}

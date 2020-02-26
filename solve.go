package main

// Solves the sudoku using goroutines, returns the board if successfully solved the game.
func (b *Board) Solve() *Board {
	possibilitiesMap := make([][][]int, 9)

	for i := 0; i < 9; i++ {
		possibilitiesMap[i] = make([][]int, 9)
		for j := 0; j < 9; j++ {
			if (*b)[i][j] != 0 {
				continue
			}

			// check possible numbers for each unfilled square
			// if any of them were not fillable, then return immediately
			if possibilities := b.possibleNumbersAt(i, j); len(possibilities) == 0 {
				return nil
			} else {
				possibilitiesMap[i][j] = possibilities
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
				if solvedBoard := b.Solve(); solvedBoard != nil {
					return solvedBoard
				}
				(*b)[i][j] = 0
				return nil
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
					if solvedBoard := b.Solve(); solvedBoard != nil {
						return solvedBoard
					}
				}
				(*b)[i][j] = 0
				return nil
			}
		}
	}

	// if none of above code returned true of false, then all numbers must have been filled
	if b.isSolved() {
		return b
	} else {
		return nil
	}
}

func (b *Board) clone() *Board {
	var newBoard Board = make([][]int, 9)

	for i := 0; i < 9; i++ {
		newBoard[i] = make([]int, 9)
		for j := 0; j < 9; j++ {
			newBoard[i][j] = (*b)[i][j]
		}
	}

	return &newBoard
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
			if !b.checkBlockValidity(i*3, j*3) {
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

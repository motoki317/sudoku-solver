package main

// Solves the sudoku using goroutines, returns the board if successfully solved the game.
func (b *Board) SolveConcurrently() *Board {
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
				if solvedBoard := b.SolveConcurrently(); solvedBoard != nil {
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
				solved := make(chan *Board, 0)

				// check each possibility concurrently
				for _, num := range possibilities {
					finalNum := num
					go func() {
						cloned := b.clone()
						(*cloned)[i][j] = finalNum
						// switch to normal (single-threaded) solving,
						// this might be faster than recursively calling SolveConcurrently() method
						solved <- cloned.Solve()
					}()
				}

				for range possibilities {
					if solvedBoard := <-solved; solvedBoard != nil {
						return solvedBoard
					}
				}
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

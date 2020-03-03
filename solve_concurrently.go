package main

// Solves the sudoku using goroutines, returns the board if successfully solved the game.
func (b *Board) SolveConcurrently() *Board {
	iMin := -1
	jMin := -1
	min := 10
	var possibilities []int

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if (*b)[i][j] != 0 {
				continue
			}

			// check possible numbers for each unfilled square
			// if any of them were not fillable, then return immediately
			if p := b.possibleNumbersAt(i, j); len(p) < min {
				min = len(p)
				possibilities = p
				iMin = i
				jMin = j
			}
		}
	}

	if possibilities == nil {
		// if all numbers have been filled
		if b.isSolved() {
			return b
		} else {
			return nil
		}
	}

	if len(possibilities) == 1 {
		(*b)[iMin][jMin] = possibilities[0]
		// fill the number and check
		if solvedBoard := b.SolveConcurrently(); solvedBoard != nil {
			return solvedBoard
		}
		(*b)[iMin][jMin] = 0
		return nil
	} else {
		solved := make(chan *Board, 0)

		// check each possibility concurrently
		for _, num := range possibilities {
			finalNum := num
			go func() {
				cloned := b.clone()
				(*cloned)[iMin][jMin] = finalNum
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

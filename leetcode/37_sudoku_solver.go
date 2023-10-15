package main

import (
	"fmt"
)

func solveSudoku(board [][]byte) {
	availableInRow := make([]uint16, 9)
	availableInColumn := make([]uint16, 9)
	availableInMatrix := make([]uint16, 9)

	for i := 0; i < 9; i++ {
		availableInRow[i] = (1 << 9) - 1
		availableInColumn[i] = (1 << 9) - 1
		availableInMatrix[i] = (1 << 9) - 1
	}
	initializeTrackers(board, availableInRow, availableInColumn, availableInMatrix)
	recurseSolve(board, availableInRow, availableInColumn, availableInMatrix)
}

func initializeTrackers(board [][]byte, availableInRow []uint16, availableInColumn []uint16, availableInMatrix []uint16) {
	for i, row := range board {
		for j, val := range row {
			if val != '.' {
				pos := (val - '0') - 1
				clearAvailability(uint8(i), uint8(j), pos, availableInRow, availableInColumn, availableInMatrix)
			}
		}
	}
}

func setBit(val *uint16, pos uint8) {
	*val |= 1 << pos
}

func clearBit(val *uint16, pos uint8) {
	*val &= ^(1 << pos)
}

func bitSet(val uint16, pos uint8) bool {
	return (val & (1 << pos)) != 0
}

func subMatrixIndex(i uint8, j uint8) uint8 {
	return ((i / 3) * 3) + (j / 3)
}

func clearAvailability(i, j, pos uint8, availableInRow []uint16, availableInColumn []uint16, availableInMatrix []uint16) {
	clearBit(&availableInRow[i], pos)
	clearBit(&availableInColumn[j], pos)
	clearBit(&availableInMatrix[subMatrixIndex(uint8(i), uint8(j))], pos)
}

func setAvailability(i, j, pos uint8, availableInRow []uint16, availableInColumn []uint16, availableInMatrix []uint16) {
	setBit(&availableInRow[i], pos)
	setBit(&availableInColumn[j], pos)
	setBit(&availableInMatrix[subMatrixIndex(i, j)], pos)
}

type optionsTracker = struct {
	I              uint8
	J              uint8
	Candidates     uint16
	CandidateCount uint8
}

func recurseSolve(board [][]byte, availableInRow []uint16, availableInColumn []uint16, availableInMatrix []uint16) bool {
	options := optionsTracker{}
	options.CandidateCount = 10
	done := true
	for submatrixIdx := uint8(0); submatrixIdx < 9; submatrixIdx++ {
		ok, subMatrixDone, submatrixOptions := getMostSpecificOption(board, submatrixIdx, availableInRow, availableInColumn, availableInMatrix)
		if !ok {
			return false
		}
		if !subMatrixDone && submatrixOptions.CandidateCount < options.CandidateCount {
			options = submatrixOptions
		}
		done = done && subMatrixDone
	}
	if done {
		return true
	}

	for pos := uint8(0); pos < 9; pos++ {
		if !bitSet(options.Candidates, pos) {
			continue
		}
		board[options.I][options.J] = pos + 1 + '0'

		clearAvailability(options.I, options.J, pos, availableInRow, availableInColumn, availableInMatrix)
		ok := recurseSolve(board, availableInRow, availableInColumn, availableInMatrix)
		if ok {
			return true
		}
		board[options.I][options.J] = '.'
		setAvailability(options.I, options.J, pos, availableInRow, availableInColumn, availableInMatrix)
	}
	return false
}

func getMostSpecificOption(board [][]byte, submatrixIdx uint8, availableInRow []uint16, availableInColumn []uint16, availableInMatrix []uint16) (bool, bool, optionsTracker) {
	complete := true

	mostSpecificOption := optionsTracker{}
	mostSpecificOption.CandidateCount = 10

	options := make([]optionsTracker, 0, 9)
	startI := uint8(submatrixIdx/3) * 3
	startJ := uint8(submatrixIdx%3) * 3
	for i := startI; i < startI+3; i++ {
		for j := startJ; j < startJ+3; j++ {
			if board[i][j] != '.' {
				continue
			} else {
				complete = false
			}
			subOptions, subOptionsCount := getAllOptions(i, j, availableInRow, availableInColumn, availableInMatrix)
			if subOptionsCount == 0 {
				return false, complete, mostSpecificOption
			}
			options = append(options, optionsTracker{i, j, subOptions, subOptionsCount})
			if subOptionsCount < mostSpecificOption.CandidateCount {
				mostSpecificOption = options[len(options)-1]
			}
		}
	}

	// Check all the masks for a number that can only be placed into one position in the matrix. If there is such a
	// number, return a mask for it.
BitPosition:
	for pos := uint8(0); pos < 9; pos++ {
		lastSet := -1
		for index, option := range options {
			if bitSet(option.Candidates, pos) {
				if lastSet != -1 {
					continue BitPosition
				}
				lastSet = index
			}
		}
		if lastSet != -1 {
			// set the mask to the single option
			mostSpecificOption = options[lastSet]
			mostSpecificOption.Candidates = 0
			setBit(&mostSpecificOption.Candidates, pos)
			mostSpecificOption.CandidateCount = 1
			break
		}
	}

	return true, complete, mostSpecificOption
}

func getAllOptions(i uint8, j uint8, availableInRow []uint16, availableInColumn []uint16, availableInMatrix []uint16) (uint16, uint8) {
	optionMask := availableInRow[i] & availableInColumn[j] & availableInMatrix[subMatrixIndex(i, j)]
	maskCopy := optionMask
	counter := uint8(0)
	for ; maskCopy != 0; maskCopy = maskCopy >> 1 {
		if maskCopy&1 == 1 {
			counter++
		}
	}
	return optionMask, counter
}

func dumpBoard(board [][]byte) {
	for _, row := range board {
		for _, val := range row {
			fmt.Printf("%c ", val)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	// board_one := [][]byte{
	// 	{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
	// 	{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
	// 	{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
	// 	{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
	// 	{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
	// 	{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
	// 	{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
	// 	{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
	// 	{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	// }
	// [["5","3","4","6","7","8","9","1","2"],
	// ["6","7","2","1","9","5","3","4","8"],
	// ["1","9","8","3","4","2","5","6","7"],
	// ["8","5","9","7","6","1","4","2","3"],
	// ["4","2","6","8","5","3","7","9","1"],
	// ["7","1","3","9","2","4","8","5","6"],
	// ["9","6","1","5","3","7","2","8","4"],
	// ["2","8","7","4","1","9","6","3","5"],
	// ["3","4","5","2","8","6","1","7","9"]]
	// solveSudoku(board_one)
	// dumpBoard(board_one)
	// fmt.Println()

	// board_two := [][]byte{
	// 	{'.', '.', '9', '7', '4', '8', '.', '.', '.'},
	// 	{'7', '.', '.', '.', '.', '.', '.', '.', '.'},
	// 	{'.', '2', '.', '1', '.', '9', '.', '.', '.'},
	// 	{'.', '.', '7', '.', '.', '.', '2', '4', '.'},
	// 	{'.', '6', '4', '.', '1', '.', '5', '9', '.'},
	// 	{'.', '9', '8', '.', '.', '.', '3', '.', '.'},
	// 	{'.', '.', '.', '8', '.', '3', '.', '2', '.'},
	// 	{'.', '.', '.', '.', '.', '.', '.', '.', '6'},
	// 	{'.', '.', '.', '2', '7', '5', '9', '.', '.'},
	// }
	// Expected
	// [["5","1","9","7","4","8","6","3","2"],
	// ["7","8","3","6","5","2","4","1","9"],
	// ["4","2","6","1","3","9","8","7","5"],
	// ["3","5","7","9","8","6","2","4","1"],
	// ["2","6","4","3","1","7","5","9","8"],
	// ["1","9","8","5","2","4","3","6","7"],
	// ["9","7","5","8","6","3","1","2","4"],
	// ["8","3","2","4","9","1","7","5","6"],
	// ["6","4","1","2","7","5","9","8","3"]]

	// . . 9 7 4 8 . . 2
	// 7 . . 6 . 2 . . 9  <- it would start here in the middle in bad case
	// . 2 . 1 . 9 . . .
	// . . 7 9 8 6 2 4 1  <- in good case it would guess this
	// 2 6 4 3 1 7 5 9 8
	// 1 9 8 5 2 4 3 6 7
	// 9 . . 8 6 3 . 2 .
	// . . 2 4 9 1 . . 6
	// . . . 2 7 5 9 . .

	// solveSudoku(board_two)
	// dumpBoard(board_two)

	messed_up := [][]byte{
		{'.', '.', '9', '7', '4', '8', '.', '.', '2'},
		{'7', '.', '.', '6', '.', '2', '.', '.', '9'},
		{'.', '2', '.', '1', '.', '9', '.', '.', '.'},
		{'.', '.', '7', '9', '8', '6', '2', '4', '1'},
		{'2', '6', '4', '3', '1', '7', '5', '9', '8'},
		{'1', '9', '8', '5', '2', '4', '3', '6', '7'},
		{'9', '.', '.', '8', '6', '3', '.', '2', '.'},
		{'.', '.', '2', '4', '9', '1', '.', '.', '6'},
		{'.', '.', '.', '2', '7', '5', '9', '.', '.'},
	}
	solveSudoku(messed_up)
	dumpBoard(messed_up)
}

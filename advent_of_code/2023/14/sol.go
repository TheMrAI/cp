package main

import (
	"bufio"
	"fmt"
	"os"
)

func RollRocksNorth(board [][]rune) [][]rune {
	nextAvailableSpace := make([]int, len(board[0]))
	for i := range board {
		for j := range board[i] {
			rollVertical(board, i, j, nextAvailableSpace, 1)
		}
	}

	return board
}

func RollRocksSouth(board [][]rune) [][]rune {
	nextAvailableSpace := make([]int, len(board[0]))
	for j := range nextAvailableSpace {
		nextAvailableSpace[j] = len(board) - 1
	}
	for i := len(board) - 1; i >= 0; i-- {
		for j := range board[i] {
			rollVertical(board, i, j, nextAvailableSpace, -1)
		}
	}

	return board
}

func rollVertical(board [][]rune, i int, j int, nextAvailableSpace []int, offset int) {
	if board[i][j] == '.' {
		return
	}
	if board[i][j] == '#' {
		nextAvailableSpace[j] = i + offset
		return
	}
	board[nextAvailableSpace[j]][j], board[i][j] = board[i][j], board[nextAvailableSpace[j]][j]
	nextAvailableSpace[j] += offset

}

func RollRocksWest(board [][]rune) [][]rune {
	nextAvailableSpace := make([]int, len(board))
	for j := range board[0] {
		for i := range board {
			rollHorizontal(board, i, j, nextAvailableSpace, 1)
		}
	}

	return board
}

func RollRocksEast(board [][]rune) [][]rune {
	nextAvailableSpace := make([]int, len(board))
	for i := range nextAvailableSpace {
		nextAvailableSpace[i] = len(board[0]) - 1
	}
	for j := len(board[0]) - 1; j >= 0; j-- {
		for i := range board {
			rollHorizontal(board, i, j, nextAvailableSpace, -1)
		}
	}

	return board
}

func rollHorizontal(board [][]rune, i int, j int, nextAvailableSpace []int, offset int) {
	if board[i][j] == '.' {
		return
	}
	if board[i][j] == '#' {
		nextAvailableSpace[i] = j + offset
		return
	}
	board[i][nextAvailableSpace[i]], board[i][j] = board[i][j], board[i][nextAvailableSpace[i]]
	nextAvailableSpace[i] += offset
}

func PartOne(board [][]rune) int {
	copied := make([][]rune, len(board))
	for i, line := range board {
		copied[i] = make([]rune, len(line))
		copy(copied[i], line)
	}
	return calculateNorthSideLoad(RollRocksNorth(copied))
}

func calculateNorthSideLoad(board [][]rune) int {
	totalLoad := 0
	for i := 1; i <= len(board); i++ {
		for _, val := range board[len(board)-i] {
			if val == 'O' {
				totalLoad += i
			}
		}
	}
	return totalLoad
}

// Well not production ready, safe struct.
// Don't matter here.
type CycleWatch struct {
	loads       []int
	detector    map[string]int
	cycleLength int
}

func (tracker *CycleWatch) CycleFound(load int, index int) bool {
	tracker.loads = append(tracker.loads, load)
	if len(tracker.loads) < 5 {
		return false
	}
	last := len(tracker.loads) - 1
	pattern := fmt.Sprintf("%v%v%v%v%v", tracker.loads[last], tracker.loads[last-1], tracker.loads[last-2], tracker.loads[last-3], tracker.loads[last-4])
	patternIndex, ok := tracker.detector[pattern]
	if ok {
		tracker.cycleLength = index - patternIndex
		return true
	}
	tracker.detector[pattern] = index
	return false
}

// No safety checks here, no logical protection either.
// This function has to be called after AddLoad returns true or it is
// undefined land.
func (tracker CycleWatch) ExtrapolateSteps(remainingSteps int) int {
	loadIndexFromCycleStart := remainingSteps % tracker.cycleLength
	// -1 because we have added the beginning of the detected cycle
	cycleStart := len(tracker.loads) - 1 - tracker.cycleLength
	fmt.Println(tracker.loads[cycleStart:])
	fmt.Println(loadIndexFromCycleStart)
	return tracker.loads[cycleStart+loadIndexFromCycleStart]
}

func PartTwo(board [][]rune, cycleCount int) int {
	copied := make([][]rune, len(board))
	for i, line := range board {
		copied[i] = make([]rune, len(line))
		copy(copied[i], line)
	}

	cycleWatch := CycleWatch{[]int{}, map[string]int{}, 0}
	for i := 0; i < cycleCount; i++ {
		RollRocksNorth(copied)
		RollRocksWest(copied)
		RollRocksSouth(copied)
		RollRocksEast(copied)

		// start checking for cycles after a thousand rounds
		if i > 1000 {
			load := calculateNorthSideLoad(copied)
			if cycleWatch.CycleFound(load, i) {
				return cycleWatch.ExtrapolateSteps(cycleCount - 1 - i)
			}
		}
	}
	return calculateNorthSideLoad(copied)
}

func ParseInput(lines []string) [][]rune {
	board := make([][]rune, len(lines))
	for i, line := range lines {
		board[i] = []rune(line)
	}
	return board
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	board := ParseInput(lines)

	fmt.Printf("Total load on north side support beans: %v\n", PartOne(board))
	fmt.Println("Part two")
	fmt.Printf("Total load on north side support beans after 1'000'000'000 cycles: %v\n", PartTwo(board, 1_000_000_000))
}

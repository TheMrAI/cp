package main

import (
	"bufio"
	"fmt"
	"os"
)

func ParseInput(lines []string) [][]byte {
	matrix := make([][]byte, len(lines))

	for i, line := range lines {
		matrix[i] = []byte(line)
	}

	return matrix
}

type Pair struct {
	I int
	J int
}

func PartOne(matrix [][]byte, stepCount int) int {
	reachablePlots := SimulateReachablePlots(matrix, stepCount)
	return reachablePlots[stepCount]
}

// This solution only works on the big test input, not on the example
// I couldn't figure it out myself.
// Reddit user charr3 https://www.reddit.com/user/charr3/ explains:
// "The main thing to notice for part 2 is that the grid is a square, and there are no obstacles in the same
// row/col of the starting point. Let f(n) be the number of spaces you can reach after n steps. Let X be the length of
// your input grid. f(n), f(n+X), f(n+2X), ...., is a quadratic, so you can find it by finding the first 3 values then
// use that to interpolate the final answer."
//
// It never occurred to me that the big input would have a different solution than one for the example.
// Even in the example I could see the diamond pattern, noted it in some commit, but never looked at the big input the
// same way!
func PartTwo(matrix [][]byte, stepCount int) int {
	fmt.Println(len(matrix))
	fmt.Println(len(matrix[0]))
	// triple the incoming grid in size
	tripledMatrix := ExpandMatrix(matrix, 9)
	// for _, l := range tripledMatrix {
	// 	fmt.Printf("%s\n", l)
	// }

	reachCounts := SimulateReachablePlots(tripledMatrix, 3*len(matrix))
	for i, val := range reachCounts {
		fmt.Printf("%3.v: %v\n", i, val)
	}

	return 0
}

func ExpandMatrix(matrix [][]byte, times int) [][]byte {
	expandedMatrix := make([][]byte, len(matrix)*times)
	sI := len(matrix) / 2
	sJ := len(matrix[0]) / 2
	// temporarily remove S for copying
	matrix[sI][sJ] = '.'
	for i := range expandedMatrix {
		expandedMatrix[i] = make([]byte, len(matrix[0])*times)
		for j := range expandedMatrix[i] {
			expandedMatrix[i][j] = matrix[i%len(matrix)][j%len(matrix[0])]
		}
	}
	matrix[sI][sJ] = 'S'
	expandedMatrix[len(expandedMatrix)/2][len(expandedMatrix[0])/2] = 'S'
	return expandedMatrix
}

func SimulateReachablePlots(matrix [][]byte, stepCount int) []int {
	reachableAfterSteps := []int{1}

	sCoord := Pair{0, 0}
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == 'S' {
				sCoord = Pair{i, j}
				break
			}
		}
	}

	offsets := [4]Pair{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	reachablePlots := map[Pair]struct{}{}
	reachablePlots[sCoord] = struct{}{}

	for i := 0; i < stepCount; i++ {
		nextReachablePlots := map[Pair]struct{}{}
		for coord := range reachablePlots {
			for _, offset := range offsets {
				nextCoord := Pair{coord.I + offset.I, coord.J + offset.J}
				if nextCoord.I < 0 || nextCoord.I >= len(matrix) || nextCoord.J < 0 || nextCoord.J >= len(matrix[0]) || matrix[nextCoord.I][nextCoord.J] == '#' {
					continue
				}
				nextReachablePlots[nextCoord] = struct{}{}
			}
		}
		// fmt.Printf("Step %v - covered %v , prev covered: %v delta: %v\n", i+1, len(nextReachablePlots), len(reachablePlots), len(nextReachablePlots)-len(reachablePlots))
		reachablePlots = nextReachablePlots
		reachableAfterSteps = append(reachableAfterSteps, len(reachablePlots))

		// fmt.Printf("Step %v - covered %v\n", i+1, len(reachablePlots))
		for l := range matrix {
			// fmt.Printf("%3.v : ", l)
			for m := range matrix[0] {
				_, ok := reachablePlots[Pair{l, m}]
				if ok {
					// fmt.Print("O")
					continue
				}
				// fmt.Printf("%c", matrix[l][m])
			}
			// fmt.Println()
		}
		// fmt.Println()
	}

	return reachableAfterSteps
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	matrix := ParseInput(lines)
	// for _, l := range matrix {
	// 	for _, c := range l {
	// 		fmt.Printf("%c ", c)
	// 	}
	// 	fmt.Println()
	// }

	fmt.Println("Part one")
	fmt.Printf("The number of reachable garden spots in 64 steps: %v\n", PartOne(matrix, 64))
	fmt.Println("Part two")
	fmt.Printf("The number of reachable garden spots in 26501365 steps: %v\n", PartTwo(matrix, 26501365))
}

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
		reachablePlots = nextReachablePlots

		fmt.Printf("Step %v - covered %v\n", i+1, len(reachablePlots))
		for l := range matrix {
			fmt.Printf("%3.v : ", l)
			for m := range matrix[0] {
				_, ok := reachablePlots[Pair{l, m}]
				if ok {
					fmt.Print("O ")
					continue
				}
				// fmt.Printf("%c ", matrix[l][m])
				fmt.Print("  ")
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return len(reachablePlots)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	matrix := ParseInput(lines)
	for _, l := range matrix {
		for _, c := range l {
			fmt.Printf("%c ", c)
		}
		fmt.Println()
	}

	fmt.Println("Part one")
	fmt.Printf("The number of reachable garden spots in 64 steps: %v\n", PartOne(matrix, 64))
}

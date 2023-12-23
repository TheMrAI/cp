package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

// not much really changes in this parsing we just make it explicit
// that we are working on bytes (also we may modify the matrix)
func ParseInput(lines []string) [][]byte {
	matrix := make([][]byte, len(lines))

	for i, line := range lines {
		matrix[i] = []byte(line)
	}

	return matrix
}

func PartOne(matrix [][]byte) int {
	// we could search for the start and end point, but I see no reason to do so now
	return LongestPathDFS(matrix, 0, 1, len(matrix)-1, len(matrix[0])-2)
}

func SlopeToOffset(slope byte) (int, int) {
	switch slope {
	case '^':
		return -1, 0
	case '>':
		return 0, 1
	case 'v':
		return 1, 0
	case '<':
		return 0, -1
	}
	log.Fatalf("Big nope")
	return 0, 0
}

func LongestPathDFS(matrix [][]byte, cI, cJ, tI, tJ int) int {
	if cI < 0 || cI >= len(matrix) || cJ < 0 || cJ >= len(matrix[0]) || matrix[cI][cJ] == '#' {
		return -1
	}
	if cI == tI && cJ == tJ {
		return 0
	}

	if matrix[cI][cJ] == '^' || matrix[cI][cJ] == '>' || matrix[cI][cJ] == 'v' || matrix[cI][cJ] == '<' {
		oI, oJ := SlopeToOffset(matrix[cI][cJ])
		length := LongestPathDFS(matrix, cI+oI, cJ+oJ, tI, tJ)
		if length == -1 {
			return -1
		}
		return length + 1
	}

	matrix[cI][cJ] = '#'
	offsetsI := []int{-1, 0, 1, 0}
	offsetsJ := []int{0, 1, 0, -1}
	pathLengths := []int{-1, -1, -1, -1}
	for i := 0; i < 4; i++ {
		nI := cI + offsetsI[i]
		nJ := cJ + offsetsJ[i]
		pathLengths[i] = LongestPathDFS(matrix, nI, nJ, tI, tJ)
	}
	matrix[cI][cJ] = '.'

	longest := math.MinInt
	for _, length := range pathLengths {
		if longest < length {
			longest = length
		}
	}

	if longest == -1 {
		return -1
	}
	return longest + 1
}

func PartTwo(matrix [][]byte) int {
	// clear all slopes
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == '^' || matrix[i][j] == '>' || matrix[i][j] == 'v' || matrix[i][j] == '<' {
				matrix[i][j] = '.'
			}
		}
	}
	return LongestPathDFS(matrix, 0, 1, len(matrix)-1, len(matrix[0])-2)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)
	for i := range data {
		for j := range data[i] {
			fmt.Printf("%c ", data[i][j])
		}
		fmt.Println()
	}

	fmt.Println("Part one")
	fmt.Printf("The length of the longest possible walk with slopes: %v\n", PartOne(data))
	fmt.Println("Part two")
	fmt.Printf("The length of the longest possible walk without slopes: %v\n", PartTwo(data))
	// Total runtime:
	// real    2m31.578s
	// user    2m31.628s
	// sys     0m0.144s
	//
	// Meh, not the best, but good enough as a start.
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func ParseInput(lines []string) [][]int {
	matrix := make([][]int, len(lines))

	numberRegexp := regexp.MustCompile("\\d")

	for i, line := range lines {
		numbers := numberRegexp.FindAllString(line, -1)
		if len(line) != len(numbers) {
			log.Fatalf("Number of parsed numbers doesn't equal the number of characters in: '%v'\n", line)
		}
		matrix[i] = make([]int, len(numbers))
		for j := range numbers {
			val, err := strconv.Atoi(numbers[j])
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'\n", numbers[j])
			}
			matrix[i][j] = val
		}
	}

	return matrix
}

func PartOne(matrix [][]int) int {
	sum := 0
	for _, v := range FindLowPoints(matrix) {
		sum += matrix[v.I][v.J] + 1

	}

	return sum
}

func PartTwo(matrix [][]int) int {
	visited := make([]bool, len(matrix)*len(matrix[0]))

	basinSizes := []int{}
	for _, v := range FindLowPoints(matrix) {
		basinSizes = append(basinSizes, countBasin(matrix, v, visited))
	}

	sort.Slice(basinSizes, func(i, j int) bool {
		return basinSizes[i] > basinSizes[j]
	})
	return basinSizes[0] * basinSizes[1] * basinSizes[2]
}

// We do not even consider the slope of the basin at all, because
// the problem clearly states that only 9 are no part of a basin,
// while every other tile must be.
func countBasin(matrix [][]int, start Index, visited []bool) int {
	basinSize := 0

	di := []int{0, 0, 1, -1}
	dj := []int{1, -1, 0, 0}

	queue := []Index{start}

	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]

		if visited[top.I*len(matrix[0])+top.J] || matrix[top.I][top.J] == 9 {
			continue
		}
		basinSize += 1
		visited[top.I*len(matrix[0])+top.J] = true

		for offsetIndex := range di {
			nI := top.I + di[offsetIndex]
			nJ := top.J + dj[offsetIndex]
			if nI < 0 || nI >= len(matrix) || nJ < 0 || nJ >= len(matrix[top.I]) {
				continue
			}
			queue = append(queue, Index{nI, nJ})
		}
	}

	return basinSize
}

type Index struct {
	I int
	J int
}

func FindLowPoints(matrix [][]int) []Index {
	lowPoints := []Index{}

	di := []int{0, 0, 1, -1}
	dj := []int{1, -1, 0, 0}

	for i := range matrix {
		for j := range matrix[i] {
			lowPoint := true
			for offsetIndex := range di {
				nI := i + di[offsetIndex]
				nJ := j + dj[offsetIndex]
				if nI < 0 || nI >= len(matrix) || nJ < 0 || nJ >= len(matrix[i]) {
					continue
				}
				if matrix[nI][nJ] <= matrix[i][j] {
					lowPoint = false
					break
				}
			}
			if lowPoint {
				lowPoints = append(lowPoints, Index{i, j})
			}
		}
	}

	return lowPoints
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	matrix := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("The sum of the low point values is: %v\n", PartOne(matrix))
	fmt.Println("Part two")
	fmt.Printf("The multiple of the sizes of the three largest basins: %v\n", PartTwo(matrix))
}

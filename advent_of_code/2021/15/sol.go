package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

func ParseInput(lines []string) [][]int {
	matrix := make([][]int, len(lines))

	numberRegexp := regexp.MustCompile("\\d")
	for i, line := range lines {
		numberText := numberRegexp.FindAllString(line, -1)
		matrix[i] = make([]int, len(numberText))
		for j, text := range numberText {
			number, err := strconv.Atoi(text)
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'\n", text)
			}
			matrix[i][j] = number
		}
	}

	return matrix
}

type Pair struct {
	I int
	J int
}

// real    0m1.689s
// user    0m1.740s
// sys     0m0.060s
func PartOne(matrix [][]int) int {
	return LowestRiskPath(matrix, Pair{0, 0}, Pair{len(matrix) - 1, len(matrix[0]) - 1})
}

// Part one and two together
// real    4m21.020s
// user    4m25.607s
// sys     0m3.812s
func PartTwo(matrix [][]int) int {
	expanded := ExpandMatrix(matrix, 5)
	// fmt.Println(expanded)
	return LowestRiskPath(expanded, Pair{0, 0}, Pair{len(expanded) - 1, len(expanded[0]) - 1})
}

func ExpandMatrix(matrix [][]int, times uint) [][]int {
	expandedMatrix := make([][]int, len(matrix)*int(times))
	for i := range expandedMatrix {
		expandedMatrix[i] = make([]int, len(matrix[0])*int(times))
	}

	for i := range matrix {
		for j := range matrix[0] {
			for repeatI := 0; repeatI < int(times); repeatI++ {
				nI := i + (repeatI * len(matrix))
				for repeatJ := 0; repeatJ < int(times); repeatJ++ {
					nJ := j + (repeatJ * len(matrix[0]))
					value := matrix[i][j] + ManhattanDistance(Pair{i, j}, Pair{nI, nJ})
					if value > 9 {
						value = ((value - 1) % 9) + 1
					}
					expandedMatrix[nI][nJ] = value
				}
			}
		}
	}
	return expandedMatrix
}

func LowestRiskPath(matrix [][]int, from Pair, to Pair) int {
	lowestRiskMatrix := make([][]int, len(matrix))
	for i := range lowestRiskMatrix {
		lowestRiskMatrix[i] = make([]int, len(matrix[i]))
		for j := range lowestRiskMatrix[i] {
			lowestRiskMatrix[i][j] = math.MaxInt
		}
	}

	toCheck := PriorityQueue{}
	heap.Init(&toCheck)
	insertOptions(matrix, &toCheck, to, Block{from, 0, ManhattanDistance(from, to), 0})

	for len(toCheck) != 0 {
		checking := heap.Pop(&toCheck).(*Block)

		// we have already found a lower risk path to this element
		if checking.TotalRisk >= lowestRiskMatrix[checking.Pos.I][checking.Pos.J] {
			continue
		}

		lowestRiskMatrix[checking.Pos.I][checking.Pos.J] = checking.TotalRisk
		if checking.Pos == to {
			continue
		}

		insertOptions(matrix, &toCheck, to, *checking)
	}
	return lowestRiskMatrix[to.I][to.J]
}

func insertOptions(matrix [][]int, toBeChecked *PriorityQueue, target Pair, current Block) {
	deltaI := []int{-1, 0, 1, 0}
	deltaJ := []int{0, 1, 0, -1}

	for i := 0; i < len(deltaI); i++ {
		dI := deltaI[i]
		dJ := deltaJ[i]
		nextLocation := Pair{current.Pos.I + dI, current.Pos.J + dJ}
		if nextLocation.I < 0 || nextLocation.I >= len(matrix) || nextLocation.J < 0 || nextLocation.J >= len(matrix[0]) {
			continue
		}
		risk := current.TotalRisk + matrix[nextLocation.I][nextLocation.J]
		heap.Push(toBeChecked, &Block{nextLocation, risk, ManhattanDistance(nextLocation, target), 0})
	}
}

type Block struct {
	Pos       Pair
	TotalRisk int
	priority  int // manhattan distance
	index     int
}

func ManhattanDistance(from, to Pair) int {
	deltaI := to.I - from.I
	if deltaI < 0 {
		deltaI = 0 - deltaI
	}
	deltaJ := to.J - from.J
	if deltaJ < 0 {
		deltaJ = 0 - deltaJ
	}
	return deltaI + deltaJ
}

type PriorityQueue []*Block

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want the lowest priorities first, because we
	// use the manhattan distance as priority
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Block)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
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
	fmt.Printf("The lowest total risk of any path from top left to bottom right corner is: %v\n", PartOne(matrix))
	fmt.Println("Part two")
	fmt.Printf("On the 5 times bigger map the lowest total risk of any path from top left to bottom right corner is: %v\n", PartTwo(matrix))
}

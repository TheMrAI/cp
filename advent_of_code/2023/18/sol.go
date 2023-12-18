package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func stringToDirection(c string) Direction {
	switch c {
	case "U":
		return Up
	case "R":
		return Right
	case "D":
		return Down
	case "L":
		return Left
	default:
		log.Fatalf("Unsupported direction received: '%s'\n", c)
	}
	return Left
}

type DigStep struct {
	Dir   Direction
	Count int
	Color string
}

func ParseInput(lines []string) []DigStep {
	data := []DigStep{}

	directionsRegexp := regexp.MustCompile("(\\w{1}) (\\d+) \\((#\\w+)\\)")
	for _, line := range lines {
		matches := directionsRegexp.FindAllStringSubmatch(line, -1)
		if len(matches) != 1 {
			log.Fatalf("Could not parse input: '%v'\n", line)
		}
		stepCount, err := strconv.Atoi(matches[0][2])
		if err != nil {
			log.Fatalf("Failed int conversion on input: '%v'\n", matches[0][2])
		}
		data = append(data, DigStep{stringToDirection(matches[0][1]), stepCount, matches[0][3]})
	}

	return data
}

type Pair struct {
	I int
	J int
}

type LineSegment struct {
	StartPoint Pair
	LastPoint  Pair
	Color      string
}

func DigDirectionsToLineSegments(startPoint Pair, digSteps []DigStep) []LineSegment {
	lineSegments := []LineSegment{}

	currentPoint := startPoint
	for _, digStep := range digSteps {
		lastPoint := currentPoint
		switch digStep.Dir {
		case Up:
			lastPoint.I -= digStep.Count
		case Right:
			lastPoint.J += digStep.Count
		case Down:
			lastPoint.I += digStep.Count
		case Left:
			lastPoint.J -= digStep.Count
		}
		lineSegments = append(lineSegments, LineSegment{currentPoint, lastPoint, digStep.Color})
		currentPoint = lastPoint
	}

	return lineSegments
}

func Min(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

func Max(lhs, rhs int) int {
	if lhs < rhs {
		return rhs
	}
	return lhs
}

func PartOne(lineSegments []LineSegment) int {
	matrix := LineSegmentsToMatrix(lineSegments)
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%c", matrix[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
	excavated := Excavate(matrix)
	for i := range excavated {
		for j := range excavated[i] {
			fmt.Printf("%c", excavated[i][j])
		}
		fmt.Println()
	}
	count := 0
	for _, row := range excavated {
		for _, val := range row {
			if val == '#' {
				count++
			}
		}
	}
	return count
}

func LineSegmentsToMatrix(lineSegments []LineSegment) [][]byte {
	minI := math.MaxInt
	maxI := math.MinInt
	minJ := math.MaxInt
	maxJ := math.MinInt
	for _, segment := range lineSegments {
		minI = Min(minI, segment.StartPoint.I)
		minI = Min(minI, segment.LastPoint.I)
		minJ = Min(minJ, segment.StartPoint.J)
		minJ = Min(minJ, segment.LastPoint.J)

		maxI = Max(maxI, segment.StartPoint.I)
		maxI = Max(maxI, segment.LastPoint.I)
		maxJ = Max(maxJ, segment.StartPoint.J)
		maxJ = Max(maxJ, segment.LastPoint.J)
	}
	// normalize points to start from 0,0
	deltaI := 0 - minI
	deltaJ := 0 - minJ
	maxI = maxI + deltaI
	maxJ = maxJ + deltaJ

	matrix := make([][]byte, maxI+1)
	for i := range matrix {
		matrix[i] = make([]byte, maxJ+1)
		for j := range matrix[i] {
			matrix[i][j] = '.'
		}
	}

	for _, lineSegment := range lineSegments {
		currentI := lineSegment.StartPoint.I + deltaI
		currentJ := lineSegment.StartPoint.J + deltaJ
		lastI := lineSegment.LastPoint.I + deltaI
		lastJ := lineSegment.LastPoint.J + deltaJ

		offsetI := 0
		if lastI-currentI != 0 {
			offsetI = 1
			if lastI-currentI < 0 {
				offsetI = -1
			}
		}

		offsetJ := 0
		if lastJ-currentJ != 0 {
			offsetJ = 1
			if lastJ-currentJ < 0 {
				offsetJ = -1
			}
		}

		for currentI != lastI || currentJ != lastJ {
			matrix[currentI][currentJ] = '#'
			currentI += offsetI
			currentJ += offsetJ
		}
	}

	return matrix
}

type DSU struct {
	parent []int
	rank   []int
}

func (dsu *DSU) MakeSet(data []int) {
	dsu.parent = make([]int, len(data))
	dsu.rank = make([]int, len(data))
	for _, val := range data {
		dsu.parent[val] = val
	}
}

func (dsu DSU) FindSet(value int) int {
	if value == dsu.parent[value] {
		return value
	}
	dsu.parent[value] = dsu.FindSet(dsu.parent[value])
	return dsu.parent[value]
}

func (dsu *DSU) UnionSets(lhs, rhs int) {
	lhs = dsu.FindSet(lhs)
	rhs = dsu.FindSet(rhs)
	if lhs != rhs {
		if dsu.rank[lhs] < dsu.rank[rhs] {
			lhs, rhs = rhs, lhs
		}
		dsu.parent[rhs] = lhs
		if dsu.rank[lhs] == dsu.rank[rhs] {
			dsu.rank[lhs]++
		}
	}
}

func Excavate(matrix [][]byte) [][]byte {
	excavated := make([][]byte, len(matrix))
	for i := range matrix {
		excavated[i] = make([]byte, len(matrix[0]))
		copy(excavated[i], matrix[i])
	}
	cellIds := make([]int, len(matrix)*len(matrix[0]))
	for i := range cellIds {
		cellIds[i] = i
	}
	dsu := DSU{}
	dsu.MakeSet(cellIds)

	initialTrenchSet := -1
	reachingEdgeSets := map[int]struct{}{}

	offsets := [4]Pair{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	for i := range matrix {
		for j := range matrix[0] {
			cellId := i*len(matrix[0]) + j
			for _, offset := range offsets {
				nI := i + offset.I
				nJ := j + offset.J
				if nI < 0 || nI >= len(matrix) || nJ < 0 || nJ >= len(matrix[0]) {
					reachingEdgeSets[dsu.FindSet(cellId)] = struct{}{}
					continue
				}
				if matrix[i][j] == matrix[nI][nJ] {
					nextCellId := nI*len(matrix[0]) + nJ
					dsu.UnionSets(cellId, nextCellId)
				}
				if matrix[i][j] == '#' {
					initialTrenchSet = dsu.FindSet(cellId)
				}
			}
		}
	}

	for i := range excavated {
		for j := range excavated[0] {
			cellId := i*len(excavated[0]) + j
			cellSet := dsu.FindSet(cellId)
			_, inOuterSets := reachingEdgeSets[cellSet]
			if cellSet == initialTrenchSet || inOuterSets {
				continue
			}
			excavated[i][j] = '#'
		}
	}

	return excavated
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)
	lineSegments := DigDirectionsToLineSegments(Pair{0, 0}, data)

	fmt.Println("Part one")
	fmt.Printf("Total cubic meters excavated: %v\n", PartOne(lineSegments))
}

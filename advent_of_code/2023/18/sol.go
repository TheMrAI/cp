package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
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
}

func (segment LineSegment) Length() int {
	// this only works because we only have horizontal or vertical lines
	// plus one adds back the end point
	return 1 + Abs(segment.StartPoint.I-segment.LastPoint.I) + Abs(segment.StartPoint.J-segment.LastPoint.J)
}

func DigDirectionsToLineSegments(startPoint Pair, digSteps []DigStep, parseFromColor bool) []LineSegment {
	lineSegments := []LineSegment{}

	currentPoint := startPoint
	for _, digStep := range digSteps {
		lastPoint := currentPoint
		digDir := digStep.Dir
		digCount := digStep.Count
		lastCharIndex := len(digStep.Color) - 1
		if parseFromColor {
			dirValue, err := strconv.Atoi(digStep.Color[lastCharIndex:])
			if err != nil {
				log.Fatalf("Failed int conversion on input: '%v'\n", digStep.Color[lastCharIndex:])
			}
			// we shift the received value forward by one so it aligns with our ordering
			digDir = Direction((dirValue + 1) % 4)
			countValue, err := strconv.ParseInt(digStep.Color[1:lastCharIndex], 16, 64)
			if err != nil {
				log.Fatalf("Failed hexadecimal int conversion on input: '%v'\n", digStep.Color[1:lastCharIndex])
			}
			digCount = int(countValue)
		}

		switch digDir {
		case Up:
			lastPoint.I -= digCount
		case Right:
			lastPoint.J += digCount
		case Down:
			lastPoint.I += digCount
		case Left:
			lastPoint.J -= digCount
		}

		lineSegments = append(lineSegments, LineSegment{currentPoint, lastPoint})
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

func PartOne(data []DigStep) int {
	lineSegments := DigDirectionsToLineSegments(Pair{0, 0}, data, false)
	matrix := LineSegmentsToMatrix(lineSegments)
	excavated := Excavate(matrix)
	// for i := range matrix {
	// 	fmt.Printf("%3.v ", i)
	// 	for j := range matrix[i] {
	// 		fmt.Printf("%c", matrix[i][j])
	// 	}
	// 	fmt.Println()
	// }
	// fmt.Println()
	// for i := range excavated {
	// 	fmt.Printf("%3.v ", i)
	// 	for j := range excavated[i] {
	// 		fmt.Printf("%c", excavated[i][j])
	// 	}
	// 	fmt.Println()
	// }
	count := 0
	for i, row := range excavated {
		rowCount := 0
		for _, val := range row {
			if val == '#' {
				count++
				rowCount++
			}
		}
		fmt.Printf("i: %v - %v\n", i, rowCount)
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

// Alright all right, all right. Grug dev has to do complexity demon because there
// isn't enough memory in the universe to load this abomination in.
// What we do is a scanline solution, from left to right, bottom to top, counting
// the internal sections of the shape without the edges.
// The algorithm may not be as efficient as it could be perhaps, but it is already rather
// difficult to get it right and runs in less than a second for part one and two together with
// compilation. So it isn't necessary to complicate any more.
func PartTwo(data []DigStep) int {
	lineSegments := DigDirectionsToLineSegments(Pair{0, 0}, data, true)
	verticalLines := LineSegmentsToVerticalLines(lineSegments)

	sort.Slice(verticalLines, func(i, j int) bool {
		return verticalLines[i].Axis < verticalLines[j].Axis
	})

	totalExcavated := 0
	for vI := range verticalLines {
		vertical := verticalLines[vI]
		if vertical.Terminating {
			continue
		}

		for scan := vertical.A; scan <= vertical.B; scan++ {
			// a horizontal starts from this point and goes to the left => not in internal section of the shape
			if (scan == vertical.A && vertical.AHDelta > 0) || (scan == vertical.B && vertical.BHDelta > 0) {
				continue
			}
			// find terminating vertical line
			tV := vI + 1
			for ; tV < len(verticalLines); tV++ {
				if verticalLines[tV].A <= scan && scan <= verticalLines[tV].B {
					break
				}
			}
			verticalLines[tV].Terminating = true
			totalExcavated += verticalLines[tV].Axis - vertical.Axis - 1
		}
	}

	// add the original trench and remove the doubly counted end points once
	for _, segment := range lineSegments {
		totalExcavated += segment.Length()
	}
	totalExcavated -= len(lineSegments)
	return totalExcavated
}

type VerticalLine struct {
	A           int // the smaller I coordinate of the vertical line
	AHDelta     int // the the offset the horizontal line attached point A would offset the Axis of the vertical line if we followed it (positive means horizontal goes right)
	B           int
	BHDelta     int
	Axis        int  // the J coordinate of the vertical line
	Terminating bool // whether or not the vertical line terminates the shapes enclosed section
}

func LineSegmentsToVerticalLines(lineSegments []LineSegment) []VerticalLine {
	verticalLines := []VerticalLine{}

	for i := range lineSegments {
		// skip horizontals
		if lineSegments[i].StartPoint.J != lineSegments[i].LastPoint.J {
			continue
		}
		axis := lineSegments[i].StartPoint.J
		a := lineSegments[i].StartPoint.I
		b := lineSegments[i].LastPoint.I

		hIndexOne := i + 1
		if hIndexOne == len(lineSegments) {
			hIndexOne = 0
		}
		hDeltaOne := lineSegments[hIndexOne].Length() - 1
		if lineSegments[hIndexOne].StartPoint.J < axis || lineSegments[hIndexOne].LastPoint.J < axis {
			hDeltaOne *= -1
		}

		hIndexTwo := i - 1
		if hIndexTwo < 0 {
			hIndexTwo = len(lineSegments) - 1
		}
		hDeltaTwo := lineSegments[hIndexTwo].Length() - 1
		if lineSegments[hIndexTwo].StartPoint.J < axis || lineSegments[hIndexTwo].LastPoint.J < axis {
			hDeltaTwo *= -1
		}

		aHDelta := hDeltaOne
		bHDelta := hDeltaTwo
		// if horizontal doesn't match the start point of our vertical then we swap the hDeltas
		if lineSegments[i].StartPoint != lineSegments[hIndexOne].StartPoint && lineSegments[i].LastPoint != lineSegments[hIndexOne].LastPoint {
			aHDelta, bHDelta = bHDelta, aHDelta
		}

		if a > b {
			// reorder data so a -> b is always increasing
			a, aHDelta, b, bHDelta = b, bHDelta, a, aHDelta
		}
		verticalLines = append(verticalLines, VerticalLine{a, aHDelta, b, bHDelta, axis, false})
	}

	return verticalLines
}

func Abs(val int) int {
	if val < 0 {
		return 0 - val
	}
	return val
}

// Total runtime with compilation and running part one and two together
// real    0m0.892s
// user    0m0.891s
// sys     0m0.066s
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Total cubic meters excavated: %v\n", PartOne(data))
	// Solution: 52240187443190
	fmt.Println("Part two")
	fmt.Printf("Total cubic meters excavated using colors for instructions: %v\n", PartTwo(data))
}

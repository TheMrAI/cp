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

// Alright all right, all right. Grug dev has to do complexity demon because there
// isn't enough memory in the universe to load this abomination in.
// What we do is a scanline solution.
// Order the segments by coordinates, then scan from top to bottom.
// On each horizontal section we decide what it does:
// - if it is disjoint or extends previous ones, it builds
// - if it is in the middle of a section it cuts
func PartTwo(data []DigStep) int {
	lineSegments := DigDirectionsToLineSegments(Pair{0, 0}, data, false)

	horizontals := []LineSegment{}
	for i := range lineSegments {
		if lineSegments[i].StartPoint.I != lineSegments[i].LastPoint.I {
			continue
		}
		// based on how we built the line segments we know that next to a horizontal section there can only be
		// two vertical ones
		horizontal := LineSegment{lineSegments[i].StartPoint, lineSegments[i].LastPoint}
		// swap the points so that J always increases from StartPoint -> LastPoint
		if horizontal.StartPoint.J > horizontal.LastPoint.J {
			horizontal.StartPoint, horizontal.LastPoint = horizontal.LastPoint, horizontal.StartPoint
		}
		horizontals = append(horizontals, horizontal)
	}

	// prepping the data for the sweep-line
	sort.Slice(horizontals, func(i, j int) bool {
		iHeight := horizontals[i].StartPoint.I
		jHeight := horizontals[j].StartPoint.I
		if iHeight == jHeight {
			return horizontals[i].StartPoint.J < horizontals[j].StartPoint.J
		}
		return iHeight < jHeight
	})

	for i := range horizontals {
		fmt.Printf("i: %v, %v\n", i, horizontals[i])
	}

	activeHorizontals, i := horizontalWithSameHeightFrom(horizontals, 0)
	totalExcavated := 0
	for _, active := range activeHorizontals {
		totalExcavated += active.Length()
	}
	for i < len(horizontals) {
		currentHeigh := horizontals[i].StartPoint.I
		for j, active := range activeHorizontals {
			width := active.Length()
			// fmt.Printf("Width: %v\n", width)
			height := Abs(active.StartPoint.I-currentHeigh) - 1
			// fmt.Printf("Active: %v, currentHeight: %v\n", active, currentHeigh)
			fmt.Printf("Excavated this round: %v (total: %v)\n", width*height, totalExcavated)
			totalExcavated += width * height
			// adjust height of active element after being accounted for
			activeHorizontals[j].StartPoint = Pair{currentHeigh, active.StartPoint.J}
			activeHorizontals[j].LastPoint = Pair{currentHeigh, active.LastPoint.J}
		}

		sameHeightSegments, sameHeightUntil := horizontalWithSameHeightFrom(horizontals, i)
		zippedHorizontals := zipLineSegmentSlices(activeHorizontals, sameHeightSegments)
		totalExcavated += countCoveredDistance(zippedHorizontals)
		fmt.Printf("Transition ex: %v\n", countCoveredDistance(zippedHorizontals))

		fmt.Printf("Zipped: %v\n", zippedHorizontals)
		// generate new activeHorizontals by merging or dividing existing ones
		// this is made easier by ordering the horizontals on the j axis as well
		nextActiveHorizontals := []LineSegment{}
		for zIndex := 0; zIndex < len(zippedHorizontals); zIndex++ {
			if len(nextActiveHorizontals) == 0 {
				nextActiveHorizontals = append(nextActiveHorizontals, zippedHorizontals[zIndex])
				continue
			}
			lastIndex := len(nextActiveHorizontals) - 1
			if nextActiveHorizontals[lastIndex].LastPoint.J < zippedHorizontals[zIndex].StartPoint.J {
				// just add
				fmt.Println("Add")
				nextActiveHorizontals = append(nextActiveHorizontals, zippedHorizontals[zIndex])
			} else if nextActiveHorizontals[lastIndex].LastPoint.J == zippedHorizontals[zIndex].StartPoint.J {
				// concatenate horizontal segment
				fmt.Println("Concat")
				nextActiveHorizontals[lastIndex].LastPoint.J = zippedHorizontals[zIndex].LastPoint.J
			} else if nextActiveHorizontals[lastIndex].StartPoint.J == zippedHorizontals[zIndex].StartPoint.J {
				// shorten segment
				fmt.Println("Short")
				nextActiveHorizontals[lastIndex].StartPoint.J = zippedHorizontals[zIndex].LastPoint.J
			} else if zippedHorizontals[zIndex].LastPoint.J == nextActiveHorizontals[lastIndex].LastPoint.J {
				// shorten segment
				fmt.Println("Short")
				nextActiveHorizontals[lastIndex].LastPoint.J = zippedHorizontals[zIndex].StartPoint.J
			} else if nextActiveHorizontals[lastIndex].StartPoint.J < zippedHorizontals[zIndex].StartPoint.J && zippedHorizontals[zIndex].LastPoint.J < nextActiveHorizontals[lastIndex].LastPoint.J {
				// chop
				fmt.Println("Chop")
				firstPart := LineSegment{nextActiveHorizontals[lastIndex].StartPoint, zippedHorizontals[zIndex].StartPoint}
				secondPart := LineSegment{zippedHorizontals[zIndex].LastPoint, nextActiveHorizontals[lastIndex].LastPoint}
				nextActiveHorizontals = nextActiveHorizontals[:len(nextActiveHorizontals)-1]
				nextActiveHorizontals = append(nextActiveHorizontals, firstPart, secondPart)
			} else {
				log.Fatalf("Unexpected state achieved. Check logic.")
			}
		}
		activeHorizontals = nextActiveHorizontals
		fmt.Printf("Next actives: %v\n", activeHorizontals)
		i = sameHeightUntil
	}

	return totalExcavated
}

func horizontalWithSameHeightFrom(segments []LineSegment, fromIndex int) ([]LineSegment, int) {
	sameHeightIndex := fromIndex
	sameHeightUntil := fromIndex + 1
	for ; sameHeightUntil < len(segments) && segments[sameHeightUntil].StartPoint.I == segments[sameHeightIndex].StartPoint.I; sameHeightUntil++ {
	}
	sameHeightSegments := segments[sameHeightIndex:sameHeightUntil]
	return sameHeightSegments, sameHeightUntil
}

func zipLineSegmentSlices(lhs []LineSegment, rhs []LineSegment) []LineSegment {
	zippedHorizontals := []LineSegment{}
	lhsIndex := 0
	rhsIndex := 0
	for lhsIndex < len(lhs) && rhsIndex < len(rhs) {
		if lhs[lhsIndex].StartPoint.J <= rhs[rhsIndex].StartPoint.J {
			zippedHorizontals = append(zippedHorizontals, lhs[lhsIndex])
			lhsIndex++
			continue
		}
		zippedHorizontals = append(zippedHorizontals, rhs[rhsIndex])
		rhsIndex++
	}
	if lhsIndex < len(lhs) {
		zippedHorizontals = append(zippedHorizontals, lhs[lhsIndex:]...)
	}
	if rhsIndex < len(rhs) {
		zippedHorizontals = append(zippedHorizontals, rhs[rhsIndex:]...)
	}
	return zippedHorizontals
}

func countCoveredDistance(segments []LineSegment) int {
	totalCovered := 0
	coveredUntil := 0
	for i := range segments {
		if i == 0 ||coveredUntil <= segments[i].StartPoint.J {
			totalCovered = segments[i].Length()
			coveredUntil = segments[i].LastPoint.J + 1
			continue
		}
		if coveredUntil >= segments[i].LastPoint.J {
			continue
		}
		totalCovered += 1 + segments[i].LastPoint.J - coveredUntil
		coveredUntil = segments[i].LastPoint.J + 1
	}
	return totalCovered
}

func Abs(val int) int {
	if val < 0 {
		return 0 - val
	}
	return val
}

func getJOnHeight(segment LineSegment, height int) int {
	if segment.StartPoint.I == height {
		return segment.StartPoint.J
	}
	return segment.LastPoint.J
}

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
	fmt.Println("Part two")
	fmt.Printf("Total cubic meters excavated using colors for instructions: %v\n", PartTwo(data))
}

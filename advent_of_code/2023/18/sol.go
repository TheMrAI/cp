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

type HorizontalSegment struct {
	Line       LineSegment // just the height of the segment, can't be index, we may merge them
	LeftIndex  int         // vertical segments can use the index
	RightIndex int
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

	horizontals := []HorizontalSegment{}
	for i := range lineSegments {
		fmt.Printf("i: %v, %v\n", i, lineSegments[i])
		if lineSegments[i].StartPoint.I != lineSegments[i].LastPoint.I {
			continue
		}
		prevIndex := len(lineSegments) - 1
		if i > 0 {
			prevIndex = i - 1
		}
		nextIndex := 0
		if i < len(lineSegments)-1 {
			nextIndex = i + 1
		}
		// based on how we built the line segments we know that next to a horizontal section there can only be
		// two vertical ones
		horizontal := LineSegment{lineSegments[i].StartPoint, lineSegments[i].LastPoint}
		if horizontal.StartPoint.J > horizontal.LastPoint.J {
			horizontal.StartPoint, horizontal.LastPoint = horizontal.LastPoint, horizontal.StartPoint
		}
		horizontals = append(horizontals, HorizontalSegment{horizontal, prevIndex, nextIndex})
	}

	sort.Slice(horizontals, func(i, j int) bool {
		iHeight := horizontals[i].Line.StartPoint.I
		jHeight := horizontals[j].Line.StartPoint.I
		if iHeight == jHeight {
			return horizontals[i].Line.StartPoint.J < horizontals[j].Line.StartPoint.J
		}
		return iHeight < jHeight
	})

	activeHorizontals, i := horizontalWithSameHeightFrom(horizontals, 0)
	totalExcavated := 0
	for i < len(horizontals) {
		currentHeigh := horizontals[i].Line.StartPoint.I
		for _, active := range activeHorizontals {
			width := active.Line.Length() - 2
			fmt.Printf("Width: %v\n", width)
			height := Abs(active.Line.StartPoint.I-currentHeigh) - 1
			fmt.Printf("Active: %v, currentHeight: %v\n", active, currentHeigh)
			fmt.Printf("Excavated: %v\n", width*height)
			totalExcavated += width * height
		}

		sameHeightSegments, sameHeightUntil := horizontalWithSameHeightFrom(horizontals, i)

		zippedHorizontals := zipHorizontalSegmentSlices(activeHorizontals, sameHeightSegments)
		// we must adjust all horizontal segments as if they were on the same height
		for l := range zippedHorizontals {
			zippedHorizontals[l].Line.StartPoint = Pair{currentHeigh, zippedHorizontals[l].Line.StartPoint.J}
			zippedHorizontals[l].Line.LastPoint = Pair{currentHeigh, zippedHorizontals[l].Line.LastPoint.J}
		}
		fmt.Printf("Zipped: %v\n", zippedHorizontals)
		// generate new activeHorizontals by merging or dividing existing ones
		// this is made easier by ordering the horizontals on the j axis as well
		nextActiveHorizontals := []HorizontalSegment{}
		nextActiveHorizontals = append(nextActiveHorizontals, zippedHorizontals[0])
		for zIndex := 1; zIndex < len(zippedHorizontals); zIndex++ {
			lastIndex := len(nextActiveHorizontals) - 1
			if nextActiveHorizontals[lastIndex].Line.LastPoint.J < zippedHorizontals[zIndex].Line.StartPoint.J {
				// just add
				nextActiveHorizontals = append(nextActiveHorizontals, zippedHorizontals[zIndex])
			} else if nextActiveHorizontals[lastIndex].Line.LastPoint.J == zippedHorizontals[zIndex].Line.StartPoint.J {
				// concatenate horizontal segment
				nextActiveHorizontals[lastIndex].Line.LastPoint.J = zippedHorizontals[zIndex].Line.LastPoint.J
				nextActiveHorizontals[lastIndex].RightIndex = zippedHorizontals[zIndex].RightIndex
			} else if nextActiveHorizontals[lastIndex].Line.StartPoint.J <= zippedHorizontals[zIndex].Line.StartPoint.J && zippedHorizontals[zIndex].Line.LastPoint.J <= nextActiveHorizontals[lastIndex].Line.LastPoint.J {
				// chop
				firstPart := HorizontalSegment{LineSegment{nextActiveHorizontals[lastIndex].Line.StartPoint, zippedHorizontals[zIndex].Line.StartPoint}, nextActiveHorizontals[lastIndex].LeftIndex, zippedHorizontals[zIndex].LeftIndex}
				secondPart := HorizontalSegment{LineSegment{zippedHorizontals[zIndex].Line.LastPoint, nextActiveHorizontals[lastIndex].Line.LastPoint}, zippedHorizontals[zIndex].RightIndex, nextActiveHorizontals[lastIndex].RightIndex}
				nextActiveHorizontals = nextActiveHorizontals[:len(nextActiveHorizontals)-1]
				if firstPart.Line.StartPoint.J != firstPart.Line.LastPoint.J {
					nextActiveHorizontals = append(nextActiveHorizontals, firstPart)
				}
				if secondPart.Line.StartPoint.J != secondPart.Line.LastPoint.J {
					nextActiveHorizontals = append(nextActiveHorizontals, secondPart)
				}
			} else {
				log.Fatalf("Unexpected state achieved. Check logic.")
			}
		}
		activeHorizontals = nextActiveHorizontals
		i = sameHeightUntil
	}

	// add original trench to excavation total
	for _, segment := range lineSegments {
		totalExcavated += segment.Length()
	}
	// counted all endpoints twice, removing them once
	totalExcavated -= len(lineSegments)
	return totalExcavated
}

func horizontalWithSameHeightFrom(segments []HorizontalSegment, from int) ([]HorizontalSegment, int) {
	sameHeightSegments := []HorizontalSegment{}
	sameHeightIndex := from
	sameHeightUntil := from + 1
	for ; sameHeightUntil < len(segments) && segments[sameHeightUntil].Line.StartPoint.I == segments[sameHeightIndex].Line.StartPoint.I; sameHeightUntil++ {
	}
	sameHeightSegments = append(sameHeightSegments, segments[sameHeightIndex:sameHeightUntil]...)
	return sameHeightSegments, sameHeightUntil
}

func zipHorizontalSegmentSlices(lhs []HorizontalSegment, rhs []HorizontalSegment) []HorizontalSegment {
	zippedHorizontals := []HorizontalSegment{}
	lhsIndex := 0
	rhsIndex := 0
	for lhsIndex < len(lhs) && rhsIndex < len(rhs) {
		if lhs[lhsIndex].Line.StartPoint.J <= rhs[rhsIndex].Line.StartPoint.J {
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

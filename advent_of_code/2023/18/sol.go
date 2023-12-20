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
	// for i := range matrix {
	// 	fmt.Printf("%3.v ", i)
	// 	for j := range matrix[i] {
	// 		fmt.Printf("%c", matrix[i][j])
	// 	}
	// 	fmt.Println()
	// }
	// fmt.Println()
	excavated := Excavate(matrix)
	for i := range excavated {
		fmt.Printf("%3.v ", i)
		for j := range excavated[i] {
			fmt.Printf("%c", excavated[i][j])
		}
		fmt.Println()
	}
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
// What we do is a scanline solution.
func PartTwo(data []DigStep) int {
	lineSegments := DigDirectionsToLineSegments(Pair{0, 0}, data, false)

	verticalSegments := []LineSegment{}
	for _, segment := range lineSegments {
		if segment.StartPoint.I == segment.LastPoint.I {
			continue
		}
		if segment.StartPoint.I < segment.LastPoint.I {
			verticalSegments = append(verticalSegments, segment)
			continue
		}
		verticalSegments = append(verticalSegments, LineSegment{segment.LastPoint, segment.StartPoint})
	}
	// prepping the data for the sweep-line
	sort.Slice(verticalSegments, func(lhs, rhs int) bool {
		return verticalSegments[lhs].StartPoint.I < verticalSegments[rhs].StartPoint.I
	})
	heightZero := verticalSegments[0].StartPoint.I
	for _, l := range verticalSegments {
		if l.StartPoint.I-heightZero == 13 || l.LastPoint.I-heightZero == 13 || l.StartPoint.I-heightZero == 12 {
			fmt.Printf("%v - %v (%v)\n", l.StartPoint.I-heightZero, l.LastPoint.I-heightZero, l.StartPoint.J)
		}
	}

	scanHeight := verticalSegments[0].StartPoint.I
	activeVerticals, i := updateActiveVerticals(verticalSegments, 0, scanHeight, []LineSegment{})
	fmt.Println(activeVerticals)
	totalExcavated := sumDistances(generateHorizontalDistances(activeVerticals))
	for len(activeVerticals) != 0 {
		prevHeight := scanHeight
		scanHeight = findNextChangeHeight(activeVerticals, scanHeight)
		fmt.Printf("Current scan height: %v\n", scanHeight)
		currentDistances := generateHorizontalDistances(activeVerticals)
		excavatedInRow := sumDistances(currentDistances)
		// fmt.Printf("Hop exc: %v\n", excavatedInRow*(scanHeight-prevHeight-1))
		totalExcavated += excavatedInRow * (scanHeight - prevHeight - 1)
		activeVerticals, i = updateActiveVerticals(verticalSegments, i, scanHeight, activeVerticals)
		fmt.Printf("Active vert: %v\n", activeVerticals)
		nextDistances := generateHorizontalDistances(activeVerticals)

		transitionRowDistances := mergeSumDistances(currentDistances, nextDistances)
		totalExcavated += transitionRowDistances
		fmt.Printf("i: %v - %v\n", scanHeight-heightZero, transitionRowDistances)
		fmt.Printf("Height: %v, PrevH: %v, Excavate in row: %v\n", scanHeight, prevHeight, transitionRowDistances)
		// fmt.Printf("Total exc: %v\n", totalExcavated)
	}

	return totalExcavated
}

// vertical lines have to appear in pairs, if not, we messed up
func generateHorizontalDistances(activeVerticals []LineSegment) []Pair {
	horizontalDistances := []Pair{}
	for i := 1; i < len(activeVerticals); i += 2 {
		horizontalDistances = append(horizontalDistances, Pair{activeVerticals[i-1].StartPoint.J, activeVerticals[i].StartPoint.J})
	}
	return horizontalDistances
}

// A bit tricky. We use the current distances as a mask and combine all ranges
// covered by it from next distances.  
func mergeSumDistances(currentDistances []Pair, nextDistances []Pair) int {
	combined := currentDistances
	combined = append(combined, nextDistances...)

	sort.Slice(combined, func(i, j int) bool {
		return combined[i].I < combined[j].I
	})

	merged := []Pair{}
	for i := 0; i < len(combined); {
		merged = append(merged, combined[i])
		last := len(merged) - 1
		j := i + 1
		for ; j < len(combined); j++ {
			if (merged[last].I <= combined[j].I && combined[j].I <= merged[last].J) || (merged[last].I <= combined[j].J && combined[j].J <= merged[last].J) {
				if merged[last].J < combined[j].J {
					merged[last].J = combined[j].J
				}
				continue
			}
			break
		}
		i = j
	}
	fmt.Println(merged)
	return sumDistances(merged)
}

func sumDistances(distances []Pair) int {
	sum := 0
	for _, distance := range distances {
		sum += 1 + distance.J - distance.I
	}
	return sum
}

func updateActiveVerticals(verticalSegments []LineSegment, i, scanHeight int, activeVerticals []LineSegment) ([]LineSegment, int) {
	nextActives := []LineSegment{}
	for ; i < len(verticalSegments); i++ {
		if verticalSegments[i].StartPoint.I > scanHeight {
			break
		}
		nextActives = append(nextActives, verticalSegments[i])
	}
	for _, activeSegment := range activeVerticals {
		if activeSegment.LastPoint.I > scanHeight {
			nextActives = append(nextActives, activeSegment)
		}
	}

	sort.Slice(nextActives, func(lhs, rhs int) bool {
		return nextActives[lhs].StartPoint.J < nextActives[rhs].StartPoint.J
	})
	return nextActives, i
}

func findNextChangeHeight(activeVerticals []LineSegment, previousHeight int) int {
	closest := activeVerticals[0].LastPoint.I
	for i := 1; i < len(activeVerticals); i++ {
		if closest-previousHeight > activeVerticals[i].LastPoint.I-previousHeight {
			closest = activeVerticals[i].LastPoint.I
		}
	}
	return closest
}

func Abs(val int) int {
	if val < 0 {
		return 0 - val
	}
	return val
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

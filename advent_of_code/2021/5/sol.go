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

type Point struct {
	X int
	Y int
}

type Line struct {
	A Point
	B Point
}

func MaxInt(lhs, rhs int) int {
	if lhs < rhs {
		return rhs
	}
	return lhs
}

func countOverlappingLinePoints(lines []Line) int {
	maxX := 0
	maxY := 0
	for _, line := range lines {
			maxX = MaxInt(maxX, line.A.X)
			maxX = MaxInt(maxX, line.B.X)
			maxY = MaxInt(maxY, line.A.Y)
			maxY = MaxInt(maxY, line.B.Y)
	}

	// we do it dum dum way, and just find the biggest numbers we have to handle
	// allocate a matrix, mark each visited point then count them up
	// not the smartest, but easy and works, good enough for this
	maxX += 1 // just so we actually fit the biggest value
	maxY += 1
	counter := make([][]int, maxY)
	for i := range counter {
		counter[i] = make([]int, maxX)
	}

	for _, line := range lines {
		markLine(counter, line)
	}

	overlapCount := 0
	for _, row := range counter {
		for _, value := range row {
			if value > 1 {
				overlapCount++
			}
		}
	}

	return overlapCount
}

func markLine(counter [][]int, line Line) {
	offsetI := line.B.Y - line.A.Y
	if offsetI != 0 {
		offsetI = offsetI / int(math.Abs(float64(offsetI)))
	}
	offsetJ := line.B.X - line.A.X
	if offsetJ != 0 {
		offsetJ = offsetJ / int(math.Abs(float64(offsetJ)))
	}
	i := line.A.Y
	j := line.A.X
	for {
		counter[i][j] += 1
		if i == line.B.Y && j == line.B.X {
			break
		}
		i += offsetI
		j += offsetJ
	}
}

func main() {
	re := regexp.MustCompile("(\\d+),(\\d+) -> (\\d+),(\\d+)")
	scanner := bufio.NewScanner(os.Stdin)

	lines := []Line{}
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) != 5 {
			log.Fatalf("Invalid input received: %v\n", line)
		}
		numbers := make([]int, 4)
		var err error
		for i, numText := range matches[1:] {
			numbers[i], err = strconv.Atoi(numText)
			if err != nil {
				log.Fatalf("Failed to convert %v to int.\n", numText)
			}
		}
		pointA := Point{numbers[0], numbers[1]}
		pointB := Point{numbers[2], numbers[3]}

		lines = append(lines, Line{pointA, pointB})
	}

	filteredLines := make([]Line, 0, len(lines))
	for _, line := range lines {
		if line.A.X == line.B.X || line.A.Y == line.B.Y {
			filteredLines = append(filteredLines, line)
		}
	}
	fmt.Println("Part one")
	fmt.Printf("Number of overlapping horizontal/vertical line points: %v\n", countOverlappingLinePoints(filteredLines))
	fmt.Println("Part two")
	fmt.Printf("Total number of overlapping line points: %v\n", countOverlappingLinePoints(lines))
}

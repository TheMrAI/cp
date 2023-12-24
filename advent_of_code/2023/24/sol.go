package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Triplet struct {
	X int
	Y int
	Z int
}

type Hail struct {
	Pos Triplet
	Vel Triplet
}

func ParseInput(lines []string) []Hail {
	hail := []Hail{}

	numberRegexp := regexp.MustCompile("-?\\d+")

	for _, line := range lines {
		parts := strings.Split(line, "@")
		if len(parts) != 2 {
			log.Fatalf("Unexpected number of sections in input line: '%v'\n", line)
		}
		position := parseNumbers(parts[0], numberRegexp)
		velocity := parseNumbers(parts[1], numberRegexp)
		if len(position) != 3 && len(velocity) != 3 {
			log.Fatalf("Unexpected number of position/velocity data: %v, %v\n", position, velocity)
		}
		hail = append(hail, Hail{Triplet{position[0], position[1], position[2]}, Triplet{velocity[0], velocity[1], velocity[2]}})
	}

	return hail
}

func parseNumbers(input string, numberRegexp *regexp.Regexp) []int {
	numberSections := numberRegexp.FindAllString(input, -1)
	numbers := []int{}
	for _, text := range numberSections {
		number, err := strconv.Atoi(text)
		if err != nil {
			fmt.Printf("Integer conversion failed on input: '%v'\n", text)
		}
		numbers = append(numbers, number)
	}
	return numbers
}

type LineEq struct {
	M float64
	C float64
}

func hailReachPointInFuture(hail Hail, x float64, y float64) bool {
	xDir := (x - float64(hail.Pos.X)) / float64(hail.Vel.X)
	yDir := (y - float64(hail.Pos.Y)) / float64(hail.Vel.Y)
	return xDir > 0 && yDir > 0
}

func PartOne(hail []Hail, rangeStart, rangeLast float64) int {
	lineEquations := make([]LineEq, len(hail))
	for i := range hail {
		lineEquations[i] = hailToLineEquation(hail[i])
	}

	collisionCountInFuture := 0
	for i := 0; i < len(lineEquations)-1; i++ {
		for j := i + 1; j < len(lineEquations); j++ {
			x, y := intersectionPoint(lineEquations[i], lineEquations[j])
			// parallel lines
			if (x == math.Inf(-1) || x == math.Inf(1)) && (y == math.Inf(-1) || y == math.Inf(1)) {
				continue
			}
			if rangeStart <= x && x <= rangeLast && rangeStart <= y && y <= rangeLast && hailReachPointInFuture(hail[i], x, y) && hailReachPointInFuture(hail[j], x, y) {
				collisionCountInFuture++
			}
		}
	}
	return collisionCountInFuture
}

// only using X, Y coordinates
func hailToLineEquation(hail Hail) LineEq {
	m := float64(hail.Vel.Y) / float64(hail.Vel.X)
	c := float64(hail.Pos.Y) - m*float64(hail.Pos.X)
	return LineEq{m, c}
}

func intersectionPoint(a, b LineEq) (float64, float64) {
	x := (b.C - a.C) / (a.M - b.M)
	y := a.M*x + a.C
	return x, y
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	hail := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("The number of hailstones that will collide between 200000000000000 - 400000000000000 based on their X,Y coordinates: %v\n", PartOne(hail, 200000000000000, 400000000000000))
}

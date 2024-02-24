package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func ParseInput(line string) (int, int, int, int) {
	numberRegexp := regexp.MustCompile(`target area: x=(-?\d+)..(-?\d+), y=(-?\d+)..(-?\d+)`)
	numbersText := numberRegexp.FindStringSubmatch(line)
	if len(numbersText) != 5 {
		log.Fatalf("Failed to parse coordinates. Expected to find 4 integers in the input: '%v'", line)
	}
	xMin, err := strconv.Atoi(numbersText[1])
	if err != nil {
		log.Fatalf("Failed to parse integer, from input: '%v'", numbersText[1])
	}
	xMax, err := strconv.Atoi(numbersText[2])
	if err != nil {
		log.Fatalf("Failed to parse integer, from input: '%v'", numbersText[2])
	}
	yMin, err := strconv.Atoi(numbersText[3])
	if err != nil {
		log.Fatalf("Failed to parse integer, from input: '%v'", numbersText[3])
	}
	yMax, err := strconv.Atoi(numbersText[4])
	if err != nil {
		log.Fatalf("Failed to parse integer, from input: '%v'", numbersText[4])
	}

	return xMin, xMax, yMin, yMax
}

func PartOne(yMin int) int {
	n := yMin
	if n < 0 {
		n = -n
	}
	return n * (n - 1) / 2
}

type XY struct {
	X int
	Y int
}

// the calculator assumes that the target area is always at some positive x
// and some negative y
// Note to self: I am such a dumb dumb. Spent around 2 hours trying to figure out
// some equation which could tell me when we would hit the zone and what limits there would
// be, when I realized oh, I would have to calculate with floating points, so I can't really
// tell if the starting velocity was an exact integer.
// So just implemented brute force in 2 minutes, and the whole thing executes:
// real    0m0.073s
// user    0m0.096s
// sys     0m0.042s
// Try not to be too smart man!
func PartTwo(xMin, xMax, yMin, yMax int) []XY {
	startingVelocities := []XY{}
	// just brute force it
	for yS := yMin; yS <= 0-yMin; yS++ {
		for xS := 0; xS <= xMax; xS++ {
			if !wouldHitPosition(xS, yS, xMin, xMax, yMin, yMax) {
				continue
			}
			startingVelocities = append(startingVelocities, XY{xS, yS})
		}
	}
	return startingVelocities
}

func wouldHitPosition(xS, yS, xMin, xMax, yMin, yMax int) bool {
	xCoord := 0
	yCoord := 0
	for yCoord >= yMin {
		xCoord += xS
		yCoord += yS
		if xMin <= xCoord && xCoord <= xMax && yMin <= yCoord && yCoord <= yMax {
			return true
		}
		if xS > 0 {
			xS -= 1
		}
		yS -= 1
	}
	return false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	line := " "
	for scanner.Scan() {
		line = scanner.Text()
	}

	xMin, xMax, yMin, yMax := ParseInput(line)
	fmt.Printf("The highest achievable y trajectory point is: %v\n", PartOne(yMin))
	startingVelocities := PartTwo(xMin, xMax, yMin, yMax)
	fmt.Printf("The number of valid initial velocities: %v\n", len(startingVelocities))
	// fmt.Println(startingVelocities)
}

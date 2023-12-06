package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ParseInput(line string) []int {
	numbers := []int{}

	numberRegexp := regexp.MustCompile("\\d+")

	for _, numberText := range numberRegexp.FindAllString(line, -1) {
		number, err := strconv.Atoi(numberText)
		if err != nil {
			log.Fatalf("Failed number conversion on input: %v\n", numberText)
		}
		numbers = append(numbers, number)
	}

	return numbers
}

func PartOne(times, distances []int) int {
	multiple := 1
	for i := range times {
		multiple *= CountWinningOptions(times[i], distances[i])
	}
	return multiple
}

// Without binary search:
// real    0m0.017s
// user    0m0.018s
// sys     0m0.000s
// With binary search:
// real    0m0.002s
// user    0m0.002s
// sys     0m0.000s
func CountWinningOptions(time, distance int) int {
	lowerBound := sort.Search(time+1, func(revTime int) bool {
		return Win(revTime, time, distance)
	})
	delta := time + 1 - lowerBound
	upperBound := sort.Search(delta, func(t int) bool {
		return !Win(t+delta+1, time, distance)
	})
	return (upperBound + delta) - lowerBound
}

func Win(revTime, totalTime, distance int) bool {
	speed := revTime
	remainingTime := totalTime - revTime
	finalDistance := speed * remainingTime
	return finalDistance > distance
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if len(lines) != 2 {
		log.Fatalf("Expected only 2 input lines. Got: %v\n", len(lines))
	}

	timesLine, ok := strings.CutPrefix(lines[0], "Time:")
	if !ok {
		log.Fatalf("Expected 'Times:' section. Got: %v\n", lines[0])
	}
	times := ParseInput(timesLine)

	distanceLine, ok := strings.CutPrefix(lines[1], "Distance:")
	if !ok {
		log.Fatalf("Expected 'Distance:' section. Got: %v\n", lines[0])
	}
	distances := ParseInput(distanceLine)

	fmt.Println("Part one")
	fmt.Printf("Multiple of number of options for winning: %v\n", PartOne(times, distances))

	fmt.Println("Part two")
	combinedTimeText := strings.ReplaceAll(timesLine, " ", "")
	combinedTime := ParseInput(combinedTimeText)

	combinedDistanceText := strings.ReplaceAll(distanceLine, " ", "")
	combinedDistance := ParseInput(combinedDistanceText)

	fmt.Printf("Multiple of number of options for winning: %v\n", CountWinningOptions(combinedTime[0], combinedDistance[0]))
}

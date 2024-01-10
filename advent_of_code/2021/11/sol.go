package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func ParseInput(lines []string) [][]int {
	data := make([][]int, len(lines))

	numberRegexp := regexp.MustCompile("\\d")

	for i, line := range lines {
		data[i] = make([]int, len(line))
		numbers := numberRegexp.FindAllString(line, -1)
		if len(numbers) == 0 {
			log.Fatalf("Unexpected input: '%v'\n", line)
		}
		for j := range numbers {
			val, err := strconv.Atoi(numbers[j])
			// I really wonder, why am I writing these branches. It should be utterly impossible to
			// get to this point. The regex only parses out single digits.
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'\n", numbers[j])
			}
			data[i][j] = val
		}
	}

	return data
}

func PartOne(matrix [][]int, simulationCount int) int {
	// make a quick copy so we don't destroy the original
	workingData := make([][]int, len(matrix))
	for i := range matrix {
		workingData[i] = make([]int, len(matrix[i]))
		copy(workingData[i], matrix[i])
	}

	flashCount := 0
	for round := 0; round < simulationCount; round++ {
		flashCount += SimulateRound(workingData)
	}

	return flashCount
}

func PartTwo(matrix [][]int) int {
	// make a quick copy so we don't destroy the original
	workingData := make([][]int, len(matrix))
	for i := range matrix {
		workingData[i] = make([]int, len(matrix[i]))
		copy(workingData[i], matrix[i])
	}

	flashCount := 0
	round := 1
	for ; ; round++ {
		flashCount = SimulateRound(workingData)
		if flashCount == len(matrix)*len(matrix[0]) {
			return round
		}
	}
}

type Pair struct {
	I int
	J int
}

func SimulateRound(matrix [][]int) int {
	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] += 1
		}
	}

	flashCount := 0
	flashers := []Pair{}
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] != 10 {
				continue
			}
			flashers = append(flashers, Pair{i, j})
		}
	}

	deltas := []Pair{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	for len(flashers) != 0 {
		toProcess := flashers[0]
		flashers = flashers[1:]

		if matrix[toProcess.I][toProcess.J] < 10 {
			continue
		}
		flashCount++
		matrix[toProcess.I][toProcess.J] = 0

		for dI := range deltas {
			location := Pair{toProcess.I + deltas[dI].I, toProcess.J + deltas[dI].J}
			if location.I < 0 || location.I >= len(matrix) || location.J < 0 || location.J >= len(matrix[0]) || matrix[location.I][location.J] == 0 {
				continue
			}
			matrix[location.I][location.J]++
			flashers = append(flashers, location)
		}
	}
	return flashCount
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)
	fmt.Println(data)

	fmt.Println("Part one")
	fmt.Printf("The total number of flashes after 100 steps: %v\n", PartOne(data, 100))
	fmt.Println("Part two")
	fmt.Printf("The step where all octopuses first synchronize: %v\n", PartTwo(data))
}

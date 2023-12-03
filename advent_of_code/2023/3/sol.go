package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type Location struct {
	I int
	J int
}

func SymbolAround(lines []string, i, j int) (bool, Location) {
	offsets := []Location{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	for _, offset := range offsets {
		candidate := Location{i + offset.I, j + offset.J}
		if candidate.I < 0 || candidate.I >= len(lines) || candidate.J < 0 || candidate.J >= len(lines[0]) {
			continue
		}
		r := rune(lines[candidate.I][candidate.J])
		if !unicode.IsDigit(r) && r != '.' {
			return true, candidate
		}
	}
	return false, Location{-1, -1}
}

// we assume that a part number can only have at most one symbol attachment
// if this wasn't the case we could extend with a map or slice (but it was the case)
type PartNumber struct {
	Value          int
	SymbolLocation Location
}

func ParsePartNumbers(lines []string) []PartNumber {
	partNumbers := []PartNumber{}

	for i, line := range lines {
		value := 0
		symbolAt := Location{-1, -1}
		for j, r := range line {
			if unicode.IsDigit(r) {
				value *= 10
				value += int(r - '0')
				symbolFound, at := SymbolAround(lines, i, j)
				if symbolFound {
					symbolAt = at
				}
				continue
			}
			if symbolAt.I != -1 && symbolAt.J != -1 {
				partNumbers = append(partNumbers, PartNumber{value, symbolAt})
			}
			value = 0
			symbolAt = Location{-1, -1}
		}
		if symbolAt.I != -1 && symbolAt.J != -1 {
			partNumbers = append(partNumbers, PartNumber{value, symbolAt})
		}
	}

	return partNumbers
}

func PartOne(partNumbers []PartNumber) int {
	partNumberSum := 0
	for _, partNumber := range partNumbers {
		partNumberSum += partNumber.Value
	}
	return partNumberSum
}

func PartTwo(partNumbers []PartNumber) int {
	gearRatioSum := 0

	ratioParts := map[Location]int{}
	for _, partNumber := range partNumbers {
		ratioPart, ok := ratioParts[partNumber.SymbolLocation]
		if ok {
			gearRatioSum += ratioPart * partNumber.Value
		}
		ratioParts[partNumber.SymbolLocation] = partNumber.Value
	}

	return gearRatioSum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	partNumbers := ParsePartNumbers(lines)

	fmt.Println("Part one")
	fmt.Printf("Sum of all part numbers: %v\n", PartOne(partNumbers))
	fmt.Println("Part two")
	fmt.Printf("Sum of all gear ratios: %v\n", PartTwo(partNumbers))
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func ExpandUniverse(universe []string, emptyColumns map[int]struct{}) []string {
	expandedUniverse := []string{}

	for _, line := range universe {
		lineEmpty := true
		for _, r := range line {
			lineEmpty = lineEmpty && (r == '.')
		}
		expandedLine := ExpandLine(line, emptyColumns)
		expandedUniverse = append(expandedUniverse, expandedLine)
		if lineEmpty {
			expandedUniverse = append(expandedUniverse, expandedLine)
		}
	}

	return expandedUniverse
}

func IdentifyEmptyColumns(universe []string) map[int]struct{} {
	emptyColumns := map[int]struct{}{}
	for j := range universe[0] {
		allEmpty := true
		for i := range universe {
			allEmpty = allEmpty && (universe[i][j] == '.')
		}
		if allEmpty {
			emptyColumns[j] = struct{}{}
		}
	}
	return emptyColumns
}

func ExpandLine(line string, emptyColumns map[int]struct{}) string {
	var b bytes.Buffer
	for i, r := range line {
		b.WriteRune(r)
		_, double := emptyColumns[i]
		if double {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type Pair struct {
	I int
	J int
}

func PartOne(universe []string) int {
	galaxies := []Pair{}
	for i := range universe {
		for j := range universe[0] {
			if universe[i][j] == '#' {
				galaxies = append(galaxies, Pair{i, j})
			}
		}
	}
	return SumManhattanDistances(galaxies)
}

func SumManhattanDistances(galaxies []Pair) int {
	distanceTotal := 0
	for i := 0; i < len(galaxies)-1; i++ {
		for j := i + 1; j < len(galaxies); j++ {
			deltaI := abs(galaxies[i].I - galaxies[j].I)
			deltaJ := abs(galaxies[i].J - galaxies[j].J)
			distanceTotal += deltaI + deltaJ
		}
	}

	return distanceTotal
}

func PartTwo(universe []string, emptyColumns map[int]struct{}, expansionRate int) int {
	galaxies := []Pair{}

	expandedRowsCount := 0
	for i, line := range universe {
		lineEmpty := true
		for _, r := range line {
			lineEmpty = lineEmpty && (r == '.')
		}
		if lineEmpty {
			expandedRowsCount++
			continue
		}
		expandedColumnCount := 0
		for j := range universe[0] {
			_, emptyColumn := emptyColumns[j]
			if emptyColumn {
				expandedColumnCount++
				continue
			}
			if universe[i][j] == '#' {
				additionalRowDistance := expandedRowsCount * expansionRate
				additionalColumnDistance := expandedColumnCount * expansionRate
				galaxies = append(galaxies, Pair{i - expandedRowsCount + additionalRowDistance, j - expandedColumnCount + additionalColumnDistance})
			}
		}
	}

	return SumManhattanDistances(galaxies)
}

func abs(value int) int {
	if value < 0 {
		return 0 - value
	}
	return value
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	universe := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		universe = append(universe, line)
	}

	emptyColumns := IdentifyEmptyColumns(universe)
	expandedUniverse := ExpandUniverse(universe, emptyColumns)

	fmt.Println("Part one")
	fmt.Printf("Sum of lengths between all pairs of galaxies: %v\n", PartOne(expandedUniverse))
	fmt.Println("Part two")
	fmt.Printf("Sum of lengths between all pairs of galaxies with huge expansion: %v\n", PartTwo(universe, emptyColumns, 1_000_000))
}

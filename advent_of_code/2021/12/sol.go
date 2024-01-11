package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type Cave struct {
	Big       bool
	Neighbors []string
}

func ParseInput(lines []string) map[string]Cave {
	graph := map[string]Cave{}
	for _, line := range lines {
		// you know what
		// this time we won't even regexp, just split on '-'
		// like a peasant and hope for the best outcome :D
		caveNames := strings.Split(line, "-")
		if len(caveNames) != 2 {
			log.Fatalf("Unexpected number of caves in the input: '%v'\n", line)
		}
		source := insertIfMissing(graph, caveNames[0])
		target := insertIfMissing(graph, caveNames[1])

		// huh it seems map values are immutable in Go
		// can't modify Neighbors directly
		updatedSource := append(source.Neighbors, caveNames[1])
		graph[caveNames[0]] = Cave{source.Big, updatedSource}
		updatedTarget := append(target.Neighbors, caveNames[0])
		graph[caveNames[1]] = Cave{target.Big, updatedTarget}
	}

	return graph
}

func insertIfMissing(system map[string]Cave, caveName string) Cave {
	entry, ok := system[caveName]
	if !ok {
		bigCave := true
		for _, c := range caveName {
			if !unicode.IsUpper(c) {
				bigCave = false
				break
			}
		}
		entry = Cave{bigCave, []string{}}
		system[caveName] = entry
	}
	return entry
}

func PartOne(caveSystem map[string]Cave, source string, target string) int {
	visited := map[string]int{}
	return DFTWalk(caveSystem, false, source, target, visited)
}

func PartTwo(caveSystem map[string]Cave, source string, target string) int {
	visited := map[string]int{}
	return DFTWalk(caveSystem, true, source, target, visited)
}

func DFTWalk(caveSystem map[string]Cave, smallVisitAvailable bool, source string, target string, visited map[string]int) int {
	if source == target {
		return 1
	}

	times, caveVisited := visited[source]
	if caveVisited && times > 0 {
		if source == "start" || !smallVisitAvailable {
			return 0
		}
		smallVisitAvailable = false
	}

	if !caveSystem[source].Big {
		// only tracking small caves
		visited[source] += 1
	}

	distinctPathCount := 0
	for _, neighbor := range caveSystem[source].Neighbors {
		distinctPathCount += DFTWalk(caveSystem, smallVisitAvailable, neighbor, target, visited)
	}

	if !caveSystem[source].Big {
		visited[source] -= 1
	}

	return distinctPathCount
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	graph := ParseInput(lines)
	for name, entry := range graph {
		fmt.Printf("Name: %v\n", name)
		fmt.Printf("\t%v\n", entry.Big)
		fmt.Printf("\t%v\n", entry.Neighbors)
	}

	fmt.Println("Part one")
	fmt.Printf("The number of possible paths: %v\n", PartOne(graph, "start", "end"))
	fmt.Println("Part two")
	fmt.Printf("The number of possible paths, with 1 small cave double visit: %v\n", PartTwo(graph, "start", "end"))
}

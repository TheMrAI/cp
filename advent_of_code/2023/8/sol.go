package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Side struct {
	Left  string
	Right string
}

func ParseInput(lines []string) ([]rune, map[string]Side) {
	charRegexp := regexp.MustCompile("\\w")
	routeRegexp := regexp.MustCompile("(\\w{3}) = \\((\\w{3}), (\\w{3})\\)")

	directions := []rune{}
	directionBytes := charRegexp.FindAll([]byte(lines[0]), -1)
	if len(directionBytes) == 0 {
		log.Fatalf("Failed to parse directions on input: %v\n", lines[0])
	}
	for _, directionByte := range directionBytes {
		if len(directionByte) != 1 {
			log.Fatalf("Expected a single byte, received: %v\n", directionByte)
		}
		directions = append(directions, rune(directionByte[0]))
	}

	routes := map[string]Side{}
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		routeTexts := routeRegexp.FindStringSubmatch(line)
		if len(routeTexts) != 4 {
			log.Fatalf("Failed to parse route on input: %v\n", line)
		}
		routes[routeTexts[1]] = Side{routeTexts[2], routeTexts[3]}
	}

	return directions, routes
}

func PartOne(directions []rune, routes map[string]Side, start, target string) int {
	stepCount := 0
	current := start
	_, ok := routes[start]
	if !ok {
		fmt.Printf("Map does not contain starting point: %v\n", start)
		return -1
	}
	_, ok = routes[target]
	if !ok {
		fmt.Printf("Map does not contain target point: %v\n", target)
		return -1
	}

	directionIndex := 0
	for current != target {
		// log.Println(current)
		switch directions[directionIndex] {
		case 'L':
			{
				current = routes[current].Left
				break
			}
		case 'R':
			{
				current = routes[current].Right
				break
			}
		default:
			{
				log.Fatalf("Invalid directions received: %v\n", directions[directionIndex])
			}
		}
		stepCount++
		directionIndex++
		directionIndex = directionIndex % len(directions)
	}

	return stepCount
}

// A number of tricks have been hidden in this one.
// Firstly the result is so huge, it takes very long to calculate and verify, so
// there must be a pattern. Alright, the only pattern that can be here is that after finding
// a 'kjZ' ending there will be a stable number of steps after which we reach another 'kjZ'. Which
// 'kjZ' exactly doesn't matter, just that the cycle length stays the same.
// Secondly, and I think this is very important, for all starting nodes each has to be at the same
// point in their cycles, one after 'kjZ'. If they were out of order in any way there would be no hope.
// With these two assumptions, we just have to count until the first 'kjZ' ending node for each starting node.
// Take the least common multiple of these step counts and that is our solution.
//
// Ps.: I had the idea there needed to be a pattern, but didn't understand the requirements for that.
// Printed out after how many steps would each node reach 'kjZ' and saw the pattern. Inferred the rest only after.
func PartTwo(directions []rune, routes map[string]Side) int {
	stepCount := 0
	directionIndex := 0

	nodes := []string{}
	for key := range routes {
		if key[len(key)-1] == 'A' {
			nodes = append(nodes, key)
		}
	}

	cycleLengths := make([]int, len(nodes))
	cyclesFound := make([]bool, len(nodes))

	for !AllEndingInZ(nodes) {
		for i, current := range nodes {
			switch directions[directionIndex] {
			case 'L':
				{
					nodes[i] = routes[current].Left
					break
				}
			case 'R':
				{
					nodes[i] = routes[current].Right
					break
				}
			default:
				{
					log.Fatalf("Invalid directions received: %v\n", directions[directionIndex])
				}
			}

			if nodes[i][len(nodes[i])-1] == 'Z' {
				cyclesFound[i] = true
				// plus one so we get the length of the cycle
				cycleLengths[i]++
			}
			if cyclesFound[i] {
				continue
			}
			cycleLengths[i]++
		}

		allCycleLengthsFound := true
		for i := range cyclesFound {
			if !cyclesFound[i] {
				allCycleLengthsFound = false
				break
			}
		}
		if allCycleLengthsFound {
			lcm := 1
			for _, cycleLength := range cycleLengths {
				lcm = EuclidLCM(lcm, cycleLength)
			}
			return lcm
		}

		stepCount++
		directionIndex++
		directionIndex = directionIndex % len(directions)
	}

	return stepCount
}

func EuclidLCM(a, b int) int {
	return a / EuclidGCD(a, b) * b
}

func EuclidGCD(a, b int) int {
	if b > a {
		a, b = b, a
	}
	if b == 0 {
		return a
	}
	for b != 0 {
		a = a % b
		a, b = b, a
	}
	return a
}

func AllEndingInZ(nodes []string) bool {
	for _, node := range nodes {
		if node[len(node)-1] != 'Z' {
			return false
		}
	}
	return true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	directions, routes := ParseInput(lines)

	fmt.Println("Part one")
	stepCount := PartOne(directions, routes, "AAA", "ZZZ")
	if stepCount != -1 {
		fmt.Printf("Steps required for route AAA -> ZZZ: %v\n", stepCount)
	}
	fmt.Println("Part two")
	fmt.Printf("Steps required for simultaneous route xyA -> kjZ: %v\n", PartTwo(directions, routes))
}

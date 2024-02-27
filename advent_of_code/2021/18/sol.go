package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type SnailNum struct {
	Value int
	Depth int
}

func Add(lhs []SnailNum, rhs []SnailNum) []SnailNum {
	result := make([]SnailNum, 0, len(lhs)+len(rhs))
	result = append(result, lhs...)
	result = append(result, rhs...)
	for i := range result {
		result[i].Depth += 1
	}
	return result
}

func Reduce(snail []SnailNum) []SnailNum {
	// fmt.Printf("Reducing: %v\n", snail)
	result := make([]SnailNum, 0, len(snail))
	result = append(result, snail...)

	changed := true
	// Sheeeesh! This is the place where I have messed up! Didn't read the instructions
	// carefully enough, and didn't realize that we want to explode all we can, then split
	// then start from the beginning if we had to split.
reduxList:
	for changed {
		changed = false
		for i := 0; i < len(result); i++ {
			if result[i].Depth >= 4 {
				result = explode(result, i)
				changed = true
			}
		}
		for i := 0; i < len(result); i++ {
			if result[i].Value >= 10 {
				result = split(result, i)
				changed = true
				continue reduxList
			}
		}
	}

	return result
}

func explode(snail []SnailNum, start int) []SnailNum {
	// fmt.Println("\t", "Exploding...")
	exploded := make([]SnailNum, 0, len(snail))
	exploded = append(exploded, snail[0:start]...)
	if len(exploded) >= 1 {
		exploded[len(exploded)-1].Value += snail[start].Value
	}
	exploded = append(exploded, SnailNum{0, snail[start].Depth - 1})
	if start+2 < len(snail) {
		exploded = append(exploded, snail[start+2:]...)
		exploded[start+1].Value += snail[start+1].Value
	}
	// fmt.Println("\t", exploded)
	return exploded
}

func split(snail []SnailNum, start int) []SnailNum {
	// fmt.Println("\tSplitting...")
	afterSplit := make([]SnailNum, 0, len(snail)+1)
	afterSplit = append(afterSplit, snail[0:start]...)
	afterSplit = append(afterSplit, SnailNum{snail[start].Value / 2, snail[start].Depth + 1}, SnailNum{Ceil(snail[start].Value, 2), snail[start].Depth + 1})
	afterSplit = append(afterSplit, snail[start+1:]...)
	// fmt.Println("\t", afterSplit)
	return afterSplit
}

func Magnitude(snail []SnailNum) int {
	lowestLevel := 0
	for i := range snail {
		if snail[i].Depth > lowestLevel {
			lowestLevel = snail[i].Depth
		}
	}
	return magnitude(snail, lowestLevel)
}

func magnitude(snail []SnailNum, currentDepth int) int {
	working := make([]SnailNum, 0, len(snail))
	working = append(working, snail...)
	for len(working) != 1 {
		nextWorking := make([]SnailNum, 0, len(working)) // allocate a bit much, but that is okay
		for i := 0; i < len(working); {
			// no need to check, whether or not there is a next element, because we shouldn't reach such a state
			if working[i].Depth == currentDepth {
				nextWorking = append(nextWorking, SnailNum{(working[i].Value * 3) + (working[i+1].Value * 2), working[i].Depth - 1})
				i += 2
			} else {
				nextWorking = append(nextWorking, working[i])
				i += 1
			}
		}
		working = nextWorking
		currentDepth--
	}
	return working[0].Value
}

func Ceil(a, b int) int {
	if a == 0 {
		return 0
	}
	return ((a - 1) / b) + 1
}

// Instead of building an unwieldy tree, we just simply
// flatten the numbers remembering their depth and value.
// This way add, explode and split operations will take ~O(N) time.
// Magnitude calculations are a bit trickier, but those only have to be
// done once.
func ParseInput(lines []string) [][]SnailNum {
	snailNums := [][]SnailNum{}

	for _, line := range lines {
		snailNums = append(snailNums, ParseSnailNum(line))
	}

	return snailNums
}

func ParseSnailNum(line string) []SnailNum {
	snails := []SnailNum{}
	depth := -1
	i := 0
	for i < len(line) {
		if line[i] == '[' {
			depth++
			i++
			continue
		} else if line[i] == ']' {
			depth--
			i++
			continue
		} else if line[i] == ',' {
			i++
			continue
		}
		// parse number and step i
		j := i + 1
		for ; j < len(line); j++ {
			if !(line[i] == '[' || line[i] == ']' || line[i] == ',') {
				break
			}
		}
		num, err := strconv.Atoi(line[i:j])
		if err != nil {
			log.Fatalf("Integer conversion failed on input: '%v'\n", num)
		}
		snails = append(snails, SnailNum{num, depth})
		i = j
	}
	return snails
}

func AdditionReductionSequence(snailNums [][]SnailNum) []SnailNum {
	snailSum := Reduce(snailNums[0])
	for i := 1; i < len(snailNums); i++ {
		snailSum = Reduce(Add(snailSum, Reduce(snailNums[i])))
		// fmt.Println(snailSum)
	}
	return snailSum
}

func PartOne(snailNums [][]SnailNum) int {
	return Magnitude(AdditionReductionSequence(snailNums))
}

func PartTwo(snailNums [][]SnailNum) int {
	largestMagnitude := 0
	for i := 0; i < len(snailNums)-1; i++ {
		for j := i + 1; j < len(snailNums); j++ {
			magnitudeLR := Magnitude(AdditionReductionSequence([][]SnailNum{snailNums[i], snailNums[j]}))
			if magnitudeLR > largestMagnitude {
				largestMagnitude = magnitudeLR
			}
			magnitudeRL := Magnitude(AdditionReductionSequence([][]SnailNum{snailNums[j], snailNums[i]}))
			if magnitudeRL > largestMagnitude {
				largestMagnitude = magnitudeRL
			}
		}
	}
	return largestMagnitude
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	snailNumbers := ParseInput(lines)
	fmt.Printf("The magnitude of the final sum: %v\n", PartOne(snailNumbers))
	fmt.Printf("The largest magnitude of any two distinct snailfish numbers is: %v\n", PartTwo(snailNumbers))
}

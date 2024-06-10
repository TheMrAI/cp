package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

type Instruction struct {
	Command string
	RegA    byte
	RegB    interface{} // this is very disturbing
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func ParseInput(lines []string) []Instruction {
	instructions := []Instruction{}

	singleArgInstruction := regexp.MustCompile(`(\w+) (\w)`)
	doubleArgInstruction := regexp.MustCompile(`(\w+) (\w) (-?\d+|\w+)`)

	for _, line := range lines {
		command := doubleArgInstruction.FindStringSubmatch(line)
		if len(command) != 0 {
			if IsLetter(command[3]) && len(command[3]) == 1 {
				value := byte(command[3][0])
				instructions = append(instructions, Instruction{command[1], command[2][0], value})
			} else {
				value, err := strconv.Atoi(command[3])
				if err != nil {
					log.Fatalf("String to integer conversion failed on input: '%s'", command[3])
				}
				instructions = append(instructions, Instruction{command[1], command[2][0], value})
			}
			continue
		}
		command = singleArgInstruction.FindStringSubmatch(line)
		instructions = append(instructions, Instruction{command[1], command[2][0], -1})
	}

	return instructions
}

func splitInstructionsIntoAluSections(instructions []Instruction) [][]Instruction {
	aluSections := [][]Instruction{}
	lastInpI := 0
	for i := 1; i < len(instructions); i++ {
		// Go, why no decent enums? WHy?
		if instructions[i].Command == "inp" {
			aluSections = append(aluSections, instructions[lastInpI:i])
			lastInpI = i
		}
	}
	aluSections = append(aluSections, instructions[lastInpI:])
	return aluSections
}

type SectionInput struct {
	State [4]int
	Value int
}

// Initially I have notices that there are sections in the instructions which are the same for all
// ALU instruction blocks. Noticed that all blocks contain the same amount of commands.
// But was lazy and didn't want to check further and just assumed that even if I had
// 9^14 options to check (an astoundingly huge number), it would be okay if the stages where cached on input.
// I was wrong, it will never finish as there are still too many options to check.
// Specifically it runs out of memory because of the caching. (64 GiBs)
// Leaving all this and its supporting code in to highlight the folly in not pursuing
// an observation until it has been exhausted.
func PartOneBruteForce(instructionSections [][]Instruction) int {
	maximumSerialNumber := 0
	sectionCache := make([]map[SectionInput]int, len(instructionSections))
	for i := range sectionCache {
		sectionCache[i] = map[SectionInput]int{}
	}

	findBiggestNumber(instructionSections[4:14], 0, [4]int{0, 0, 0, 0}, sectionCache, &maximumSerialNumber)
	// for i := range sectionCache {
	// 	fmt.Printf("%v\n", sectionCache[i])
	// }
	return maximumSerialNumber
}

// Yeah Go is garbage, don't make anybody write common functions, just because they are trivial to write....
func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func findBiggestNumber(instructionSections [][]Instruction, depth int, state [4]int, sectionCache []map[SectionInput]int, maximumSerialNumber *int) bool {
	if state[3] > 900_000_000_000 {
		return false
	}
	if depth >= len(instructionSections) {
		// return false
		return getValue(byte('z'), state[:]) == 0
	}

	for value := 9; value >= 1; value-- {
		sectionInput := SectionInput{state, value}
		_, seen := sectionCache[depth][sectionInput]
		if seen {
			sectionCache[depth][sectionInput] += 1
			continue
		}
		state = simulateAluSection(instructionSections[depth], sectionInput)
		*maximumSerialNumber = (*maximumSerialNumber * 10) + value
		found := findBiggestNumber(instructionSections, depth+1, state, sectionCache, maximumSerialNumber)
		if found {
			return found
		}
		*maximumSerialNumber = (*maximumSerialNumber - value) / 10

		sectionCache[depth][sectionInput] = 1
	}
	return false
}

func simulateAluSection(instructions []Instruction, sectionInput SectionInput) [4]int {
	for i := range instructions {
		simulateAluInstruction(instructions[i], sectionInput.State[:], sectionInput.Value)
	}
	return sectionInput.State
}

func simulateAluInstruction(instruction Instruction, state []int, input int) {
	registerInd := registerToIndex(instruction.RegA)

	switch instruction.Command {
	case "inp":
		{
			state[registerInd] = input
		}
	case "add":
		{
			state[registerInd] = state[registerInd] + getValue(instruction.RegB, state)
		}
	case "mul":
		{
			state[registerInd] = state[registerInd] * getValue(instruction.RegB, state)
		}
	case "div":
		{
			state[registerInd] = state[registerInd] / getValue(instruction.RegB, state)
		}
	case "mod":
		{
			state[registerInd] = state[registerInd] % getValue(instruction.RegB, state)
		}
	case "eql":
		{
			if state[registerInd] == getValue(instruction.RegB, state) {
				state[registerInd] = 1
			} else {
				state[registerInd] = 0
			}
		}
	default:
		{
			log.Fatalf("Unexpected command: '%s'", instruction.Command)
		}
	}
}

func getValue(register interface{}, state []int) int {
	switch register.(type) {
	case byte:
		{
			return state[registerToIndex(register.(byte))]
		}
	case int:
		{
			return register.(int)
		}
	default:
		{
			log.Fatalf("Unexpected register type: '%v'", register)
		}
	}
	return -1
}

func registerToIndex(register byte) int {
	return int(register) - int('w')
}

// Checked the solution discussion and everyone was saying there is a pattern.
// Went back and looked into it more. (All following observations are for my case)
// Noticed that all instructions are the same and only some constants change between blocks.
// Furthermore noticed that each block starts by:
//   - zeroing x => x is irrelevant between blocks
//   - input is always read into w register => w is irrelevant between blocks
//   - y is always zeroed in the block before it gets used => y is irrelevant between blocks
//   - only z persists between blocks and is defined by starting w
//
// z is initially set by the first layer and then just propagates.
// With these we know that only input 'w' matters for each layer, otherwise the whole thing is constant.
// This still means that we have 9^14 options to check, but at least now there are only 2 values to possibly
// cache at each layer, z and input.
func calculateZ(z int, input int, c1, c2, c3 int) int {
	x := (z % 26) + 14 + c2
	z = z / c1
	if x != input {
		z = (z * 26) + input + c3
	}
	return z
}

func PartOne() int {
	return 0
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	instructions := ParseInput(lines)

	aluSections := splitInstructionsIntoAluSections(instructions)

	fmt.Printf("The largest model number accepted by MONAD is: %d\n", PartOneBruteForce(aluSections))
}

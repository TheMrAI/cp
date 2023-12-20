package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Module interface {
	Process(from string, highPulse bool) (bool, bool) // pulse sent, low/high
	Destinations() []string
}

type FlipFlop struct {
	On      bool
	Outputs []string
}

func (f *FlipFlop) Process(from string, highPulse bool) (bool, bool) {
	if highPulse {
		return false, false
	}
	f.On = !f.On
	return true, f.On
}

func (f *FlipFlop) Destinations() []string {
	return f.Outputs
}

type Conjunction struct {
	Memory  map[string]bool
	Outputs []string
}

func (c *Conjunction) Process(from string, highPulse bool) (bool, bool) {
	c.Memory[from] = highPulse
	allHigh := true
	for _, pulse := range c.Memory {
		if !pulse {
			allHigh = false
		}
	}
	if allHigh {
		return true, false
	}
	return true, true
}

func (c *Conjunction) Destinations() []string {
	return c.Outputs
}

type Broadcaster struct {
	Outputs []string
}

func (b *Broadcaster) Process(from string, highPulse bool) (bool, bool) {
	return true, highPulse
}

func (b *Broadcaster) Destinations() []string {
	return b.Outputs
}

type Button struct {
	Outputs []string
}

func (b *Button) Destinations() []string {
	return b.Outputs
}

func (b *Button) Process(from string, highPulse bool) (bool, bool) {
	return true, false
}

func ParseInput(lines []string) map[string]Module {
	moduleSystem := map[string]Module{}

	moduleTypeRegexp := regexp.MustCompile("([\\%\\&]{1})(\\w+)")
	destinationListRegexp := regexp.MustCompile("\\w+")

	// // these need to be maintained separately, because we only know their input connections
	// // after processing the whole system
	// conjunctionModules := map[string]ConjunctionData{}

	for _, line := range lines {
		parts := strings.Split(line, "->")
		if len(parts) != 2 {
			log.Fatalf("Unexpected input: '%v'\n", line)
		}
		parts[0] = strings.Trim(parts[0], " ")
		var module Module
		destinations := destinationListRegexp.FindAllString(parts[1], -1)
		if len(destinations) == 0 {
			log.Fatalf("Expected a non empty destination list after '->', in input: '%v'\n", line)
		}
		moduleIdParts := moduleTypeRegexp.FindStringSubmatch(parts[0])
		moduleName := "broadcaster"
		if len(moduleIdParts) == 0 {
			if parts[0] == "broadcaster" {
				module = &Broadcaster{destinations}
			} else {
				log.Fatalf("Unexpected module type section received: '%v'\n", parts[0])
			}
		} else {
			moduleName = moduleIdParts[2]
			switch moduleIdParts[1] {
			case "%":
				module = &FlipFlop{false, destinations}
			case "&":
				module = &Conjunction{map[string]bool{}, destinations}
			default:
				log.Fatalf("Unsupported module identifier received: '%v'\n", moduleIdParts[1])
			}
		}
		moduleSystem[moduleName] = module
	}
	// add button
	moduleSystem["button"] = &Button{[]string{"broadcaster"}}
	// init all Conjunction module memories now that we have the full system
	for k, v := range moduleSystem {
		for _, destination := range v.Destinations() {
			target := moduleSystem[destination]
			conjunction, ok := target.(*Conjunction)
			if ok {
				conjunction.Memory[k] = false
			}
		}
	}

	return moduleSystem
}

func PartOne(moduleSystem map[string]Module, simulationRounds int) int {
	lowSignalCount := 0
	highSignalCount := 0
	for i := 0; i < simulationRounds; i++ {
		low, high := Simulate(moduleSystem)
		lowSignalCount += low
		highSignalCount += high
	}
	return lowSignalCount * highSignalCount
}

type CycleData struct {
	CycleStartI int
	CycleLength int
}

// Tried brute forcing it to a 100 million button presses, but nothing happened.
// Wasn't even sure I have understood the problem, the wording is confusing.
// Had to check the solutions for some hints, and people were talking about Conjunction
// module chains, because these act like memory.
// Looked into it and my input has 9 of them in this exact configuration:
// &th -> &rz ->
// &pn -> &lf ->
//
//	-> &lb -> rx
//
// &vd -> &br ->
// &fz -> &fk ->
//
//	^ and     ^ have no other inputs!
//
// So based on the conjunction rule that it only outputs low if all of its inputs are high
// we will need an input chain like so:
// L    L     H    L
// This means if I figure out at what cycle lengths do th, pn, vd and fz produce a L output I could
// be able to tell when rx should get its input. The answer should be lcm(cycle th, cycle pn, cycle vd, cycle fz).
// Now since there are no other conjunctions in the system and they otherwise don't depend on one another, between
// button presses their inputs should be repeated. So the said cycle should be mean button press.

// 26_276_404_587_447 - too low
// 250_924_073_918_341 <= the right answer!!!
func PartTwo(moduleSystem map[string]Module) int {
	cycleWatch := map[string]CycleData{
		"th": {0, 0},
		"pn": {0, 0},
		"vd": {0, 0},
		"fz": {0, 0},
	}
	for i := 1; ; i++ {
		queue := []Signal{}
		// press button
		processSignal(moduleSystem, &queue, "elf_hand", "button", true)
		for len(queue) != 0 {
			signalToProcess := queue[0]
			queue = queue[1:]
			low, _ := processSignal(moduleSystem, &queue, signalToProcess.From, signalToProcess.To, signalToProcess.Value)
			data, watching := cycleWatch[signalToProcess.To]
			// && value == -1
			if low >= 1 && watching {
				cycleLength := i - data.CycleStartI
				// fmt.Printf("Found: %v at %v\n", signalToProcess.To, cycleLength)
				cycleWatch[signalToProcess.To] = CycleData{i, cycleLength}
			}
		}
		// This is a quick hack for detecting the cycle. If we were to print out the data we would see, that after
		// the first detected numbers, the cycle lengths will stabilize. That is what we want. Instead of actually
		// detecting a cycle we just wait a few rounds.
		if i > 100_000 && cycleWatch["th"].CycleLength != 0 && cycleWatch["pn"].CycleLength != 0 && cycleWatch["vd"].CycleLength != 0 && cycleWatch["fz"].CycleLength != 0 {
			break
		}
	}

	return EuclidLCM(EuclidLCM(cycleWatch["th"].CycleLength, cycleWatch["pn"].CycleLength), EuclidLCM(cycleWatch["vd"].CycleLength, cycleWatch["fz"].CycleLength))
}

type Signal struct {
	From  string
	To    string
	Value bool
}

func Simulate(moduleSystem map[string]Module) (int, int) { // low, high, rx hit with low
	queue := []Signal{}
	// press button
	lowSignalCount, highSignalCount := processSignal(moduleSystem, &queue, "elf_hand", "button", true)
	for len(queue) != 0 {
		signalToProcess := queue[0]
		queue = queue[1:]
		low, high := processSignal(moduleSystem, &queue, signalToProcess.From, signalToProcess.To, signalToProcess.Value)
		lowSignalCount += low
		highSignalCount += high
	}

	return lowSignalCount, highSignalCount
}

func processSignal(moduleSystem map[string]Module, queue *[]Signal, from string, to string, signal bool) (int, int) {
	module, valid := moduleSystem[to]
	if !valid {
		return 0, 0
	}
	pulseSent, pulseType := module.Process(from, signal)
	if !pulseSent {
		return 0, 0
	}
	count := len(moduleSystem[to].Destinations())
	for _, destination := range moduleSystem[to].Destinations() {
		*queue = append(*queue, Signal{to, destination, pulseType})
	}
	if pulseType {
		return 0, count
	}
	return count, 0
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	moduleSystem := ParseInput(lines)
	// fmt.Println(moduleSystem)
	// for k, v := range moduleSystem {
	// 	fmt.Printf("%v, %v\n", k, v)
	// }

	fmt.Println("Part one")
	fmt.Printf("The multiple of low and high signal counts after 1000 simulation rounds: %v\n", PartOne(moduleSystem, 1000))
	fmt.Println("Part two")
	fmt.Printf("The fewest number of button presses required to deliver a single low pulse to the module named rx: %v\n", PartTwo(moduleSystem))
}

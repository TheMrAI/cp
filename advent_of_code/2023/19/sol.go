package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Comparator string

const (
	Less Comparator = "<"
	More            = ">"
	Pass            = "PASS"
)

type Check struct {
	Field      string
	CompMethod Comparator
	Value      int
	NextCheck  string
}

type Part struct {
	X int
	M int
	A int
	S int
}

func (part Part) GetValueForField(fieldName string) int {
	switch fieldName {
	case "x":
		{
			return part.X
		}
	case "m":
		{
			return part.M
		}
	case "a":
		{
			return part.A
		}
	case "s":
		{
			return part.S
		}
	default:
		{
			log.Fatalf("Unexpected fieldName received: '%v'\n", fieldName)
		}
	}
	return -1
}

func ParseInput(lines []string) (map[string][]Check, []Part) {
	workflows := map[string][]Check{}

	workflowLine := regexp.MustCompile("(\\w+){(.*)}")
	workflowCheck := regexp.MustCompile("(\\w+)(<|>)(-?\\d+):(\\w+)")
	partRegexp := regexp.MustCompile("x=(-?\\d+),m=(-?\\d+),a=(-?\\d+),s=(-?\\d+)")

	i := 0
	for ; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			break
		}
		workflowParts := workflowLine.FindStringSubmatch(line)
		name := workflowParts[1]
		checksText := strings.Split(workflowParts[2], ",")
		checks := []Check{}
		for _, checkText := range checksText {
			checkParts := workflowCheck.FindStringSubmatch(checkText)
			if len(checkParts) == 0 {
				checks = append(checks, Check{"", Pass, 0, checkText})
				continue
			}
			limit, err := strconv.Atoi(checkParts[3])
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'", checkParts[3])
			}
			// well this string to Comparator abuse is sad, but here it is okay
			checks = append(checks, Check{checkParts[1], Comparator(checkParts[2]), limit, checkParts[4]})
		}

		workflows[name] = checks
	}

	i++ // skip the empty line
	parts := []Part{}
	for ; i < len(lines); i++ {
		line := lines[i]
		valueTexts := partRegexp.FindStringSubmatch(line)
		if len(valueTexts) == 0 {
			log.Fatalf("Failed to parse input: '%v'\n", line)
		}
		values := []int{}
		for k := 1; k < len(valueTexts); k++ {
			valueText := valueTexts[k]
			value, err := strconv.Atoi(valueText)
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'\n", valueText)
			}
			values = append(values, value)
		}
		parts = append(parts, Part{values[0], values[1], values[2], values[3]})
	}

	return workflows, parts
}

func PartOne(workflows map[string][]Check, parts []Part) int {
	accepted, _ := SortParts(workflows, parts)
	total := 0
	for _, part := range accepted {
		total += part.A + part.M + part.S + part.X
	}
	return total
}

func SortParts(workflows map[string][]Check, parts []Part) ([]Part, []Part) {
	accepted := []Part{}
	rejected := []Part{}

	for _, part := range parts {
		workflowId := "in"
		checkIndex := 0
		for workflowId != "A" && workflowId != "R" {
			check := workflows[workflowId][checkIndex]
			if check.CompMethod == "PASS" {
				workflowId = check.NextCheck
				checkIndex = 0
				continue
			}
			if check.CompMethod == Less && part.GetValueForField(check.Field) < check.Value {
				workflowId = check.NextCheck
				checkIndex = 0
				continue
			}
			if check.CompMethod == More && part.GetValueForField(check.Field) > check.Value {
				workflowId = check.NextCheck
				checkIndex = 0
				continue
			}
			checkIndex++
		}
		if workflowId == "A" {
			accepted = append(accepted, part)
			continue
		}
		rejected = append(rejected, part)
	}

	return accepted, rejected
}

type Range struct {
	Start int
	Last  int // inclusive!
}

func (r Range) Valid() bool {
	if r.Start <= r.Last {
		return true
	}
	return false
}

func (r Range) Distance() int {
	return 1 + r.Last - r.Start
}

func (r *Range) Adjust(compMethod Comparator, value int, shouldPass bool) {
	switch compMethod {
	case Less:
		{
			if shouldPass {
				r.Last = value - 1
				return
			}
			r.Start = value
		}
	case More:
		{
			if shouldPass {
				r.Start = value + 1
				return
			}
			r.Last = value
		}
	case Pass:
		{
			return // if pass the range doesn't change
		}
	default:
		{
			log.Fatalf("Unexpected comparator received: '%v'\n", compMethod)
		}
	}
}

type Constraints struct {
	X Range
	M Range
	A Range
	S Range
}

func (c Constraints) Valid() bool {
	return c.X.Valid() && c.M.Valid() && c.A.Valid() && c.S.Valid()
}

func (c Constraints) GetField(fieldName string) Range {
	switch fieldName {
	case "x":
		{
			return c.X
		}
	case "m":
		{
			return c.M
		}
	case "a":
		{
			return c.A
		}
	case "s":
		{
			return c.S
		}
	default:
		{
			log.Fatalf("Unexpected fieldName received: '%v'\n", fieldName)
		}
	}
	return Range{}
}

func (c *Constraints) SetField(fieldName string, r Range) {
	switch fieldName {
	case "x":
		{
			c.X = r
		}
	case "m":
		{
			c.M = r
		}
	case "a":
		{
			c.A = r
		}
	case "s":
		{
			c.S = r
		}
	default:
		{
			log.Fatalf("Unexpected fieldName received: '%v'\n", fieldName)
		}
	}
}

func PartTwo(workflows map[string][]Check, rangeStart, rangeLast int) int {
	return countDistinctCombinations(workflows, "in", 0, Constraints{Range{rangeStart, rangeLast}, Range{rangeStart, rangeLast}, Range{rangeStart, rangeLast}, Range{rangeStart, rangeLast}})
}

func countDistinctCombinations(workflows map[string][]Check, currentWorkflow string, checkIndex int, currentConstraints Constraints) int {
	if !currentConstraints.Valid() {
		return 0
	}
	if currentWorkflow == "A" {
		return currentConstraints.X.Distance() * currentConstraints.M.Distance() * currentConstraints.A.Distance() * currentConstraints.S.Distance()
	}
	if currentWorkflow == "R" {
		return 0
	}
	check := workflows[currentWorkflow][checkIndex]
	// a final transition state reached
	if len(check.Field) == 0 {
		return countDistinctCombinations(workflows, check.NextCheck, 0, currentConstraints)
	}

	// generate adjusted passing and failing ranges and branch off to count them all
	constraintsAfterPass, constraintsAfterFail := GeneratePassFailConstraints(currentConstraints, check.Field, check.CompMethod, check.Value)
	optionsCount := countDistinctCombinations(workflows, check.NextCheck, 0, constraintsAfterPass) + countDistinctCombinations(workflows, currentWorkflow, checkIndex+1, constraintsAfterFail)

	return optionsCount
}

func GeneratePassFailConstraints(c Constraints, field string, compMethod Comparator, value int) (Constraints, Constraints) {
	passed := Constraints{c.X, c.M, c.A, c.S}
	passingRange := passed.GetField(field)
	passingRange.Adjust(compMethod, value, true)
	passed.SetField(field, passingRange)

	failed := Constraints{c.X, c.M, c.A, c.S}
	failingRange := failed.GetField(field)
	failingRange.Adjust(compMethod, value, false)
	failed.SetField(field, failingRange)
	return passed, failed
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	workflows, parts := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Sum of all accepted rating numbers: %v\n", PartOne(workflows, parts))
	fmt.Println("Part two")
	fmt.Printf("Total number of accepted distinct possible combinations: %v\n", PartTwo(workflows, 1, 4000))
}

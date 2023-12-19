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

type Workflows map[string][]Check

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

func ParseInput(lines []string) (Workflows, []Part) {
	workflows := Workflows{}

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

func PartOne(workflows Workflows, parts []Part) int {
	accepted, _ := SortParts(workflows, parts)
	total := 0
	for _, part := range accepted {
		total += part.A + part.M + part.S + part.X
	}
	return total
}

func SortParts(workflows Workflows, parts []Part) ([]Part, []Part) {
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
}

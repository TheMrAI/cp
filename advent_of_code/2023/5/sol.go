package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

func ParseSeeds(line string) []int {
	numberRegexp := regexp.MustCompile("\\d+")

	seeds := []int{}
	for _, seedText := range numberRegexp.FindAllString(line, -1) {
		number, err := strconv.Atoi(seedText)
		if err != nil {
			log.Fatalf("Failed number conversion for input: %v\n", seedText)
		}
		seeds = append(seeds, number)
	}

	return seeds
}

type Type int

const (
	Undefined Type = iota
	Seed
	Soil
	Fertilizer
	Water
	Light
	Temperature
	Humidity
	Location
)

func ToType(typeStr string) Type {
	switch typeStr {
	case "seed":
		{
			return Seed
		}
	case "soil":
		{
			return Soil
		}
	case "fertilizer":
		{
			return Fertilizer
		}
	case "water":
		{
			return Water
		}
	case "light":
		{
			return Light
		}
	case "temperature":
		{
			return Temperature
		}
	case "humidity":
		{
			return Humidity
		}
	case "location":
		{
			return Location
		}
	default:
		{
			return Undefined
		}
	}
}

type RangeData struct {
	Start int
	End   int
}

type RangeMapping struct {
	Range  RangeData
	Offset int
}

type Mapping struct {
	To            Type
	RangeMappings []RangeMapping
}

func (mapping *Mapping) Remap(value int) (Type, int) {
	remapped := value
	for _, rangeMapping := range mapping.RangeMappings {
		if !(rangeMapping.Range.Start <= remapped && remapped < rangeMapping.Range.End) {
			continue
		}
		remapped += rangeMapping.Offset
		break
	}
	return mapping.To, remapped
}

func ParseMap(lines []string) (Type, Mapping) {
	fromToRegexp := regexp.MustCompile("(\\w+)-to-(\\w+) map:")
	numberRegexp := regexp.MustCompile("\\d+")

	mapping := Mapping{}
	fromToMatches := fromToRegexp.FindStringSubmatch(lines[0])
	if len(fromToMatches) != 3 {
		log.Fatalf("Failed to parse: %v\n", lines[0])
	}
	mappingFrom := ToType(fromToMatches[1])
	mapping.To = ToType(fromToMatches[2])
	if mappingFrom == Undefined || mapping.To == Undefined {
		log.Fatalf("Failed to parse mapping types from: %v\n", fromToMatches[0])
	}

	for i := 1; i < len(lines); i++ {
		numbersText := numberRegexp.FindAllString(lines[i], -1)
		if len(numbersText) != 3 {
			log.Fatalf("Invalid input received for range data: %v\n", lines[i])
		}
		numbers := []int{}
		for _, numberText := range numbersText {
			number, err := strconv.Atoi(numberText)
			if err != nil {
				log.Fatalf("Failed number conversion on input: %v\n", numberText)
			}
			numbers = append(numbers, number)
		}
		rangeMapping := RangeMapping{RangeData{numbers[1], numbers[1] + numbers[2]}, numbers[0] - numbers[1]}
		mapping.RangeMappings = append(mapping.RangeMappings, rangeMapping)
	}

	return mappingFrom, mapping
}

func ParseMappings(lines []string) map[Type]Mapping {
	mappings := map[Type]Mapping{}

	mappingText := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			if len(mappingText) != 0 {
				from, mapping := ParseMap(mappingText)
				mappings[from] = mapping
			}
			mappingText = nil
			continue
		}
		mappingText = append(mappingText, line)
	}
	if len(mappingText) != 0 {
		from, mapping := ParseMap(mappingText)
		mappings[from] = mapping
	}

	return mappings
}

func PartOne(seeds []int, mappings map[Type]Mapping) int {
	lowestLocationNumber := math.MaxInt
	for _, seed := range seeds {
		currentType := Seed
		value := seed
		for {
			mapping, ok := mappings[currentType]
			if !ok {
				break
			}
			currentType, value = mapping.Remap(value)
		}
		if value < lowestLocationNumber {
			lowestLocationNumber = value
		}
	}
	return lowestLocationNumber
}

func PartTwo(seedRanges []RangeData, mappings map[Type]Mapping) int {
	lowestLocationNumber := math.MaxInt
	for _, rangeData := range seedRanges {
		for startValue := rangeData.Start; startValue < rangeData.End; startValue++ {
			currentType := Seed
			value := startValue
			for {
				mapping, ok := mappings[currentType]
				if !ok {
					break
				}
				currentType, value = mapping.Remap(value)
			}
			if value < lowestLocationNumber {
				lowestLocationNumber = value
			}
		}
	}
	return lowestLocationNumber
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	seeds := ParseSeeds(lines[0])
	mappings := ParseMappings(lines[1:])

	fmt.Println("Part one")
	fmt.Printf("Lowest location number: %v\n", PartOne(seeds, mappings))

	seedRanges := []RangeData{}
	for i := 0; i < len(seeds); i += 2 {
		seedRanges = append(seedRanges, RangeData{seeds[i], seeds[i] + seeds[i+1]})
	}
	fmt.Println("Part two")
	fmt.Printf("Lowest location number: %v\n", PartTwo(seedRanges, mappings))
}

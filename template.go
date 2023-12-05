package main

import (
	"bufio"
	"os"
	"regexp"
)

func ParseInput(lines []string) []int {
	data := []int{}

	numberRegexp := regexp.MustCompile("\\d+")

	for _, line := range lines {

	}

	return data
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)
}

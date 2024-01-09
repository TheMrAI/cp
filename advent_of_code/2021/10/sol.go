package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func PartOne(lines []string) int {
	syntaxErrorScore := 0

	for i := range lines {
		_, check := Validate(lines[i])
		if check == -1 || check == len(lines[i]) {
			continue
		}
		score := Score(lines[i][check])
		if score == -1 {
			log.Fatalf("Unexpected score found: '%v'\n", score)
		}
		syntaxErrorScore += score
	}

	return syntaxErrorScore
}

func Score(c byte) int {
	switch c {
	case ')':
		{
			return 3
		}
	case ']':
		{
			return 57
		}
	case '}':
		{
			return 1197
		}
	case '>':
		{
			return 25137
		}
	}
	return -1
}

func PartTwo(lines []string) int {
	completionStringScores := []int{}
	for i := range lines {
		remainder, check := Validate(lines[i])
		if check == len(lines[i]) {
			completionStringScores = append(completionStringScores, ScoreCompletionString(GenerateCompletionString(remainder)))
		}
	}

	sort.Slice(completionStringScores, func(i, j int) bool {
		return completionStringScores[i] < completionStringScores[j]
	})
	fmt.Println(completionStringScores)
	return completionStringScores[len(completionStringScores)/2]
}

func GenerateCompletionString(stack []byte) []byte {
	completionString := []byte{}

	for len(stack) != 0 {
		toComplete := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		closing := ClosingBrace(toComplete)
		if closing == 'z' {
			log.Fatalf("Unexpected opening brace type detected: '%v'\n", toComplete)
		}
		completionString = append(completionString, closing)
	}

	return completionString
}

func ClosingBrace(c byte) byte {
	switch c {
	case '(':
		{
			return ')'
		}
	case '[':
		{
			return ']'
		}
	case '{':
		{
			return '}'
		}
	case '<':
		{
			return '>'
		}
	}
	return 'z'
}

func ScoreCompletionString(completionString []byte) int {
	score := 0
	for i := range completionString {
		score *= 5
		val := CompletionScore(completionString[i])
		if val == -1 {
			log.Fatalf("Invalid input detected in completion string: '%v'\n", completionString)
		}
		score += val
	}
	return score
}

func CompletionScore(c byte) int {
	switch c {
	case ')':
		{
			return 1
		}
	case ']':
		{
			return 2
		}
	case '}':
		{
			return 3
		}
	case '>':
		{
			return 4
		}
	}
	return -1
}

// Result of -1 means valid, while len(line), means incomplete
func Validate(line string) ([]byte, int) {
	stack := []byte{}
	for i := range line {
		if line[i] == '(' || line[i] == '[' || line[i] == '{' || line[i] == '<' {
			stack = append(stack, line[i])
			continue
		}
		if line[i] == ')' || line[i] == ']' || line[i] == '}' || line[i] == '>' {
			if len(stack) == 0 {
				return stack, i
			}
			candidate := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !BraceMatches(candidate, line[i]) {
				return stack, i
			}
			continue
		}
		log.Fatalf("Invalid character received: '%c'\n", line[i])
	}

	if len(stack) == 0 {
		return stack, -1
	}
	return stack, len(line)
}

func BraceMatches(opening byte, closing byte) bool {
	switch opening {
	case '(':
		{
			return closing == ')'
		}
	case '[':
		{
			return closing == ']'
		}
	case '{':
		{
			return closing == '}'
		}
	case '<':
		{
			return closing == '>'
		}
	}
	return false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Finally one without input parsing!!! Eeehhaaa
	fmt.Println("Part one")
	fmt.Printf("The total of the first syntax error scores is: %v\n", PartOne(lines))
	fmt.Println("Part two")
	fmt.Printf("The middle score of sorted completions string scores is: %v\n", PartTwo(lines))
}

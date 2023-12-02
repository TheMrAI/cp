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

type Selection map[string]int

type Game struct {
	Id         int
	Selections []Selection
}

func ParseGames(lines []string) []Game {
	gameIdRegexp := regexp.MustCompile("Game (\\d+)")
	gameRegexp := regexp.MustCompile("(\\d+) (red|green|blue)")

	games := []Game{}
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid game input received: %v\n", line)
		}
		gameIdData := gameIdRegexp.FindStringSubmatch(parts[0])
		if len(gameIdData) != 2 {
			log.Fatalf("Invalid game id section: %v\n", gameIdData)
		}
		gameId, err := strconv.Atoi(gameIdData[1])
		if err != nil {
			log.Fatal(err)
		}

		selectionTexts := strings.Split(parts[1], ";")
		selections := []Selection{}
		for _, selectionText := range selectionTexts {
			allCubes := gameRegexp.FindAllStringSubmatch(selectionText, -1)
			selection := make(Selection)
			for _, cube := range allCubes {
				count, err := strconv.Atoi(cube[1])
				if err != nil {
					log.Fatal(err)
				}
				selection[cube[2]] += count
			}
			selections = append(selections, selection)
		}

		games = append(games, Game{gameId, selections})
	}
	return games
}

func PartOne(games []Game, constraints Selection) int {
	idSum := 0

	for _, game := range games {
		allValid := true

	check:
		for _, selection := range game.Selections {
			for color, required := range selection {
				total, ok := constraints[color]
				if !ok {
					continue
				}
				if required > total {
					allValid = false
					break check
				}
			}
		}

		if allValid {
			idSum += game.Id
		}
	}

	return idSum
}

func PartTwo(games []Game) int {
	minimumPowerSum := 0

	for _, game := range games {
		minRequired := Selection{}
		for _, selection := range game.Selections {
			for color, required := range selection {
				have := minRequired[color]
				if required > have {
					minRequired[color] = required
				}
			}
		}
		power := 1
		for _, val := range minRequired {
			power *= val
		}
		minimumPowerSum += power
	}

	return minimumPowerSum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	games := ParseGames(lines)
	constraints := Selection{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	fmt.Println("Part one")
	fmt.Printf("Sum of possible game IDs with constraints: %v => %v\n", constraints, PartOne(games, constraints))
	fmt.Println("Part two")
	fmt.Printf("Sum of minimum powers for all games: %v\n", PartTwo(games))
}

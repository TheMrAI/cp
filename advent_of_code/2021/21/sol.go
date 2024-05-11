package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Player struct {
	Id    int
	Pos   int
	Score int
}

func ParseInput(lines []string) []Player {
	players := []Player{}

	playersRegexp := regexp.MustCompile(`Player (\d+) starting position: (\d+)`)

	for _, line := range lines {
		match := playersRegexp.FindStringSubmatch(line)
		id, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatalf("String to integer conversion failed on input: '%v'", match[1])
		}
		start, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatalf("String to integer conversion failed on input: '%v'", match[2])
		}
		players = append(players, Player{
			id, start, 0,
		})
	}

	return players
}

func SimulateUntilWithDeterministicDice(players []Player, limit int) ([]Player, int) {
	playerIndex := 0
	dice := 0
	diceRollCount := 0

	for {
		total := 0
		for count := 0; count < 3; count++ {
			dice++
			diceRollCount++
			if dice > 100 {
				dice = dice % 100
			}
			total += dice
		}
		players[playerIndex].Pos += total
		players[playerIndex].Pos = ((players[playerIndex].Pos - 1) % 10) + 1
		players[playerIndex].Score += players[playerIndex].Pos
		if players[playerIndex].Score >= limit {
			return players, diceRollCount
		}
		playerIndex = (playerIndex + 1) % len(players)
	}
	return players, diceRollCount
}

func PartOne(players []Player) int {
	results, diceRollCount := SimulateUntilWithDeterministicDice(players, 1_000)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results[1].Score * diceRollCount
}

type DiceCount struct {
	DiceSum int
	Count   int
}

var dieTotals = [...]DiceCount{{3, 1}, {4, 3}, {5, 6}, {6, 7}, {7, 6}, {8, 3}, {9, 1}}

func SimulatePlayerWinsUntilWithDiracDice(playerOne, playerTwo Player, limit int) (int, int) {
	if playerOne.Score >= limit {
		return 1, 0
	} else if playerTwo.Score >= limit {
		return 0, 1
	}
	playerOneWins, playerTwoWins := 0, 0

	for _, dirac := range dieTotals {
		nextPlayer := Player{playerOne.Id, playerOne.Pos, playerOne.Score}
		nextPlayer.Pos += dirac.DiceSum
		nextPlayer.Pos = ((nextPlayer.Pos - 1) % 10) + 1
		nextPlayer.Score += nextPlayer.Pos
		// swap the players, so that we only need to move one of them
		scoreOne, scoreTwo := SimulatePlayerWinsUntilWithDiracDice(playerTwo, nextPlayer, limit)
		playerOneWins += scoreTwo * dirac.Count
		playerTwoWins += scoreOne * dirac.Count
	}
	return playerOneWins, playerTwoWins
}

func PartTwo(playerOne Player, playerTwo Player, limit int) int {
	oneWins, twoWins := SimulatePlayerWinsUntilWithDiracDice(playerOne, playerTwo, limit)
	if oneWins > twoWins {
		return oneWins
	}
	return twoWins
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	players := ParseInput(lines)
	playerOne := Player{players[0].Id, players[0].Pos, 0}
	playerTwo := Player{players[1].Id, players[1].Pos, 0}

	fmt.Printf("The score of the losing player multiplied with dice roll count using deterministic die: %v\n", PartOne(players))
	fmt.Printf("The number of universes the dominant player wins in: %v\n", PartTwo(playerOne, playerTwo, 21))
}

// 13843825473600
// 444356092776315

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Shape int

type DesiredOutcome string

type Outcome string

type GameMetrics struct {
	playerScore   int
	opponentScore int
	rounds        map[int]Outcome // map of round number to outcome
}

const (
	ShapeRock     Shape = 1
	ShapePaper    Shape = 2
	ShapeScissors Shape = 3

	DesiredOutcomeLose DesiredOutcome = "lose"
	DesiredOutcomeDraw DesiredOutcome = "draw"
	DesiredOutcomeWin  DesiredOutcome = "win"

	OutcomePlayer   Outcome = "player"
	OutcomeOpponent Outcome = "opponent"
	OutcomeDraw     Outcome = "draw"

	OutcomeScoreLost = 0
	OutcomeScoreDraw = 3
	OutcomeScoreWin  = 6
)

func getStrategyGuide() map[string]Shape {
	return map[string]Shape{
		// Opponent map
		"A": ShapeRock,
		"B": ShapePaper,
		"C": ShapeScissors,

		// Player map
		"X": ShapeRock,
		"Y": ShapePaper,
		"Z": ShapeScissors,
	}
}

func getDesiredOutcome(input string) DesiredOutcome {
	desiredOutcomes := map[string]DesiredOutcome{
		"X": DesiredOutcomeLose,
		"Y": DesiredOutcomeDraw,
		"Z": DesiredOutcomeWin,
	}
	return desiredOutcomes[input]
}

func getPlayerShape(opponentShape Shape, input string) Shape {
	desiredOutcome := getDesiredOutcome(input)

	switch desiredOutcome {
	case DesiredOutcomeDraw:
		return opponentShape
	case DesiredOutcomeWin:
		switch opponentShape {
		case ShapeRock:
			return ShapePaper
		case ShapePaper:
			return ShapeScissors
		case ShapeScissors:
			return ShapeRock
		}
	case DesiredOutcomeLose:
		switch opponentShape {
		case ShapeRock:
			return ShapeScissors
		case ShapePaper:
			return ShapeRock
		case ShapeScissors:
			return ShapePaper
		}
	}

	return 0
}

func getShapeName(shape Shape) string {
	switch shape {
	case ShapeRock:
		return "Rock"
	case ShapePaper:
		return "Paper"
	case ShapeScissors:
		return "Scissors"
	default:
		return ""
	}
}

func getRoundOutcome(opponentShape Shape, playerShape Shape) Outcome {
	switch opponentShape {
	case playerShape:
		return OutcomeDraw
	case ShapeRock:
		if playerShape == ShapeScissors {
			return OutcomeOpponent
		}

		if playerShape == ShapePaper {
			return OutcomePlayer
		}
	case ShapePaper:
		if playerShape == ShapeRock {
			return OutcomeOpponent
		}

		if playerShape == ShapeScissors {
			return OutcomePlayer
		}
	case ShapeScissors:
		if playerShape == ShapePaper {
			return OutcomeOpponent
		}

		if playerShape == ShapeRock {
			return OutcomePlayer
		}
	default:
		log.Panicf("Invalid shape: %v", opponentShape)
	}

	return ""
}

func readInput(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {
	logRoundResults := false
	gameMetrics := &GameMetrics{
		playerScore:   0,
		opponentScore: 0,
		rounds:        make(map[int]Outcome),
	}
	strategyGuide := getStrategyGuide()

	input, err := readInput("input.txt")
	if err != nil {
		log.Panicf("Error reading input: %v", err)
	}

	gameRounds := strings.Split(input, "\n")

	for idx, round := range gameRounds {
		opponentScore := 0
		playerScore := 0

		shapes := strings.Split(round, " ")

		opponentShape := strategyGuide[shapes[0]]
		playerShape := getPlayerShape(opponentShape, shapes[1])

		roundOutcome := getRoundOutcome(opponentShape, playerShape)
		switch roundOutcome {
		case OutcomeOpponent:
			opponentScore += int(opponentShape) + OutcomeScoreWin
			playerScore += int(playerShape) + OutcomeScoreLost
		case OutcomePlayer:
			playerScore += int(playerShape) + OutcomeScoreWin
			opponentScore += int(opponentShape) + OutcomeScoreLost
		case OutcomeDraw:
			playerScore += int(playerShape) + OutcomeScoreDraw
			opponentScore += int(opponentShape) + OutcomeScoreDraw
		default:
			log.Panicf("Invalid round outcome: %v", roundOutcome)
		}

		gameMetrics.opponentScore += opponentScore
		gameMetrics.playerScore += playerScore
		gameMetrics.rounds[idx] = roundOutcome

		if logRoundResults {
			fmt.Printf("Round %v outcome -> %v\n", idx+1, roundOutcome)
			fmt.Printf("Opponent Shape -> %v, Player shape -> %v\n", getShapeName(opponentShape), getShapeName(playerShape))
			fmt.Printf("Opponent score -> %v, Player score -> %v\n", opponentScore, playerScore)
			fmt.Printf("======\n")
		}
	}

	fmt.Printf("Opponent score: %v. Player score: %v\n", gameMetrics.opponentScore, gameMetrics.playerScore)
}

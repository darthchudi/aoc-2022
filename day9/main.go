package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Move struct {
	direction string
	count     int
}

type Position struct {
	row    int
	column int
}

func computeKnotPosition(headPosition Position, tailPosition Position) *Position {
	newTailPosition := &Position{
		column: tailPosition.column,
		row:    tailPosition.row,
	}

	isOnSameRow := headPosition.row == tailPosition.row
	isOnSameColumn := headPosition.column == tailPosition.column

	isOverlapping := isOnSameRow && isOnSameColumn
	isHorizontallyAdjacent := isOnSameRow && math.Abs(float64(headPosition.column-tailPosition.column)) == 1
	isVerticallyAdjacent := isOnSameColumn && math.Abs(float64(headPosition.row-tailPosition.row)) == 1
	isDiagonallyAdjacent := math.Abs(float64(headPosition.column-tailPosition.column)) == 1 && math.Abs(float64(headPosition.row-tailPosition.row)) == 1

	isTouching := isOverlapping || isHorizontallyAdjacent || isVerticallyAdjacent || isDiagonallyAdjacent
	if isTouching {
		return newTailPosition
	}

	switch {
	case isOnSameRow:
		isTwoStepsBehind := (headPosition.column - tailPosition.column) == 2
		if isTwoStepsBehind {
			newTailPosition.column = tailPosition.column + 1
			return newTailPosition
		}

		isTwoStepsAhead := (tailPosition.column - headPosition.column) == 2
		if isTwoStepsAhead {
			newTailPosition.column = tailPosition.column - 1
			return newTailPosition
		}
	case isOnSameColumn:
		isTwoStepsAbove := (headPosition.row - tailPosition.row) == 2
		if isTwoStepsAbove {
			newTailPosition.row = tailPosition.row + 1
			return newTailPosition
		}

		isTwoStepsBelow := (tailPosition.row - headPosition.row) == 2
		if isTwoStepsBelow {
			newTailPosition.row = tailPosition.row - 1
			return newTailPosition
		}
	}

	// The tail isn't in the same row or column as the head, so we have to move diagonally.
	isAbove := headPosition.row > tailPosition.row
	isBelow := headPosition.row < tailPosition.row

	isToTheRight := tailPosition.column > headPosition.column
	isToTheLeft := tailPosition.column < headPosition.column

	if isAbove {
		// Move diagonally downwards
		newTailPosition.row = tailPosition.row + 1
	} else if isBelow {
		// Move diagonally upwards
		newTailPosition.row -= 1

	}

	if isToTheLeft {
		// Move diagonally to the right
		newTailPosition.column += 1
	} else if isToTheRight {
		// Move diagonally to the left
		newTailPosition.column -= 1
	}

	return newTailPosition
}

func computeMultipleKnotPositions(headPosition Position, knots []*Position) []*Position {
	result := make([]*Position, len(knots))

	for idx, knot := range knots {
		if idx == 0 {
			// Use the actual head knot to compute the new position
			// of the first knot.
			newPosition := computeKnotPosition(headPosition, *knot)
			result[idx] = newPosition
			continue
		}

		// Use the previous knot to compute the new position of the current knot.
		previousKnot := result[idx-1]
		newPosition := computeKnotPosition(*previousKnot, *knot)
		result[idx] = newPosition
	}

	return result
}

func markVisitedPosition(visitedPositions map[string]bool, row, column int) {
	key := fmt.Sprintf("row-%v-column-%v", row, column)
	visitedPositions[key] = true
}

func getGridOutput(grid [][]int, headPosition *Position, tailPosition *Position, tailVisitedPositions map[string]bool, knots []*Position) string {
	seenKnotsOutputMap := map[string]bool{}
	includeVisitedSpotsInOutput := false
	includeKnotsInOutput := false

	output := ""

	for rowIdx, row := range grid {
		for columnIdx, _ := range row {
			if rowIdx == headPosition.row && columnIdx == headPosition.column {
				output += "H"
				continue
			}

			if rowIdx == tailPosition.row && columnIdx == tailPosition.column {
				output += "T"
				continue
			}

			containsKnot := false
			if includeKnotsInOutput {
				for idx, knot := range knots {
					// Check if the knot is in the current position
					if rowIdx == knot.row && columnIdx == knot.column {
						key := fmt.Sprintf("%d-%d", rowIdx, columnIdx)

						// Check if we've already seen and marked a knot in this position
						if _, ok := seenKnotsOutputMap[key]; ok {
							continue
						}

						output += strconv.Itoa(idx + 1)
						containsKnot = true
						seenKnotsOutputMap[key] = true
						continue
					}
				}
			}

			if containsKnot {
				continue
			}

			if !includeVisitedSpotsInOutput {
				output += "."
				continue
			}

			key := fmt.Sprintf("row-%v-column-%v", rowIdx, columnIdx)
			_, ok := tailVisitedPositions[key]
			if ok {
				output += "#"
			} else {
				output += "."
			}
		}

		output += "\n"
	}

	return output
}

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	var moves []*Move

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		count, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("failed to convert to int: %s", err)
		}

		move := &Move{
			direction: parts[0],
			count:     count,
		}
		moves = append(moves, move)
	}

	maxVerticalMove := 0
	maxHorizontalMove := 0
	for _, move := range moves {
		switch move.direction {
		case "U", "D":
			if move.count > maxVerticalMove {
				maxVerticalMove = move.count
			}
		case "L", "R":
			if move.count > maxHorizontalMove {
				maxHorizontalMove = move.count
			}
		}
	}

	// Create the grid based on the max moves
	rowCount := maxVerticalMove + 1
	columnCount := maxHorizontalMove + 1

	grid := make([][]int, rowCount)
	for idx, _ := range grid {
		grid[idx] = make([]int, columnCount)
	}

	// Set the starting position for the head and tail knots
	headPosition := &Position{
		row:    len(grid) - 1,
		column: 0,
	}
	tailPosition := &Position{
		row:    len(grid) - 1,
		column: 0,
	}
	tailVisitedPositions := map[string]bool{}

	// Create 9 knots
	knots := make([]*Position, 9)
	for idx, _ := range knots {
		knots[idx] = &Position{
			row:    len(grid) - 1,
			column: 0,
		}
	}
	lastKnotVisitedPositions := map[string]bool{}

	for _, move := range moves {
		switch move.direction {
		case "R":
			// move right
			// Each step in a move is handled separately
			for i := 0; i < move.count; i++ {
				headPosition.column += 1
				tailPosition = computeKnotPosition(*headPosition, *tailPosition)
				markVisitedPosition(tailVisitedPositions, tailPosition.row, tailPosition.column)

				knots = computeMultipleKnotPositions(*headPosition, knots)
				lastKnot := knots[len(knots)-1]
				markVisitedPosition(lastKnotVisitedPositions, lastKnot.row, lastKnot.column)
			}
		case "L":
			// move left
			for i := 0; i < move.count; i++ {
				headPosition.column -= 1
				tailPosition = computeKnotPosition(*headPosition, *tailPosition)
				markVisitedPosition(tailVisitedPositions, tailPosition.row, tailPosition.column)

				knots = computeMultipleKnotPositions(*headPosition, knots)
				lastKnot := knots[len(knots)-1]
				markVisitedPosition(lastKnotVisitedPositions, lastKnot.row, lastKnot.column)
			}
		case "U":
			// move up
			for i := 0; i < move.count; i++ {
				headPosition.row -= 1
				tailPosition = computeKnotPosition(*headPosition, *tailPosition)
				markVisitedPosition(tailVisitedPositions, tailPosition.row, tailPosition.column)

				knots = computeMultipleKnotPositions(*headPosition, knots)
				lastKnot := knots[len(knots)-1]
				markVisitedPosition(lastKnotVisitedPositions, lastKnot.row, lastKnot.column)
			}
		case "D":
			// move down
			for i := 0; i < move.count; i++ {
				headPosition.row += 1
				tailPosition = computeKnotPosition(*headPosition, *tailPosition)
				markVisitedPosition(tailVisitedPositions, tailPosition.row, tailPosition.column)

				knots = computeMultipleKnotPositions(*headPosition, knots)
				lastKnot := knots[len(knots)-1]
				markVisitedPosition(lastKnotVisitedPositions, lastKnot.row, lastKnot.column)
			}
		}
	}

	shouldPrintGridOutput := true
	fmt.Println("Visited positions by tail knot: ", len(tailVisitedPositions))
	fmt.Println("Visited positions by last knot: ", len(lastKnotVisitedPositions))

	if !shouldPrintGridOutput {
		return
	}

	output := getGridOutput(grid, headPosition, tailPosition, tailVisitedPositions, knots)
	fmt.Println(output)
}

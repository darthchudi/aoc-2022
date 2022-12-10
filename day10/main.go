package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	row    int
	column int
}

type Sprite struct {
	start int
	end   int
}

type Output struct {
	value    string
	position *Position
}

// Move moves the sprite based on the register count
// The register count identifies the center position of
// the 3px wide sprite
func (s *Sprite) Move(registerCount int) {
	s.start = registerCount - 1
	s.end = registerCount + 1
}

func (o *Output) Draw(sprite *Sprite, currentCycle, registerCount int) {
	value := "."

	for i := sprite.start; i <= sprite.end; i++ {
		if i == o.position.column {
			value = "#"
		}
	}

	o.value += value
	o.position.column += 1

	if currentCycle%40 == 0 {
		// Wrap to a new row because the CRT is 40px wide and 6px high
		o.value += "\n"
		o.position.row += 1
		o.position.column = 0

	}
}

func main() {
	currentCycle := 0
	registerCount := 1

	signalStrengths := map[int]int{}

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	output := &Output{
		value: "",
		position: &Position{
			row:    0,
			column: 0,
		},
	}

	sprite := &Sprite{
		start: 0,
		end:   2,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		instruction := parts[0]

		switch {
		case instruction == "addx":
			argument, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatalf("failed parsing argument: %s", err)
			}

			// Two cycles
			for i := 0; i < 2; i++ {
				currentCycle += 1

				output.Draw(sprite, currentCycle, registerCount)

				isRelevantCycle := currentCycle%20 == 0
				if isRelevantCycle {
					signalStrengths[currentCycle] = currentCycle * registerCount
				}
			}

			registerCount += argument
			sprite.Move(registerCount)
		case instruction == "noop":
			// One cycle
			for i := 0; i < 1; i++ {
				currentCycle += 1

				output.Draw(sprite, currentCycle, registerCount)

				isRelevantCycle := currentCycle%20 == 0
				if isRelevantCycle {
					signalStrengths[currentCycle] = currentCycle * registerCount
				}
			}
		default:
			log.Fatalf("unknown instruction: %s", instruction)
		}
	}

	sum := signalStrengths[20] + signalStrengths[60] + signalStrengths[100] + signalStrengths[140] + signalStrengths[180] + signalStrengths[220]
	fmt.Println("Sum:", sum)

	fmt.Println(output.value)
}

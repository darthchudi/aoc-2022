package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Instruction struct {
	count       int
	source      int
	destination int
}

type CrateConfig map[int][]string

func parseInstructions(instructionLines []string) ([]*Instruction, error) {
	var instructions []*Instruction

	re := regexp.MustCompile("move ([0-9]+) from ([0-9]+) to ([0-9]+)")

	for _, line := range instructionLines {
		matches := re.FindStringSubmatch(line)
		matches = matches[1:]

		count, err := strconv.Atoi(matches[0])
		if err != nil {
			return nil, err
		}

		source, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}

		destination, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, err
		}

		instruction := &Instruction{
			count:       count,
			source:      source,
			destination: destination,
		}

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

func parseCrates(crateLines []string, maxColumnsStr string) (CrateConfig, error) {
	crates := map[int][]string{}

	maxColumns, err := strconv.Atoi(maxColumnsStr)
	if err != nil {
		return nil, err
	}

	characterRegex := regexp.MustCompile(`([A-Z]+)`)

	// Map of the index of the crate to the column number
	chracterIndexToCrateNumber := map[int]int{}

	currentCharacterIndex := 1
	for i := 0; i < maxColumns; i++ {
		chracterIndexToCrateNumber[currentCharacterIndex] = i + 1
		currentCharacterIndex += 4 // 4 characters per crate including spaces
	}

	for _, line := range crateLines {
		for idx, character := range line {
			if characterRegex.MatchString(string(character)) {
				crateNumber, ok := chracterIndexToCrateNumber[idx]
				if !ok {
					log.Panicf("crate number not found for index: %v", idx)
				}
				crates[crateNumber] = append(crates[crateNumber], string(character))
			}
		}
	}

	return crates, nil
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Panicf("failed to open file: %v", err)
	}
	defer file.Close()

	isPart1 := false

	var crateLines []string
	var instructionLines []string
	var maxColumns string

	crateColumnLine := regexp.MustCompile(`^ [0-9]+`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "move") {
			instructionLines = append(instructionLines, line)
			continue
		}

		if crateColumnLine.MatchString(line) {
			maxColumns = string(line[len(line)-1])
			continue
		}

		if strings.Contains(line, "[") {
			crateLines = append(crateLines, line)
			continue
		}
	}

	instructions, err := parseInstructions(instructionLines)
	if err != nil {
		log.Panicf("failed to parse instructions: %v", err)
	}

	crates, err := parseCrates(crateLines, maxColumns)
	if err != nil {
		log.Panicf("failed to parse crates: %v", err)
	}

	for _, instruction := range instructions {
		source := crates[instruction.source]
		destination := crates[instruction.destination]

		itemsToMove := source[:instruction.count]

		updatedSource := append([]string{}, source[instruction.count:]...)
		updatedDestination := append([]string{}, itemsToMove...)

		if isPart1 {
			// Solutions for part 1 needs to be reversed
			updatedDestination = reverse(updatedDestination)
		}
		updatedDestination = append(updatedDestination, destination...)

		crates[instruction.source] = updatedSource
		crates[instruction.destination] = updatedDestination
	}

	keys := []int{}
	for key, _ := range crates {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	result := ""
	for _, key := range keys {
		result += crates[key][0]
	}

	fmt.Println("Result: ", result)
}

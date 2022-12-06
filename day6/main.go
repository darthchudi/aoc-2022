package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func hasUniqueCharacters(input []rune) bool {
	seen := map[rune]bool{}

	for _, character := range input {
		if seen[character] {
			return false
		}

		seen[character] = true
	}

	return true
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	dataStream := string(bytes)

	var stack []rune
	resultIndex := 0
	start := 0
	sequenceLength := 14

	for index, character := range dataStream {
		stack = append(stack, character)

		if len(stack) >= sequenceLength {
			mostRecentCharactersInSequence := stack[start : index+1]

			if hasUniqueCharacters(mostRecentCharactersInSequence) {
				resultIndex = index
				fmt.Println("Found unique characters: ", string(mostRecentCharactersInSequence), " at index: ", resultIndex)
				break
			}

			start += 1 // move the start index forward
		}
	}

	fmt.Printf("Processed %v characters\n", resultIndex+1)
}

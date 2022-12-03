package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

type RuckSack struct {
	items         string
	duplicateItem rune
}

func getItemPriority(item rune) int {
	priority := 0

	switch {
	case unicode.IsLower(item):
		priority = (int(item) % 97) + 1 // priority 1 -> 26
	case unicode.IsUpper(item):
		priority = (int(item) % 65) + 27 // priority 27 -> 52
	}

	return priority
}

func getDuplicateItemInCompartments(value string) rune {
	half := len(value) / 2
	firstCompartment, secondCompartment := value[:half], value[half:]

	existingItems := map[string]bool{}
	for _, item := range firstCompartment {
		existingItems[string(item)] = true
	}

	// Check for duplicates in the second compartment
	for _, item := range secondCompartment {
		_, ok := existingItems[string(item)]
		if ok {
			return item
		}
	}

	return 0
}

func findDuplicateItemInElfGroup(elfGroup []string) rune {
	firstElf := elfGroup[0]
	elfGroup = elfGroup[1:]

	// Build a lookup map of items in the first elf's bag
	lookup := map[rune]int{}
	for _, item := range firstElf {
		if _, ok := lookup[item]; ok {
			continue
		}

		lookup[item] = 1
	}

	for _, elf := range elfGroup {
		seenItemsInLookup := map[rune]bool{}

		for _, item := range elf {
			if seenItemsInLookup[item] {
				continue
			}

			if _, ok := lookup[item]; ok {
				lookup[item]++
				seenItemsInLookup[item] = true
			}
		}
	}

	// Find the item that appears in all three elf groups
	for item, count := range lookup {
		if count == 3 {
			return item
		}
	}

	return 0
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Panicf("failed to open file: %v", err)
	}

	currentIndex := 0
	elfGroups := make([][]string, 1)

	r := bufio.NewScanner(file)
	for r.Scan() {
		if len(elfGroups[currentIndex]) == 3 {
			currentIndex++
			elfGroups = append(elfGroups, []string{})
		}

		elfGroups[currentIndex] = append(elfGroups[currentIndex], r.Text())
	}

	sum := 0
	for _, elfGroup := range elfGroups {
		duplicateItem := findDuplicateItemInElfGroup(elfGroup)
		priority := getItemPriority(duplicateItem)

		sum += priority
	}

	fmt.Println(sum)
}

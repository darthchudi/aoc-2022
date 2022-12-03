package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ElfStat struct {
	index    int
	calories int
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

func getElfCalories(input string) ([][]int, error) {
	splitLines := strings.Split(input, "\n")

	currentIndex := 0
	elfCalories := make([][]int, 1)

	for _, line := range splitLines {
		if line == "" {
			// Move to the next elf and grow the array
			currentIndex++
			updatedElfCalories := make([][]int, currentIndex+1)
			copy(updatedElfCalories, elfCalories)
			elfCalories = updatedElfCalories

			continue
		}

		value, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		elfCalories[currentIndex] = append(elfCalories[currentIndex], value)
	}

	return elfCalories, nil
}

func sumIntegers(integers []int) int {
	var sum int
	for _, integer := range integers {
		sum += integer
	}
	return sum
}

func getElfStats(elfCalories [][]int) []*ElfStat {
	var elfStats []*ElfStat

	for idx, elf := range elfCalories {
		elfStats = append(elfStats, &ElfStat{
			index:    idx,
			calories: sumIntegers(elf),
		})
	}

	sort.Slice(elfStats, func(i, j int) bool {
		return elfStats[i].calories > elfStats[j].calories
	})

	return elfStats
}

func findTopElvesWithCalories(elfStats []*ElfStat, n int) ([]*ElfStat, int, error) {
	if n > len(elfStats) {
		return nil, 0, fmt.Errorf("n is greater than the number of elves")
	}

	var total int
	topElves := elfStats[:n]

	for _, elf := range topElves {
		total += elf.calories
	}

	return topElves, total, nil
}

func main() {
	input, err := readInput("input.txt")
	if err != nil {
		log.Panicf("Error reading input: %v", err)
	}

	elfCalories, err := getElfCalories(input)
	if err != nil {
		log.Panicf("Error getting elf calories: %v", err)
	}

	elfStats := getElfStats(elfCalories)

	_, total, err := findTopElvesWithCalories(elfStats, 3)
	if err != nil {
		log.Panicf("Error finding top elves: %v", err)
	}

	fmt.Printf("total calories by top 3 elves: %v\n\n", total)

	elfIndex, mostCalories := elfStats[0].index, elfStats[0].calories
	fmt.Printf("The most calories are %v calories carried by elf #%v\n", mostCalories, elfIndex+1)
}

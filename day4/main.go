package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type SectionAssignment struct {
	start int
	end   int
}

func getSectionAssignment(assignmentStr string) (*SectionAssignment, error) {
	assignmentValues := strings.Split(assignmentStr, "-")

	start, err := strconv.Atoi(assignmentValues[0])
	if err != nil {
		return nil, err
	}

	end, err := strconv.Atoi(assignmentValues[1])
	if err != nil {
		return nil, err
	}

	return &SectionAssignment{
		start: start,
		end:   end,
	}, nil
}

func isFullyContainedPair(firstElfSectionAssignment *SectionAssignment, secondElfSectionAssignment *SectionAssignment) bool {
	// Check if the second elf's section is fully contained in the first elf's section
	if firstElfSectionAssignment.start <= secondElfSectionAssignment.start && firstElfSectionAssignment.end >= secondElfSectionAssignment.end {
		return true
	}

	// Check if the first elf's section is fully contained in the second elf's section
	if secondElfSectionAssignment.start <= firstElfSectionAssignment.start && secondElfSectionAssignment.end >= firstElfSectionAssignment.end {
		return true
	}

	return false
}

func isOverlappingPair(firstElfSectionAssignment *SectionAssignment, secondElfSectionAssignment *SectionAssignment) bool {
	// Check if the first elf's section overlaps with the second elf's section
	if firstElfSectionAssignment.start <= secondElfSectionAssignment.start && firstElfSectionAssignment.end >= secondElfSectionAssignment.start {
		return true
	}

	// Check if the second elf's section overlaps with the first elf's section
	if secondElfSectionAssignment.start <= firstElfSectionAssignment.start && secondElfSectionAssignment.end >= firstElfSectionAssignment.start {
		return true
	}

	return false
}

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		log.Panicf("failed to open file: %v", err)
	}
	defer file.Close()

	fullyContainedPairs := 0
	overlappingPairs := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		elfPair := strings.Split(scanner.Text(), ",")

		firstElfSectionAssignment, err := getSectionAssignment(elfPair[0])
		if err != nil {
			log.Panicf("failed to parse first elf section assignment: %v", err)
		}

		secondElfSectionAssignment, err := getSectionAssignment(elfPair[1])
		if err != nil {
			log.Panicf("failed to parse second elf section assignment: %v", err)
		}

		if isFullyContainedPair(firstElfSectionAssignment, secondElfSectionAssignment) {
			fullyContainedPairs++
		}

		if isOverlappingPair(firstElfSectionAssignment, secondElfSectionAssignment) {
			overlappingPairs++
		}
	}

	log.Printf("fully contained pairs: %d\n", fullyContainedPairs)
	log.Printf("overlapping pairs: %d\n", overlappingPairs)
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Operation struct {
	Operator             string
	Delta                int
	isRecursiveOperation bool
}

type Test struct {
	Delta  int // Delta is the amount to divide by
	OnPass int // monkey to throw to if the test passes
	OnFail int // monkey to throw to if the test fails
}

type Monkey struct {
	index          int
	items          []int
	operation      *Operation
	test           *Test
	inspectedItems int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	isPartOne := false

	var monkeys []*Monkey
	currentMonkey := &Monkey{}
	monkeyMap := map[int]*Monkey{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			// Store the current monkey and move to the next monkey
			// as a line break separates monkeys
			monkeys = append(monkeys, currentMonkey)
			monkeyMap[currentMonkey.index] = currentMonkey
			currentMonkey = &Monkey{}
			continue
		}

		switch {
		case strings.HasPrefix(line, "Monkey"):
			parts := strings.Split(line, " ")
			indexParts := strings.Split(parts[1], ":")

			index, err := strconv.Atoi(indexParts[0])
			if err != nil {
				log.Fatalf("failed parsing monkey index: %s", err)
			}

			currentMonkey.index = index
		case strings.HasPrefix(line, "Starting items:"):
			parts := strings.Split(line, "Starting items: ")
			itemsStringSlice := strings.Split(parts[1], ", ")

			for _, itemString := range itemsStringSlice {
				item, err := strconv.Atoi(itemString)
				if err != nil {
					log.Fatalf("failed parsing monkey items: %s", err)
				}
				currentMonkey.items = append(currentMonkey.items, item)
			}
		case strings.HasPrefix(line, "Operation:"):
			operationRegex := regexp.MustCompile(`Operation: new = old (.+) (.+)`)
			matches := operationRegex.FindStringSubmatch(line)
			matches = matches[1:]

			operator, deltaStr := matches[0], matches[1]
			delta := 0
			isRecursiveOperation := false

			switch deltaStr {
			case "old":
				isRecursiveOperation = true
			default:
				delta, err = strconv.Atoi(deltaStr)
				if err != nil {
					log.Fatalf("failed parsing operation delta: %s", err)
				}
			}

			currentMonkey.operation = &Operation{
				Operator:             operator,
				Delta:                delta,
				isRecursiveOperation: isRecursiveOperation,
			}
		case strings.HasPrefix(line, "Test:"):
			testRegex := regexp.MustCompile(`Test: divisible by (.+)`)
			matches := testRegex.FindStringSubmatch(line)
			matches = matches[1:]

			delta, err := strconv.Atoi(matches[0])
			if err != nil {
				log.Fatalf("failed parsing test delta: %s", err)
			}

			currentMonkey.test = &Test{
				Delta:  delta,
				OnPass: 0,
				OnFail: 0,
			}
		case strings.HasPrefix(line, "If true:"):
			onTestPassRegex := regexp.MustCompile(`If true: throw to monkey (.+)`)
			matches := onTestPassRegex.FindStringSubmatch(line)
			matches = matches[1:]

			value, err := strconv.Atoi(matches[0])
			if err != nil {
				log.Fatalf("failed parsing on test pass value: %s", err)
			}

			currentMonkey.test.OnPass = value
		case strings.HasPrefix(line, "If false:"):
			onTestFailRegex := regexp.MustCompile(`If false: throw to monkey (.+)`)
			matches := onTestFailRegex.FindStringSubmatch(line)
			matches = matches[1:]

			value, err := strconv.Atoi(matches[0])
			if err != nil {
				log.Fatalf("failed parsing on test fail value: %s", err)
			}

			currentMonkey.test.OnFail = value
		default:
			log.Fatalf("Invalid line in input: %v", line)
		}
	}
	// Store the last monkey
	monkeys = append(monkeys, currentMonkey)
	monkeyMap[currentMonkey.index] = currentMonkey

	dividend := 1
	for _, monkey := range monkeys {
		dividend = dividend * monkey.test.Delta
	}

	rounds := 10000
	for i := 0; i < rounds; i++ {
		for _, monkey := range monkeys {
			// On each round, inspect all items belonging to
			// a monkey and carry out an operation on the items
			// given their worry level
			for _, item := range monkey.items {
				worryLevel := item

				delta := monkey.operation.Delta
				if monkey.operation.isRecursiveOperation {
					delta = worryLevel
				}

				switch monkey.operation.Operator {
				case "+":
					worryLevel = worryLevel + delta
				case "*":
					worryLevel = worryLevel * delta
				default:
					log.Fatalf("Invalid monkey operator: %v", monkey.operation.Operator)
				}

				if isPartOne {
					// After a monkey inspects an item, the worry level is divided by 3
					// and rounded down to the nearest integer
					worryLevel = int(math.Floor(float64(worryLevel) / 3))
				} else {
					worryLevel = worryLevel % dividend
				}

				// Carry out the monkey's test to determine which monkey
				// to throw the item to
				nextItemLocation := -1
				testPassed := (worryLevel % monkey.test.Delta) == 0
				if testPassed {
					nextItemLocation = monkey.test.OnPass
				} else {
					nextItemLocation = monkey.test.OnFail
				}

				// Throw the item with the new worry level to the next monkey
				nextMonkey := monkeyMap[nextItemLocation]
				nextMonkey.items = append(nextMonkey.items, worryLevel)

				// Update the number of inspected items
				monkey.inspectedItems += 1
			}

			// Reset the monkey's items after each round as
			// it throws all items to the next monkey after testing
			monkey.items = []int{}
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspectedItems > monkeys[j].inspectedItems
	})

	monkeyBusiness := monkeys[0].inspectedItems * monkeys[1].inspectedItems
	fmt.Printf("Monkey business after %v rounds: %v\n", rounds, monkeyBusiness)
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getVisibleTreesAtEdge(rows [][]int) int {
	// Get the trees at the left edge
	var leftEdgeTrees []int
	for _, row := range rows {
		leftEdgeTree := row[0]
		leftEdgeTrees = append(leftEdgeTrees, leftEdgeTree)
	}

	// Get the trees at the top edge
	topEdgeTrees := rows[0][1:]

	// Get the trees at the right edge
	var rightEdgeTrees []int
	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		rightEdgeTree := row[len(row)-1]
		rightEdgeTrees = append(rightEdgeTrees, rightEdgeTree)
	}

	// Get the trees at the bottom edge
	bottomEdgeTrees := rows[len(rows)-1][1 : len(rows)-1]

	sumOfTreesAtEdge := len(leftEdgeTrees) + len(topEdgeTrees) + len(rightEdgeTrees) + len(bottomEdgeTrees)
	return sumOfTreesAtEdge
}

func getVisibleInteriorTrees(rows [][]int) int {
	visibleInteriorTrees := 0

	for rowIndex, row := range rows {
		// Skip the first row and last row as they
		// are the top and bottom edges
		if rowIndex == 0 || rowIndex == len(rows)-1 {
			continue
		}

		for treeIndex, tree := range row {
			// Skip the first tree and last tree as they are the left and right edges
			if treeIndex == 0 || treeIndex == len(row)-1 {
				continue
			}

			// Check if the trees to the left are smaller
			// to determine if the current tree is visible
			isVisibleFromLeft := true
			for i := treeIndex - 1; i >= 0; i-- {
				if row[i] >= tree {
					isVisibleFromLeft = false
					break
				}
			}

			// Check if the trees to the right are smaller
			isVisibleFromRight := true
			for i := treeIndex + 1; i < len(row); i++ {
				if row[i] >= tree {
					isVisibleFromRight = false
					break
				}
			}

			// Check if the trees above are smaller
			isVisibleFromTop := true
			for i := rowIndex - 1; i >= 0; i-- {
				if rows[i][treeIndex] >= tree {
					isVisibleFromTop = false
					break
				}
			}

			// Check if the trees below are smaller
			isVisibleFromBottom := true
			for i := rowIndex + 1; i < len(rows); i++ {
				if rows[i][treeIndex] >= tree {
					isVisibleFromBottom = false
					break
				}
			}

			if isVisibleFromLeft || isVisibleFromRight || isVisibleFromTop || isVisibleFromBottom {
				visibleInteriorTrees++
			}
		}
	}

	return visibleInteriorTrees
}

func getScenicScores(rows [][]int) []int {
	var scenicScores []int

	for rowIndex, row := range rows {
		// Skip the first row and last row as they
		// are the top and bottom edges
		if rowIndex == 0 || rowIndex == len(rows)-1 {
			continue
		}

		for treeIndex, tree := range row {
			// Skip the first tree and last tree as they are the left and right edges
			if treeIndex == 0 || treeIndex == len(row)-1 {
				continue
			}

			// Count the trees to the left are smaller or equal
			leftScenicScore := 0
			for i := treeIndex - 1; i >= 0; i-- {
				if row[i] >= tree {
					leftScenicScore++
					break
				}

				leftScenicScore++
			}

			// Count the trees to the right are smaller or equal
			rightScenicScore := 0
			for i := treeIndex + 1; i < len(row); i++ {
				if row[i] >= tree {
					rightScenicScore++
					break
				}

				rightScenicScore++
			}

			// Count the trees above are smaller or equal
			topScenicScore := 0
			for i := rowIndex - 1; i >= 0; i-- {
				if rows[i][treeIndex] >= tree {
					topScenicScore++
					break
				}

				topScenicScore++
			}

			// Check if the trees below are smaller
			bottomScenicScore := 0
			for i := rowIndex + 1; i < len(rows); i++ {
				if rows[i][treeIndex] >= tree {
					bottomScenicScore++
					break
				}

				bottomScenicScore++
			}

			scenicScore := leftScenicScore * rightScenicScore * topScenicScore * bottomScenicScore
			scenicScores = append(scenicScores, scenicScore)
		}
	}

	return scenicScores
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var rows [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var trees []int

		line := scanner.Text()
		splitTrees := strings.Split(line, "")

		for _, treeStr := range splitTrees {
			tree, err := strconv.Atoi(treeStr)
			if err != nil {
				log.Fatalf("failed to parse tree: %v", err)
			}

			trees = append(trees, tree)
		}

		rows = append(rows, trees)
	}

	sumOfTreesAtEdge := getVisibleTreesAtEdge(rows)
	visibleInteriorTrees := getVisibleInteriorTrees(rows)
	scenicScores := getScenicScores(rows)
	sort.Ints(scenicScores)

	fmt.Println("Sum of trees at edge: ", sumOfTreesAtEdge)
	fmt.Println("Visible interior trees: ", visibleInteriorTrees)

	fmt.Println("Total visible trees: ", visibleInteriorTrees+sumOfTreesAtEdge)
	fmt.Println("Highest scenic score: ", scenicScores[len(scenicScores)-1])
}

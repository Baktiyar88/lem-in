package main

import (
	"fmt"
	"lem-in/utils"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: go run ./cmd/main_new.go test_cases/example00.txt")
	}

	inputFile := os.Args[1]

	if !strings.HasSuffix(inputFile, ".txt") {
		log.Fatal("usage: go run ./cmd/main_new.go test_cases/example00.txt")
	}

	info, err := os.Stat(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if info.Size() == 0 {
		log.Fatal("File is empty")
	}

	graph, err := utils.ParseInput(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Используем объединенный алгоритм Дейкстры и Суурбалле
	fmt.Println("Using Dijkstra and Suurballe algorithms...")
	paths := utils.FindOptimalPaths(graph)

	if len(paths) == 0 {
		log.Fatal("No paths found from start to end.")
	}

	fmt.Printf("Found %d optimal paths\n", len(paths))

	// Создаем совместимый набор путей
	allCompatibleSets := utils.GetCompatiblePaths(paths)

	var bestPaths [][]*utils.Room
	var bestDistribution []int
	minTurns := -1

	for _, candidate := range allCompatibleSets {
		if len(candidate) == 0 {
			continue
		}
		distribution := utils.DistributeAnts(graph.AntCount, candidate)
		turns := 0
		for i := range candidate {
			t := len(candidate[i]) + distribution[i] - 2
			if t > turns {
				turns = t
			}
		}
		if minTurns == -1 || turns < minTurns {
			minTurns = turns
			bestPaths = candidate
			bestDistribution = distribution
		}
	}

	if len(bestPaths) == 0 {
		log.Fatal("No compatible paths found for simulation.")
	}

	utils.PrintInputFile(inputFile)
	fmt.Println("\nResult:")
	utils.SimulateAnts(bestPaths, bestDistribution)
}

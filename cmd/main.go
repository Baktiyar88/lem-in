package main

import (
	"fmt"
	"os"
	"strings"

	"lem-in/logic"
)

// main runs the simulation in CLI mode only. It requires a path to an input file.
// Читает путь к файлу из аргументов, печатает исходный ввод
// и результат симуляции (или ошибку) в требуемом формате.
func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "usage: lem-in <input_file>")
		os.Exit(1)
	}

	var input string
	var inputLines []string

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot read file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	input = string(data)
	inputLines = strings.Split(strings.TrimSpace(input), "\n")

	// Run the simulation
	result := logic.RunSimulation(input)
	if result.Error != "" {
		fmt.Println(result.Error)
		os.Exit(1)
	}

	// Print input lines
	for _, line := range inputLines {
		fmt.Println(line)
	}
	// Print blank line
	fmt.Println()
	// Print moves
	for _, move := range result.Output {
		fmt.Println(move)
	}
}

package utils

import (
	"fmt"
	"strings"
)

func FindPath(graph *Graph) [][]*Room {
	if graph.Start == nil || graph.End == nil {
		fmt.Println("Start or End room is not defined")
		return nil
	}

	var allPaths [][]*Room

	var dfs func(current *Room, target *Room, path []*Room, visited map[*Room]bool)
	dfs = func(current *Room, target *Room, path []*Room, visited map[*Room]bool) {

		path = append(path, current)

		if current == target {

			pathCopy := make([]*Room, len(path))
			copy(pathCopy, path)
			allPaths = append(allPaths, pathCopy)
			return
		}

		visited[current] = true

		for _, linkedRoom := range current.Links {
			if !visited[linkedRoom] {
				dfs(linkedRoom, target, path, visited)
			}
		}

		delete(visited, current)
	}

	visited := make(map[*Room]bool)
	dfs(graph.Start, graph.End, []*Room{}, visited)

	if len(allPaths) == 0 {
		fmt.Println("No path found from start to end.")
		return nil
	}

	return allPaths
}

func ArePathsCompatible(p1, p2 []*Room) bool {
	rooms := make(map[*Room]bool)
	for _, r := range p1[1 : len(p1)-1] {
		rooms[r] = true
	}
	for _, r := range p2[1 : len(p2)-1] {
		if rooms[r] {
			return false
		}
	}
	return true
}

func GetCompatiblePaths(allPaths [][]*Room) [][][]*Room {
	var result [][][]*Room

	var backtrack func(start int, current [][]*Room)
	backtrack = func(start int, current [][]*Room) {
		result = append(result, append([][]*Room{}, current...))
		for i := start; i < len(allPaths); i++ {
			ok := true
			for _, path := range current {
				if !ArePathsCompatible(path, allPaths[i]) {
					ok = false
					break
				}
			}
			if ok {
				backtrack(i+1, append(current, allPaths[i]))
			}
		}
	}
	backtrack(0, [][]*Room{})
	return result
}

func DistributeAnts(ants int, paths [][]*Room) []int {
	counts := make([]int, len(paths))
	if len(paths) == 0 {
		return counts
	}

	lengths := make([]int, len(paths))
	for i, path := range paths {
		lengths[i] = len(path)
	}

	for ants > 0 {
		minTurns := -1
		index := -1
		for i := range paths {
			turns := lengths[i] + counts[i]
			if index == -1 || turns < minTurns {
				minTurns = turns
				index = i
			}
		}
		if index == -1 {
			break // safety net
		}
		counts[index]++
		ants--
	}
	return counts
}

func SimulateAnts(paths [][]*Room, distribution []int) {
	totalAnts := 0
	for _, count := range distribution {
		totalAnts += count
	}

	antsInPaths := []Ant{}
	antCounter := 1

	occupied := map[*Room]bool{}

	for {
		line := []string{}
		occupied = map[*Room]bool{}

		for i := range antsInPaths {
			ant := &antsInPaths[i]
			if ant.Pos < len(ant.Path)-1 {
				nextRoom := ant.Path[ant.Pos+1]

				if nextRoom.IsEnd || !occupied[nextRoom] {
					ant.Pos++
					line = append(line, fmt.Sprintf("L%d-%s", ant.Id, nextRoom.Name))
					if !nextRoom.IsEnd {
						occupied[nextRoom] = true
					}
				} else {

					if ant.Pos > 0 && !ant.Path[ant.Pos].IsEnd && !ant.Path[ant.Pos].IsStart {
						occupied[ant.Path[ant.Pos]] = true
					}
				}
			}
		}

		for i, path := range paths {
			if distribution[i] > 0 {
				firstRoom := path[1]
				if !occupied[firstRoom] {
					antsInPaths = append(antsInPaths, Ant{
						Id:   antCounter,
						Path: path,
						Pos:  1,
					})
					line = append(line, fmt.Sprintf("L%d-%s", antCounter, firstRoom.Name))
					occupied[firstRoom] = true
					antCounter++
					distribution[i]--
				}
			}
		}

		if len(line) == 0 {
			break
		}

		fmt.Println(strings.Join(line, " "))
	}
}

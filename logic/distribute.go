package logic

import (
	"sort"
)

// calcTime возвращает минимальное число ходов (turns) для заданных путей
// и количества муравьёв. Формула учитывает выравнивание длин путей.
func calcTime(paths []Path, ants int) int {
	if len(paths) == 0 {
		return 0
	}
	lengths := make([]int, len(paths))
	minL, maxL := 0, 0
	for i, p := range paths {
		l := len(p) - 1
		lengths[i] = l
		if i == 0 || l < minL {
			minL = l
		}
		if l > maxL {
			maxL = l
		}
	}
	sum := 0
	for _, l := range lengths {
		sum += maxL - l
	}
	antsPer := maxL - minL + (ants-sum)/len(paths)
	if (ants-sum)%len(paths) > 0 {
		antsPer++
	}
	return minL + antsPer - 1
}

// calcTimeAndDistribute считает число ходов и распределяет муравьёв
// по путям, когда пути получены методом Суурбалле (сортировка по длине).
func calcTimeAndDistribute(paths []Path, ants int) (int, []int) {
	if len(paths) == 0 {
		return 0, []int{}
	}
	pathInfos := make([]struct {
		path  Path
		index int
	}, len(paths))
	for i, p := range paths {
		pathInfos[i] = struct {
			path  Path
			index int
		}{path: p, index: i}
	}
	sort.Slice(pathInfos, func(i, j int) bool {
		return len(pathInfos[i].path) < len(pathInfos[j].path)
	})
	lengths := make([]int, len(paths))
	for i, pi := range pathInfos {
		lengths[i] = len(pi.path) - 1
	}
	result := make([]int, len(paths))
	if lengths[0] == 1 {
		result[pathInfos[0].index] = ants
		return 1, result
	}
	maxLength := lengths[len(lengths)-1]
	remainingAnts := ants
	for i := 0; i < len(paths); i++ {
		baseAnts := maxLength - lengths[i]
		result[pathInfos[i].index] = baseAnts
		remainingAnts -= baseAnts
	}
	if remainingAnts > 0 {
		antsPerPath := remainingAnts / len(paths)
		remainder := remainingAnts % len(paths)
		for i := 0; i < len(paths); i++ {
			result[pathInfos[i].index] += antsPerPath
			if i < remainder {
				result[pathInfos[i].index]++
			}
		}
	}
	steps := lengths[0] + (ants+len(paths)-1)/len(paths) - 1
	return steps, result
}

// distributeAnts распределяет муравьёв по путям для набора путей из DFS
// при известном числе ходов (turns), стремясь сбалансировать нагрузки.
func distributeAnts(paths []Path, ants, turns int) []int {
	n := len(paths)
	counts := make([]int, n)
	for i, p := range paths {
		length := len(p) - 1
		c := turns - (length - 1)
		if c < 0 {
			c = 0
		}
		counts[i] = c
	}
	sum := 0
	for _, c := range counts {
		sum += c
	}
	for sum < ants {
		idx := 0
		for i := 1; i < n; i++ {
			if len(paths[i]) > len(paths[idx]) {
				idx = i
			}
		}
		counts[idx]++
		sum++
	}
	for sum > ants {
		idx := 0
		for i := 1; i < n; i++ {
			if counts[i] < counts[idx] {
				idx = i
			}
		}
		if counts[idx] > 0 {
			counts[idx]--
			sum--
		} else {
			break
		}
	}
	return counts
}

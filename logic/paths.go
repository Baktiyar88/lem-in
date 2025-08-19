package logic

import (
	"sort"
)

// searchShortPath находит кратчайший путь (по количеству рёбер)
// с использованием упрощённой идеи алгоритма Суурбалле.
func searchShortPath(g *Graph) (Path, bool) {
	q := &sortedQueue{}
	startRoom := g.Rooms[g.Start]
	endRoom := g.Rooms[g.End]
	for _, room := range g.Rooms {
		room.Visit = false
		room.Parent = nil
		room.Weight = 0
	}
	startRoom.Visit = true
	q.Enqueue(startRoom, 0)

	for len(q.items) > 0 && !endRoom.Visit {
		current := q.Dequeue()
		currentRoom := current.Room
		for _, nextName := range g.Links[currentRoom.Name] {
			next := g.Rooms[nextName]
			weight := 1
			if !next.Visit {
				if next.Separated && next != endRoom {
					continue
				}
				next.Visit = true
				next.Parent = currentRoom
				next.Weight = current.Weight + weight
				q.Enqueue(next, next.Weight)
			} else if current.Weight+weight < next.Weight {
				next.Parent = currentRoom
				next.Weight = current.Weight + weight
				q.Enqueue(next, next.Weight)
			}
		}
	}

	if !endRoom.Visit {
		return nil, false
	}

	path := []string{}
	r := endRoom
	for r != startRoom {
		path = append([]string{r.Name}, path...)
		r = r.Parent
	}
	path = append([]string{startRoom.Name}, path...)

	for i := 0; i < len(path)-1; i++ {
		from, to := path[i], path[i+1]
		g.Links[from] = removeLink(g.Links[from], to)
		g.Links[to] = removeLink(g.Links[to], from)
		if i > 0 && i < len(path)-1 {
			g.Rooms[path[i]].Separated = true
		}
	}

	return path, true
}

func removeLink(links []string, target string) []string {
	for i, link := range links {
		if link == target {
			return append(links[:i], links[i+1:]...)
		}
	}
	return links
}

// findDisjointPaths ищет набор попарно непересекающихся путей
// с помощью последовательного применения searchShortPath.
func findDisjointPaths(g *Graph) []Path {
	var paths []Path
	gCopy := &Graph{
		Rooms:   make(map[string]*Room),
		Links:   make(map[string][]string),
		Start:   g.Start,
		End:     g.End,
		NumAnts: g.NumAnts,
	}
	for name, room := range g.Rooms {
		gCopy.Rooms[name] = &Room{Name: room.Name, X: room.X, Y: room.Y}
	}
	for name, links := range g.Links {
		gCopy.Links[name] = make([]string, len(links))
		copy(gCopy.Links[name], links)
	}

	for {
		path, found := searchShortPath(gCopy)
		if !found {
			break
		}
		paths = append(paths, path)
	}
	return paths
}

// dfsPaths собирает все простые пути от старта к финишу (DFS),
// затем сортирует их по длине по возрастанию.
func dfsPaths(g *Graph) []Path {
	start, end := g.Start, g.End
	visited := make(map[string]bool)
	var current Path
	var all []Path
	var dfs func(string)
	dfs = func(node string) {
		if node == end {
			tmp := make(Path, len(current)+1)
			copy(tmp, current)
			tmp[len(current)] = node
			all = append(all, tmp)
			return
		}
		visited[node] = true
		current = append(current, node)
		for _, nb := range g.Links[node] {
			if !visited[nb] {
				dfs(nb)
			}
		}
		visited[node] = false
		current = current[:len(current)-1]
	}
	dfs(start)
	// Стабильная сортировка по длине сохраняет порядок обнаружения
	// (зависит от порядка рёбер во входе) для путей одинаковой длины.
	sort.SliceStable(all, func(i, j int) bool { return len(all[i]) < len(all[j]) })
	return all
}

// choosePathsDFS перебирает комбинации путей и выбирает набор,
// минимизирующий число ходов (turns) при заданном числе муравьёв.
func choosePathsDFS(paths []Path, ants int) []Path {
	if len(paths) == 0 {
		return nil
	}
	best := []Path{paths[0]}
	bestTime := calcTime(best, ants)
	n := len(paths)
	for mask := 1; mask < (1 << n); mask++ {
		var comb []Path
		used := make(map[string]bool)
		valid := true
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			p := paths[i]
			for j := 1; j < len(p)-1; j++ {
				if used[p[j]] {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
			for j := 1; j < len(p)-1; j++ {
				used[p[j]] = true
			}
			comb = append(comb, p)
		}
		if valid && len(comb) > 0 {
			t := calcTime(comb, ants)
			if t < bestTime {
				bestTime = t
				best = comb
			}
		}
	}
	return best
}

// choosePathsHybrid выбирает стратегию: для сложных графов
// пытается Суурбалле, в противном случае — DFS, затем ищет
// оптимальную комбинацию путей.
func choosePathsHybrid(g *Graph, ants int) []Path {
	// Простая стратегия: для небольших входов предпочитаем DFS,
	// для крупных/плотных графов или большого числа муравьёв — Суурбалле.
	linkCount := 0
	for _, links := range g.Links {
		linkCount += len(links)
	}
	complexity := float64(linkCount) / float64(len(g.Rooms))

	if ants > 100 || len(g.Rooms) > 30 || complexity > 2.5 {
		paths := findDisjointPaths(g)
		if len(paths) > 0 {
			return paths
		}
	}

	// По умолчанию — DFS со стабильной сортировкой и выбором лучшей комбинации
	paths := dfsPaths(g)
	if len(paths) == 0 {
		return nil
	}
	return choosePathsDFS(paths, ants)
}

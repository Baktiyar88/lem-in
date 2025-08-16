package utils

import (
	"container/heap"
	"fmt"
	"strings"
)

// PriorityQueue для алгоритма Дейкстры
type PriorityQueue []*PathNode

type PathNode struct {
	room     *Room
	priority int
	path     []*Room
	index    int
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*PathNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// DijkstraShortestPath находит кратчайший путь от start до end
func DijkstraShortestPath(graph *Graph) []*Room {
	if graph.Start == nil || graph.End == nil {
		return nil
	}

	// Инициализация
	distances := make(map[*Room]int)
	visited := make(map[*Room]bool)

	for _, room := range graph.Rooms {
		distances[room] = 1<<31 - 1 // Бесконечность
	}
	distances[graph.Start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PathNode{
		room:     graph.Start,
		priority: 0,
		path:     []*Room{graph.Start},
	})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*PathNode)

		if visited[current.room] {
			continue
		}
		visited[current.room] = true

		if current.room == graph.End {
			return current.path
		}

		for _, neighbor := range current.room.Links {
			if visited[neighbor] {
				continue
			}

			newDistance := distances[current.room] + 1
			if newDistance < distances[neighbor] {
				distances[neighbor] = newDistance

				newPath := make([]*Room, len(current.path))
				copy(newPath, current.path)
				newPath = append(newPath, neighbor)

				heap.Push(pq, &PathNode{
					room:     neighbor,
					priority: newDistance,
					path:     newPath,
				})
			}
		}
	}

	return nil
}

// SuurballeAlgorithm находит два непересекающихся пути с минимальной общей длиной
func SuurballeAlgorithm(graph *Graph) ([][]*Room, error) {
	if graph.Start == nil || graph.End == nil {
		return nil, fmt.Errorf("start or end room not defined")
	}

	// Шаг 1: Находим кратчайший путь P1
	path1 := DijkstraShortestPath(graph)
	if path1 == nil {
		return nil, fmt.Errorf("no path found from start to end")
	}

	// Шаг 2: Создаем обратный граф и изменяем веса
	reversedGraph := createReversedGraph(graph, path1)

	// Шаг 3: Находим кратчайший путь в обратном графе
	path2 := DijkstraShortestPath(reversedGraph)
	if path2 == nil {
		return nil, fmt.Errorf("no second path found")
	}

	// Шаг 4: Удаляем общие ребра и объединяем пути
	resultPaths := removeCommonEdges(path1, path2)

	return resultPaths, nil
}

// createReversedGraph создает обратный граф с измененными весами
func createReversedGraph(originalGraph *Graph, shortestPath []*Room) *Graph {
	// Создаем новый граф
	newGraph := &Graph{
		AntCount: originalGraph.AntCount,
		Rooms:    make(map[string]*Room),
	}

	// Копируем все комнаты
	for name, room := range originalGraph.Rooms {
		newRoom := &Room{
			Name:    room.Name,
			X:       room.X,
			Y:       room.Y,
			IsStart: room.IsStart,
			IsEnd:   room.IsEnd,
			Links:   []*Room{},
		}
		newGraph.Rooms[name] = newRoom
	}

	// Устанавливаем указатели на start и end
	newGraph.Start = newGraph.Rooms[originalGraph.Start.Name]
	newGraph.End = newGraph.Rooms[originalGraph.End.Name]

	// Создаем множество ребер из кратчайшего пути
	pathEdges := make(map[string]bool)
	for i := 0; i < len(shortestPath)-1; i++ {
		edge := getEdgeKey(shortestPath[i], shortestPath[i+1])
		pathEdges[edge] = true
	}

	// Добавляем обратные ребра и изменяем веса
	for name, room := range originalGraph.Rooms {
		newRoom := newGraph.Rooms[name]
		for _, linkedRoom := range room.Links {
			edge := getEdgeKey(room, linkedRoom)
			reverseEdge := getEdgeKey(linkedRoom, room)

			if pathEdges[edge] {
				// Ребро в кратчайшем пути - добавляем обратное ребро с весом 0
				newLinkedRoom := newGraph.Rooms[linkedRoom.Name]
				newRoom.Links = append(newRoom.Links, newLinkedRoom)
			} else if !pathEdges[reverseEdge] {
				// Обычное ребро - добавляем как есть
				newLinkedRoom := newGraph.Rooms[linkedRoom.Name]
				newRoom.Links = append(newRoom.Links, newLinkedRoom)
			}
		}
	}

	return newGraph
}

// getEdgeKey создает уникальный ключ для ребра
func getEdgeKey(room1, room2 *Room) string {
	if room1.Name < room2.Name {
		return room1.Name + "-" + room2.Name
	}
	return room2.Name + "-" + room1.Name
}

// removeCommonEdges удаляет общие ребра и объединяет пути
func removeCommonEdges(path1, path2 []*Room) [][]*Room {
	// Создаем множества ребер для каждого пути
	edges1 := make(map[string]bool)
	edges2 := make(map[string]bool)

	for i := 0; i < len(path1)-1; i++ {
		edge := getEdgeKey(path1[i], path1[i+1])
		edges1[edge] = true
	}

	for i := 0; i < len(path2)-1; i++ {
		edge := getEdgeKey(path2[i], path2[i+1])
		edges2[edge] = true
	}

	// Находим общие ребра
	commonEdges := make(map[string]bool)
	for edge := range edges1 {
		if edges2[edge] {
			commonEdges[edge] = true
		}
	}

	// Создаем новые пути без общих ребер
	newPath1 := removeEdgesFromPath(path1, commonEdges)
	newPath2 := removeEdgesFromPath(path2, commonEdges)

	return [][]*Room{newPath1, newPath2}
}

// removeEdgesFromPath удаляет указанные ребра из пути
func removeEdgesFromPath(path []*Room, edgesToRemove map[string]bool) []*Room {
	if len(path) <= 1 {
		return path
	}

	var result []*Room
	result = append(result, path[0])

	for i := 0; i < len(path)-1; i++ {
		edge := getEdgeKey(path[i], path[i+1])
		if !edgesToRemove[edge] {
			result = append(result, path[i+1])
		}
	}

	// Убеждаемся, что путь содержит хотя бы start и end
	if len(result) < 2 {
		return path // Возвращаем оригинальный путь если что-то пошло не так
	}

	return result
}

// FindOptimalPaths - основная функция для поиска оптимальных путей
func FindOptimalPaths(graph *Graph) [][]*Room {
	// Сначала пробуем алгоритм Суурбалле для поиска непересекающихся путей
	suurballePaths, err := SuurballeAlgorithm(graph)
	if err == nil && len(suurballePaths) > 0 {
		// Проверяем, что пути содержат хотя бы 2 комнаты
		validPaths := [][]*Room{}
		for _, path := range suurballePaths {
			if len(path) > 1 {
				validPaths = append(validPaths, path)
			}
		}
		if len(validPaths) > 0 {
			return validPaths
		}
	}

	// Если алгоритм Суурбалле не сработал, используем Дейкстру
	dijkstraPath := DijkstraShortestPath(graph)
	if dijkstraPath != nil {
		return [][]*Room{dijkstraPath}
	}

	// Если ничего не найдено, возвращаем пустой массив
	return [][]*Room{}
}

// ArePathsCompatible проверяет совместимость путей (не пересекаются)
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

// GetCompatiblePaths находит все совместимые наборы путей
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

// DistributeAnts распределяет муравьев по путям оптимально
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

// SimulateAnts симулирует движение муравьев по путям
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

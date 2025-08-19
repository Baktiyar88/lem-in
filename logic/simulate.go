package logic

import (
	"fmt"
	"sort"
	"strings"
)

// RunSimulation — входная точка движка. Парсит вход, выбирает пути
// и выполняет пошаговую симуляцию перемещения муравьёв.
func RunSimulation(input string) Response {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	g, err := parseLines(lines)
	if err != nil {
		return Response{Error: "ERROR: invalid data format"}
	}
	paths := choosePathsHybrid(g, g.NumAnts)
	if len(paths) == 0 {
		return Response{Error: "ERROR: no valid paths found"}
	}
	isSuurballe := float64(len(g.Links))/float64(len(g.Rooms)) > 2.0 || g.NumAnts > 20
	moves := moveAnts(paths, g.NumAnts, isSuurballe)
	return Response{Output: moves}
}

// moveAnts выполняет пошаговую симуляцию и возвращает срез строк ходов
// в формате L<id>-<room> для каждого шага, объединённых пробелами.
func moveAnts(paths []Path, ants int, isSuurballe bool) []string {
	if len(paths) == 0 {
		return nil
	}

	var counts []int
	if isSuurballe {
		_, counts = calcTimeAndDistribute(paths, ants)
	} else {
		turns := calcTime(paths, ants)
		counts = distributeAnts(paths, ants, turns)
	}

	type antOnPath struct {
		id      int
		pathIdx int
		pos     int
	}

	var active []antOnPath
	moves := []string{}
	nextID := 1
	finished := 0

	totalAntsToLaunch := 0
	for _, count := range counts {
		totalAntsToLaunch += count
	}

	if totalAntsToLaunch > ants {
		scale := float64(ants) / float64(totalAntsToLaunch)
		remainder := ants
		for i := range counts {
			counts[i] = int(float64(counts[i]) * scale)
			remainder -= counts[i]
		}
		for i := 0; i < remainder && i < len(counts); i++ {
			counts[i]++
		}
	}

	for finished < ants {
		roomOcc := make(map[string]bool)
		var step []string
		var nextActive []antOnPath

		// Движение активных муравьев
		for _, a := range active {
			path := paths[a.pathIdx]
			nextPos := a.pos + 1
			if nextPos < len(path) {
				room := path[nextPos]
				if nextPos == len(path)-1 || !roomOcc[room] {
					if a.pos > 0 {
						roomOcc[path[a.pos]] = false
					}
					if nextPos != len(path)-1 {
						roomOcc[room] = true
					}
					a.pos = nextPos
					step = append(step, fmt.Sprintf("L%d-%s", a.id, room))
					if nextPos == len(path)-1 {
						finished++
					} else {
						nextActive = append(nextActive, a)
					}
				} else {
					nextActive = append(nextActive, a)
				}
			}
		}
		active = nextActive

		// Запуск новых муравьев
		for i, path := range paths {
			if counts[i] == 0 || nextID > ants {
				continue
			}
			if len(path) < 2 {
				continue
			}
			room := path[1]
			// Для прямого пути: если муравьев <= 2, отправляем всех сразу, иначе по одному
			if len(path) == 2 {
				if ants <= 2 {
					// Отправляем всех муравьев в один шаг
					for counts[i] > 0 && nextID <= ants {
						step = append(step, fmt.Sprintf("L%d-%s", nextID, room))
						finished++
						counts[i]--
						nextID++
					}
				} else if !roomOcc[room] {
					// Отправляем одного муравья за шаг
					step = append(step, fmt.Sprintf("L%d-%s", nextID, room))
					roomOcc[room] = true
					finished++
					counts[i]--
					nextID++
				}
			} else if !roomOcc[room] {
				step = append(step, fmt.Sprintf("L%d-%s", nextID, room))
				if len(path) > 2 {
					roomOcc[room] = true
					active = append(active, antOnPath{id: nextID, pathIdx: i, pos: 1})
				} else {
					finished++
				}
				counts[i]--
				nextID++
			}
		}

		if len(step) > 0 {
			sort.Slice(step, func(i, j int) bool {
				var ni, nj int
				fmt.Sscanf(step[i], "L%d-", &ni)
				fmt.Sscanf(step[j], "L%d-", &nj)
				return ni < nj
			})
			moves = append(moves, strings.Join(step, " "))
		}

		// Если все муравьи достигли конца, выходим
		if finished == ants {
			break
		}

		// Защита от бесконечного цикла
		if len(moves) > ants*100 {
			break
		}
	}

	return moves
}

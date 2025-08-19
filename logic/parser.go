package logic

import (
	"fmt"
	"strconv"
	"strings"
)

// parseLines парсит входные строки в структуру Graph.
// Поддерживаются комментарии, директивы \"##start\"/\"##end\",
// декларации комнат и рёбер. Валидирует формат и обязательные сущности.
func parseLines(lines []string) (*Graph, error) {
	g := &Graph{
		Rooms: make(map[string]*Room),
		Links: make(map[string][]string),
		Input: lines,
	}
	antsParsed := false
	parsingRooms := true
	countstart := 0
	countend := 0

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Пропускаем пустые строки
		if line == "" {
			continue
		}

		// Обрабатываем специальные команды ##start и ##end
		if line == "##start" || line == "##end" {
			if line == "##start" {
				countstart++
			} else if line == "##end" {
				countend++
			}
			continue
		}

		// Пропускаем обычные комментарии (начинающиеся с #, но не ##start/##end)
		if strings.HasPrefix(line, "#") {
			continue
		}

		if !antsParsed {
			n, err := strconv.Atoi(line)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid number of ants at line %d: %s", i+1, line)
			}
			g.NumAnts = n
			antsParsed = true
			continue
		}

		if strings.Contains(line, "-") {
			parsingRooms = false
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid link format at line %d: %s", i+1, line)
			}
			a, b := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if a == b {
				return nil, fmt.Errorf("invalid link format, self-loop detected at line %d: %s", i+1, line)
			}
			if _, okA := g.Rooms[a]; !okA {
				return nil, fmt.Errorf("invalid link format, room %s not found at line %d: %s", a, i+1, line)
			}
			if _, okB := g.Rooms[b]; !okB {
				return nil, fmt.Errorf("invalid link format, room %s not found at line %d: %s", b, i+1, line)
			}
			g.Links[a] = append(g.Links[a], b)
			g.Links[b] = append(g.Links[b], a)
		} else if parsingRooms {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				// Если это не комната (неправильный формат), пропускаем
				continue
			}

			name := fields[0]

			// ВАЖНО: Проверяем имя комнаты ПЕРЕД парсингом координат
			if strings.HasPrefix(name, "L") {
				return nil, fmt.Errorf("invalid room name at line %d: room name cannot start with 'L': %s", i+1, name)
			}
			if strings.HasPrefix(name, "#") {
				return nil, fmt.Errorf("invalid room name at line %d: room name cannot start with '#': %s", i+1, name)
			}
			if strings.Contains(name, " ") {
				return nil, fmt.Errorf("invalid room name at line %d: room name cannot contain spaces: %s", i+1, name)
			}

			x, err1 := strconv.Atoi(fields[1])
			y, err2 := strconv.Atoi(fields[2])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid room coordinates at line %d: %s", i+1, line)
			}
			g.Rooms[name] = &Room{Name: name, X: x, Y: y}
		}
	}

	foundStart, foundEnd := false, false
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "##start" && i+1 < len(lines) {
			nxt := strings.TrimSpace(lines[i+1])
			if !strings.HasPrefix(nxt, "#") && strings.Contains(nxt, " ") {
				g.Start = strings.Fields(nxt)[0]
				foundStart = true
			}
		} else if line == "##end" && i+1 < len(lines) {
			nxt := strings.TrimSpace(lines[i+1])
			if !strings.HasPrefix(nxt, "#") && strings.Contains(nxt, " ") {
				g.End = strings.Fields(nxt)[0]
				foundEnd = true
			}
		}
	}

	if !antsParsed {
		return nil, fmt.Errorf("missing number of ants")
	}
	if !foundStart || !foundEnd {
		return nil, fmt.Errorf("missing start or end")
	}
	if countstart != 1 || countend != 1 {
		return nil, fmt.Errorf("exactly one start and one end room are required")
	}
	if _, ok := g.Rooms[g.Start]; !ok {
		return nil, fmt.Errorf("start room not declared")
	}
	if _, ok := g.Rooms[g.End]; !ok {
		return nil, fmt.Errorf("end room not declared")
	}
	return g, nil
}

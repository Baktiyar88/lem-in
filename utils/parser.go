package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ParseInput(filename string) (*Graph, error) {

	inputfile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer inputfile.Close()

	scanner := bufio.NewScanner(inputfile)

	graph := &Graph{
		Rooms: make(map[string]*Room),
	}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		antCount, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("ERROR:invalid ant count: %s", err)
		}

		if antCount <= 0 {
			return nil, fmt.Errorf("ERROR:ant count must be greater than 0")
		}

		graph.AntCount = antCount
		break
	}

	var nextIsStart, nextIsEnd bool

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				nextIsStart = true
			} else if line == "##end" {
				nextIsEnd = true
			}
			continue
		}

		if !strings.Contains(line, "-") {
			room, err := parseRoom(line)
			if err != nil {
				return nil, err
			}
			graph.Rooms[room.Name] = room
			room.IsStart = nextIsStart
			room.IsEnd = nextIsEnd
			if nextIsStart {
				graph.Start = room
				nextIsStart = false
			}
			if nextIsEnd {
				graph.End = room
				nextIsEnd = false
			}
		} else {
			err := parseLink(line, graph)
			if err != nil {
				return nil, err
			}
		}
	}

	if graph.Start == nil {
		return nil, fmt.Errorf("no start room found")
	}
	if graph.End == nil {
		return nil, fmt.Errorf("no end room found")
	}

	return graph, nil

}

func parseRoom(line string) (*Room, error) {

	roomsFields := strings.Fields(line)

	if strings.HasPrefix(roomsFields[0], "L") || strings.HasSuffix(roomsFields[0], "#") {
		return nil, fmt.Errorf("invalid room name: %s", roomsFields[0])
	}

	if len(roomsFields) != 3 {
		return nil, fmt.Errorf("invalid rooms: %s", roomsFields)
	}

	room := &Room{}

	xCoordinate, err := strconv.Atoi(roomsFields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid x coordination: %s", err)
	}

	yCoordinate, err := strconv.Atoi(roomsFields[2])
	if err != nil {
		return nil, fmt.Errorf("invalid y coordination: %s", err)
	}

	room.Name = roomsFields[0]
	room.X = xCoordinate
	room.Y = yCoordinate

	return room, nil

}

func parseLink(line string, graph *Graph) error {

	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid link format: %s", line)
	}

	room1, exists1 := graph.Rooms[parts[0]]
	room2, exists2 := graph.Rooms[parts[1]]

	if !exists1 || !exists2 {
		return fmt.Errorf("room not found in link: %s", line)
	}

	room1.Links = append(room1.Links, room2)
	room2.Links = append(room2.Links, room1)

	return nil
}

func PrintInputFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Package logic реализует ядро симуляции задачи "lem-in":
// парсинг входных данных, поиск путей и пошаговую симуляцию
// перемещения муравьёв для формирования корректного вывода.
package logic

import (
	"sort"
)

// Response содержит либо последовательность ходов симуляции, либо текст ошибки.
type Response struct {
	Error  string
	Output []string
}

// Room описывает вершину графа с координатами и служебными полями
// для алгоритмов поиска путей (посещённость, родитель, вес, разделение).
type Room struct {
	Name      string
	X, Y      int
	Visit     bool  // For Suurballe's algorithm
	Parent    *Room // For path reconstruction
	Weight    int   // For path weights
	Separated bool  // For marking rooms in disjoint paths
}

// Graph хранит распарсенную конфигурацию муравейника: комнаты, связи,
// старт/финиш, число муравьёв и исходные строки ввода.
type Graph struct {
	Rooms   map[string]*Room    // name -> room
	Links   map[string][]string // adjacency list
	Start   string              // name of start room
	End     string              // name of end room
	NumAnts int                 // number of ants
	Input   []string            // raw input lines (trimmed)
}

// Path — последовательность имён комнат от старта к финишу.
type Path []string

// Priority queue for Suurballe's algorithm
type queueEntry struct {
	Room   *Room
	Weight int
}

type sortedQueue struct {
	items []queueEntry
}

func (q *sortedQueue) Enqueue(room *Room, weight int) {
	q.items = append(q.items, queueEntry{Room: room, Weight: weight})
	sort.Slice(q.items, func(i, j int) bool { return q.items[i].Weight < q.items[j].Weight })
}

func (q *sortedQueue) Dequeue() *queueEntry {
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return &item
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Data structures for parsing
type Room struct {
	Name string
	X    int
	Y    int
}

type InputData struct {
	Ants     int
	Rooms    map[string]Room
	Start    string
	End      string
	Links    [][2]string
	RawLines []string
}

// Graph for max-flow using node splitting (vertex capacity = 1)
type edge struct {
	to   int
	rev  int
	cap  int
	room string // destination room name for link edges; empty for internal in->out edges
}

type Dinic struct {
	n   int
	g   [][]edge
	lvl []int
	it  []int
}

func newDinic(n int) *Dinic {
	g := make([][]edge, n)
	return &Dinic{n: n, g: g, lvl: make([]int, n), it: make([]int, n)}
}

func (d *Dinic) addEdge(u, v, c int, roomName string) {
	d.g[u] = append(d.g[u], edge{to: v, rev: len(d.g[v]), cap: c, room: roomName})
	d.g[v] = append(d.g[v], edge{to: u, rev: len(d.g[u]) - 1, cap: 0, room: roomName})
}

func (d *Dinic) bfs(s, t int) bool {
	for i := range d.lvl {
		d.lvl[i] = -1
	}
	q := make([]int, 0, d.n)
	d.lvl[s] = 0
	q = append(q, s)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.lvl[e.to] < 0 {
				d.lvl[e.to] = d.lvl[v] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.lvl[t] >= 0
}

func (d *Dinic) dfs(v, t, f int) int {
	if v == t {
		return f
	}
	for ; d.it[v] < len(d.g[v]); d.it[v]++ {
		i := d.it[v]
		e := &d.g[v][i]
		if e.cap > 0 && d.lvl[v] < d.lvl[e.to] {
			dmin := f
			if e.cap < dmin {
				dmin = e.cap
			}
			ret := d.dfs(e.to, t, dmin)
			if ret > 0 {
				e.cap -= ret
				rev := &d.g[e.to][e.rev]
				rev.cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *Dinic) maxFlow(s, t int) int {
	flow := 0
	for d.bfs(s, t) {
		for i := range d.it {
			d.it[i] = 0
		}
		for {
			f := d.dfs(s, t, 1<<30)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

// Parsing functions
func parseInput(path string) (*InputData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("ERROR: invalid data format")
	}

	idx := 0
	// Read ants count
	antsLine := strings.TrimSpace(lines[idx])
	idx++
	ants, err := strconv.Atoi(antsLine)
	if err != nil || ants < 0 {
		return nil, errors.New("ERROR: invalid data format")
	}

	rooms := make(map[string]Room)
	var startName, endName string
	links := make([][2]string, 0)

	expectStart := false
	expectEnd := false

	// Remaining lines
	for ; idx < len(lines); idx++ {
		raw := lines[idx]
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				expectStart = true
				expectEnd = false
			} else if line == "##end" {
				expectEnd = true
				expectStart = false
			}
			continue
		}
		// link or room
		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			// link
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, errors.New("ERROR: invalid data format")
			}
			a, b := parts[0], parts[1]
			if a == b {
				return nil, errors.New("ERROR: invalid data format")
			}
			links = append(links, [2]string{a, b})
			continue
		}
		// room line: name x y
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, errors.New("ERROR: invalid data format")
		}
		name := fields[0]
		if name == "" || strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") || strings.Contains(name, " ") {
			return nil, errors.New("ERROR: invalid data format")
		}
		x, err1 := strconv.Atoi(fields[1])
		y, err2 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil {
			return nil, errors.New("ERROR: invalid data format")
		}
		if _, exists := rooms[name]; exists {
			return nil, errors.New("ERROR: invalid data format")
		}
		rooms[name] = Room{Name: name, X: x, Y: y}
		if expectStart {
			startName = name
			expectStart = false
		} else if expectEnd {
			endName = name
			expectEnd = false
		}
	}

	if startName == "" || endName == "" {
		return nil, errors.New("ERROR: invalid data format")
	}
	// Validate links refer to known rooms and no duplicates
	type pair struct{ a, b string }
	seen := map[pair]struct{}{}
	uniq := make([][2]string, 0, len(links))
	for _, l := range links {
		a, b := l[0], l[1]
		if _, ok := rooms[a]; !ok {
			return nil, errors.New("ERROR: invalid data format")
		}
		if _, ok := rooms[b]; !ok {
			return nil, errors.New("ERROR: invalid data format")
		}
		p := pair{a, b}
		q := pair{b, a}
		if _, ok := seen[p]; ok {
			continue
		}
		if _, ok := seen[q]; ok {
			continue
		}
		seen[p] = struct{}{}
		uniq = append(uniq, [2]string{a, b})
	}

	return &InputData{Ants: ants, Rooms: rooms, Start: startName, End: endName, Links: uniq, RawLines: lines}, nil
}

// Build flow network with node splitting for all rooms except start and end.
// For each room R, create R_in and R_out with capacity 1 edge from in->out
// For start and end, do not split and allow multiple ants.
type graphIndex struct {
	in  int
	out int
}

type Network struct {
	D         *Dinic
	RoomToIdx map[string]graphIndex
	Start     int
	End       int
}

func buildNetwork(inp *InputData) *Network {
	// Assign indices
	// Each normal room uses 2 nodes (in/out), start and end use 1
	roomToIdx := make(map[string]graphIndex)
	index := 0
	// pre-assign start and end single nodes
	startNode := index
	index++
	endNode := index
	index++

	// Assign for other rooms except start/end
	for name := range inp.Rooms {
		if name == inp.Start || name == inp.End {
			continue
		}
		in := index
		out := index + 1
		index += 2
		roomToIdx[name] = graphIndex{in: in, out: out}
	}

	// Build Dinic
	d := newDinic(index)

	// Add capacity 1 in->out for normal rooms
	for _, gi := range roomToIdx {
		d.addEdge(gi.in, gi.out, 1, "")
	}

	// Map start/end
	// For convenience, keep roomToIdx entries for start/end pointing to Start/End as both in and out
	roomToIdx[inp.Start] = graphIndex{in: startNode, out: startNode}
	roomToIdx[inp.End] = graphIndex{in: endNode, out: endNode}

	// Add edges for links. For an undirected link A-B, add directed edges both ways with capacity 1.
	// From A_out to B_in and B_out to A_in. For start/end, use their single nodes.
	for _, l := range inp.Links {
		a, b := l[0], l[1]
		ai := roomToIdx[a]
		bi := roomToIdx[b]
		// define from out to in
		var aOut int
		var bOut int
		var aIn int
		var bIn int
		if a == inp.Start || a == inp.End {
			aOut = ai.out
			aIn = ai.in
		} else {
			aOut = ai.out
			aIn = ai.in
		}
		if b == inp.Start || b == inp.End {
			bOut = bi.out
			bIn = bi.in
		} else {
			bOut = bi.out
			bIn = bi.in
		}
		d.addEdge(aOut, bIn, 1, b)
		d.addEdge(bOut, aIn, 1, a)
	}

	return &Network{D: d, RoomToIdx: roomToIdx, Start: startNode, End: endNode}
}

// Extract disjoint paths from residual network after max flow by following edges with reverse capacity > 0
func extractPaths(net *Network, inp *InputData) [][]string {
	// Consume flow edges (reverse.cap > 0) to reconstruct disjoint paths.
	paths := [][]string{}

	// Helper to try build one path and consume it
	for {
		// visited nodes to avoid loops while tracing single path
		visited := make([]bool, net.D.n)
		seq := []string{inp.Start}
		cur := net.Start
		ok := false
		for {
			if cur == net.End {
				ok = true
				break
			}
			visited[cur] = true
			advanced := false
			for i := range net.D.g[cur] {
				e := &net.D.g[cur][i]
				// If there is flow on cur->e.to, the reverse edge must have positive capacity
				rev := &net.D.g[e.to][e.rev]
				if rev.cap > 0 {
					// Consume this unit of flow for path extraction
					rev.cap -= 1
					e.cap += 1
					// Append room name and advance
					if e.room != "" && (len(seq) == 0 || seq[len(seq)-1] != e.room) {
						seq = append(seq, e.room)
					}
					cur = e.to
					advanced = true
					break
				}
			}
			if !advanced {
				break
			}
		}
		if !ok {
			break
		}
		paths = append(paths, seq)
	}
	return paths
}

// Distribute ants across paths to minimize makespan: assign ants to paths by increasing (len-1 + assigned)
func distributeAnts(numAnts int, paths [][]string) []int {
	n := len(paths)
	assign := make([]int, n)
	if n == 0 {
		return assign
	}
	// Precompute path lengths in edges (rooms - 1)
	lens := make([]int, n)
	for i, p := range paths {
		lens[i] = len(p) - 1
	}
	// Greedy allocation
	for a := 0; a < numAnts; a++ {
		best := 0
		bestVal := lens[0] + assign[0]
		for i := 1; i < n; i++ {
			v := lens[i] + assign[i]
			if v < bestVal {
				best = i
				bestVal = v
			}
		}
		assign[best]++
	}
	return assign
}

// Simulate ant movements step-by-step enforcing one ant per room and per tunnel per turn
func simulate(ants int, paths [][]string, assign []int) [][][2]int {
	// For output, we keep for each turn a list of (antId, destIndexInPath)
	// We'll map to room names later
	moves := [][][2]int{}

	type antState struct {
		pathIdx int
		posIdx  int // index in path, start at 0 (start room), goal is last index
	}

	// Build ants along paths
	states := make([]antState, 0, ants)
	antToPath := make([]int, 0, ants)
	for pi, count := range assign {
		for i := 0; i < count; i++ {
			states = append(states, antState{pathIdx: pi, posIdx: 0})
			antToPath = append(antToPath, pi)
		}
	}

	// For deterministic output, sort ants by path length then id grouping
	sort.Slice(states, func(i, j int) bool {
		li := len(paths[states[i].pathIdx])
		lj := len(paths[states[j].pathIdx])
		if li == lj {
			return states[i].pathIdx < states[j].pathIdx
		}
		return li < lj
	})

	// Room occupancy map excluding start/end
	occupied := map[string]int{} // room -> antId currently occupying

	// Start queue per path for ants not yet entered the first edge
	// But our states all start at position 0 and will move when first room available

	// Continue until all ants reached end
	finished := 0
	stateIdxToAntId := make([]int, len(states))
	for i := range states {
		stateIdxToAntId[i] = i + 1
	}

	for finished < ants {
		turnMoves := [][2]int{}
		// Update in reverse order along each path to avoid blocking within same turn
		// For each path, attempt to move ants from end to start
		// Build list of indices per path
		perPath := make([][]int, len(paths))
		for si, st := range states {
			if st.posIdx < len(paths[st.pathIdx])-1 { // not finished
				perPath[st.pathIdx] = append(perPath[st.pathIdx], si)
			}
		}
		for pi := range perPath {
			// Order ants on this path by decreasing posIdx
			sort.Slice(perPath[pi], func(a, b int) bool {
				return states[perPath[pi][a]].posIdx > states[perPath[pi][b]].posIdx
			})
			for _, si := range perPath[pi] {
				st := states[si]
				curIdx := st.posIdx
				nextIdx := curIdx + 1
				if nextIdx >= len(paths[pi]) {
					continue
				}
				destRoom := paths[pi][nextIdx]
				// End room is unlimited
				if destRoom != paths[pi][len(paths[pi])-1] {
					if _, ok := occupied[destRoom]; ok {
						continue
					}
				}
				// Free to move
				// Vacate current room if not start and not end
				if curIdx > 0 && paths[pi][curIdx] != paths[pi][len(paths[pi])-1] {
					delete(occupied, paths[pi][curIdx])
				}
				// Occupy dest if not end
				if destRoom != paths[pi][len(paths[pi])-1] {
					occupied[destRoom] = stateIdxToAntId[si]
				}
				states[si].posIdx = nextIdx
				turnMoves = append(turnMoves, [2]int{stateIdxToAntId[si], nextIdx})
				if nextIdx == len(paths[pi])-1 {
					finished++
				}
			}
		}
		if len(turnMoves) == 0 {
			// Deadlock (should not happen if paths are simple). Break to avoid infinite loop.
			break
		}
		// Sort by ant id for stable output ordering within a turn
		sort.Slice(turnMoves, func(i, j int) bool { return turnMoves[i][0] < turnMoves[j][0] })
		moves = append(moves, turnMoves)
	}
	return moves
}

func printOutput(inp *InputData, paths [][]string, assign []int, moves [][][2]int) {
	// Print original input as required
	for _, l := range inp.RawLines {
		fmt.Println(l)
	}
	fmt.Println()
	// Print moves
	// We need mapping from (pathIdx, index) to room name for each ant.
	// But we only stored (antId, destIndexInPath). We also need to know which path each ant follows.
	// Reconstruct antId -> pathIdx via assign ordering. We also need same order as simulate built stateIdxToAntId.
	antToPath := map[int]int{}
	id := 1
	for pi, cnt := range assign {
		for i := 0; i < cnt; i++ {
			antToPath[id] = pi
			id++
		}
	}
	for _, turn := range moves {
		parts := make([]string, len(turn))
		for i, mv := range turn {
			antId := mv[0]
			destIdx := mv[1]
			pi := antToPath[antId]
			room := paths[pi][destIdx]
			parts[i] = fmt.Sprintf("L%d-%s", antId, room)
		}
		fmt.Println(strings.Join(parts, " "))
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}
	path := os.Args[1]
	inp, err := parseInput(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Print the number of ants and the rest of the file, matching examples: the file itself already contains ants? The spec shows first line is number_of_ants as part of the file.
	// Our parser expects first line ants in the file; RawLines includes it.

	// Build network and compute max flow
	net := buildNetwork(inp)
	flow := net.D.maxFlow(net.Start, net.End)
	if flow == 0 {
		// No path: per spec, should output error? The examples for invalid have ERROR. Here, treat as error.
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}
	// Extract paths from residual network
	paths := extractPaths(net, inp)
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}

	// Sort paths by length ascending for stable distribution
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })

	assign := distributeAnts(inp.Ants, paths)
	moves := simulate(inp.Ants, paths, assign)

	printOutput(inp, paths, assign, moves)
}

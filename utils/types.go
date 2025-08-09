package utils

type Room struct {
	Name    string
	X, Y    int
	Links   []*Room
	IsStart bool
	IsEnd   bool
}

type Graph struct {
	AntCount int
	Rooms    map[string]*Room
	Start    *Room
	End      *Room
}

type Ant struct {
	Id   int
	Path []*Room
	Pos  int
}

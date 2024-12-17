package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"
)

type pos struct {
	x int
	y int
}

type Input struct {
	field [][]string
	start pos
	end   pos
}

type CostKey struct {
	pos pos
	dir pos
}

type HeapValue struct {
	cost  float64
	path  []pos
	dir   pos
	index int
}

type NextValue struct {
	nextPos pos
	dir     pos
	cost    float64
}

type PriorityQueue []*HeapValue

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*HeapValue)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

var directions []pos = []pos{
	{1, 0},
	{0, -1},
	{-1, 0},
	{0, 1},
}

func main() {
	data := read()
	printField(&data.field, data.start, data.end, []pos{})
	findMinPath(&data.field, data.start, data.end)
}

func findMinPath(field *[][]string, start, end pos) {
	costs := map[CostKey]float64{}
	costs[CostKey{start, directions[0]}] = 0
	currPath := []pos{start}
	stack := PriorityQueue{
		{
			0,
			currPath,
			directions[0],
			0,
		},
	}
	heap.Init(&stack)
	min_path := []pos{}
	min_cost := math.Inf(1)
	for stack.Len() != 0 {

		newVal := heap.Pop(&stack).(*HeapValue)
		currPos := newVal.path[len(newVal.path)-1]
		if currPos == end {
			min_cost = math.Min(min_cost, newVal.cost)
			if newVal.cost <= min_cost {
				appendUniqPath(&min_path, &newVal.path)
			}
			continue
		}
		for _, next := range []NextValue{
			{
				pos{currPos.x + newVal.dir.x, currPos.y + newVal.dir.y},
				newVal.dir,
				newVal.cost + 1.0,
			},
			{
				currPos,
				pos{newVal.dir.y * -1, newVal.dir.x},
				newVal.cost + 1000.0,
			},
			{
				currPos,
				pos{newVal.dir.y, newVal.dir.x * -1},
				newVal.cost + 1000.0,
			},
		} {
			if (*field)[next.nextPos.y][next.nextPos.x] == "#" {
				continue
			}
			costKey := CostKey{next.nextPos, next.dir}
			if next.cost <= getValueOrDefault(costs, costKey) {
				costs[costKey] = next.cost
				newPath := []pos{}
				newPath = append(newPath, newVal.path...)
				newPath = append(newPath, next.nextPos)
				heap.Push(&stack, &HeapValue{
					next.cost,
					newPath,
					next.dir,
					stack.Len(),
				})
			}
		}
	}
	fmt.Println(min_cost)
	fmt.Println(len(min_path))
}

func appendUniqPath(initial *[]pos, new *[]pos) *[]pos {
	for _, el := range *new {
		if !slices.Contains(*initial, el) {
			*initial = append(*initial, el)
		}
	}
	return initial
}

func getValueOrDefault(dict map[CostKey]float64, key CostKey) float64 {
	if val, ok := dict[key]; ok {
		return val
	}
	return math.Inf(1)
}

func printField(field *[][]string, start pos, end pos, visited []pos) {
	for y, rows := range *field {
		for x, col := range rows {
			curr := pos{x, y}
			if slices.Contains(visited, curr) {
				fmt.Print("X")
			} else {
				fmt.Print(col)
			}
		}
		fmt.Println()
	}
}

func canMove(field *[][]string, curr, dir pos, visitted []pos) bool {
	nextX, nextY := curr.x+dir.x, curr.y+dir.y
	if nextX < 0 || nextX >= len((*field)[0]) || nextY < 0 || nextY >= len(*field) {
		return false
	}
	if slices.Contains(visitted, pos{nextX, nextY}) {
		return false
	}
	if ((*field)[nextY][nextX]) == "." {
		return true
	}
	return false
}

func read() Input {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	field := [][]string{}
	start := pos{}
	end := pos{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		if index := strings.Index(line, "S"); index != -1 {
			start = pos{index, len(field)}
			line = strings.Replace(line, "S", ".", 1)
		} else if index := strings.Index(line, "E"); index != -1 {
			end = pos{index, len(field)}
			line = strings.Replace(line, "E", ".", 1)
		}
		field = append(field, strings.Split(line, ""))
	}
	return Input{
		field,
		start,
		end,
	}
}

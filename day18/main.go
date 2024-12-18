package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type pos struct {
	x int
	y int
}

var directions []pos = []pos{
	{1, 0},
	{0, -1},
	{-1, 0},
	{0, 1},
}

type Input struct {
	field []pos
}
type CostKey struct {
	pos pos
}

type HeapValue struct {
	cost  float64
	path  []pos
	index int
}

var fieldSize = 71

func main() {
	data := read()
	fmt.Println(findMinPath(data.field[:2873], pos{0, 0}, pos{fieldSize - 1, fieldSize - 1}))
	fmt.Println(findCutOff(data.field, pos{0, 0}, pos{fieldSize - 1, fieldSize - 1}))
}

func findCutOff(allBytes []pos, start, end pos) pos {
	wg := sync.WaitGroup{}
	results := make([]int, len(allBytes))
	mu := sync.Mutex{}
	inner := func(bytes []pos) {
		defer wg.Done()
		minPath := findMinPath(bytes, start, end)
		mu.Lock()
		results[len(bytes)] = minPath
		mu.Unlock()
	}
	for i := range len(allBytes) - 1 {
		wg.Add(1)
		go inner(allBytes[:i])
	}
	wg.Wait()
	for i, el := range results {
		if el == math.MaxInt {
			return allBytes[i-1]
		}
	}
	return end
}

func isValid(p pos) bool {
	if p.x < 0 || p.x >= fieldSize || p.y < 0 || p.y >= fieldSize {
		return false
	}
	return true
}

func findMinPath(field []pos, start, end pos) int {
	cost := map[pos]int{}
	for y := range fieldSize {
		for x := range fieldSize {
			cost[pos{x, y}] = math.MaxInt
		}
	}
	cost[start] = 0
	q := make([]pos, 0, 50)
	q = append(q, start)
	for len(q) > 0 {
		current := q[0]
		q = q[1:]
		for _, dir := range directions {
			next := pos{current.x + dir.x, current.y + dir.y}
			if !isValid(next) {
				continue
			}
			if slices.Contains(field, next) {
				continue
			}
			if cost[next] <= cost[current]+1 {
				continue
			}

			cost[next] = cost[current] + 1
			q = append(q, next)
		}
	}
	return cost[end]
}

func read() Input {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	field := []pos{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		vals := strings.Split(line, ",")
		x, err := strconv.ParseInt(vals[0], 10, 0)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(vals[1], 10, 0)
		if err != nil {
			panic(err)
		}
		field = append(field, pos{int(x), int(y)})

	}
	return Input{
		field,
	}
}

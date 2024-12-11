package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

var directions []Point = []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func main() {
	data := read()
	fmt.Println(calculateScores(data, true))
}

func dfs(data [][]int, x, y int, direction Point, seen map[Point]bool, all bool) int {
	h, w := len(data), len(data[0])
	score := 0
	newX, newY := x+direction.x, y+direction.y
	if newX < 0 || newX >= w || newY < 0 || newY >= h {
		return score
	}
	if data[newY][newX]-data[y][x] == 1 {
		if data[newY][newX] == 9 {
			if _, ok := seen[Point{newX, newY}]; ok {
				if all {
					return score + 1
				}
				return score
			} else {
				seen[Point{newX, newY}] = true
				return score + 1
			}
		}
		for _, direction := range directions {
			score = score + dfs(data, newX, newY, direction, seen, all)
		}
	}
	return score
}

func calculatePathScore(data [][]int, x, y int, all bool) int {
	score := 0
	seen := map[Point]bool{}
	for _, direction := range directions {
		subScore := dfs(data, x, y, direction, seen, all)
		fmt.Println(subScore, x, y, direction)
		score = score + subScore
	}
	return score
}

func calculateScores(data [][]int, all bool) int {
	score := 0
	for y, row := range data {
		for x, el := range row {
			if el == 0 {
				startingPointScore := calculatePathScore(data, x, y, all)

				score = score + startingPointScore
			}
		}
	}
	return score
}

func read() [][]int {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data := [][]int{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		dataLine := []int{}
		for _, el := range strings.Split(line, "") {
			parsed, err := strconv.ParseInt(el, 10, 0)
			if err != nil {
				panic(err)
			}
			dataLine = append(dataLine, int(parsed))
		}
		data = append(data, dataLine)
	}
	return data
}

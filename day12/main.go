package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type direction struct{ x, y int }

var directions []direction = []direction{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

func main() {
	data := read()
	fmt.Println(part1(&data))
}

func part1(data *[][]string) int {
	seen := map[direction]bool{}
	groups := [][]direction{}
	for r, rows := range *data {
		for c := range rows {
			if _, ok := seen[direction{c, r}]; ok {
				continue
			}
			group := []direction{{c, r}}
			seen[direction{c, r}] = true
			dfs(data, &seen, c, r, &group)
			groups = append(groups, group)
		}
	}
	checksum := 0
	for _, group := range groups {
		area := len(group)
		sides := calculateRegionSides(group)
		fmt.Println(area, sides, area*sides)
		checksum += area * sides
	}
	return checksum
}

func calculateFencePerimeter(data *[][]string, pos direction) int {
	fenceCount := 0
	rows := len((*data))
	cols := len((*data)[0])
	el := (*data)[pos.y][pos.x]
	for _, dir := range directions {
		dx, dy := pos.x+dir.x, pos.y+dir.y
		if dx >= 0 && dx < cols && dy >= 0 && dy < rows {
			if el != (*data)[dy][dx] {
				fenceCount++
			}
		} else {
			fenceCount++
		}
	}
	return fenceCount
}

func calculateRegionSides(region []direction) int {
	sides := 0
	for _, dir := range directions {
		seen := map[direction]bool{}
		for _, point := range region {
			if seen[point] {
				continue
			}
			x, y := point.x, point.y
			neighbor := direction{x + dir.x, y + dir.y}

			if slices.Contains(region, neighbor) {
				continue
			}
			sides++

			for delta := -1; delta <= 1; delta += 2 {
				nx, ny := x, y
				for {
					next := direction{nx + dir.x, ny + dir.y}
					if !slices.Contains(region, direction{nx, ny}) || slices.Contains(region, next) {
						break
					}
					seen[direction{nx, ny}] = true
					nx += dir.y * delta
					ny += dir.x * delta

				}
			}
		}
	}
	return sides
}

func contains(region []direction, point direction) bool {
	for _, p := range region {
		if p == point {
			return true
		}
	}
	return false
}

func dfs(data *[][]string, seen *map[direction]bool, x, y int, group *[]direction) {
	rows := len(*data)
	cols := len((*data)[0])
	el := (*data)[y][x]
	for _, dir := range directions {
		dx, dy := x+dir.x, y+dir.y
		if dx >= 0 && dx < cols && dy >= 0 && dy < rows && el == (*data)[dy][dx] {
			if _, ok := (*seen)[direction{dx, dy}]; ok {
				continue
			}
			(*seen)[direction{dx, dy}] = true
			*group = append(*group, direction{dx, dy})
			dfs(data, seen, dx, dy, group)
		}
	}
}

func read() [][]string {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data := [][]string{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		dataLine := []string{}
		for _, el := range strings.Split(line, "") {
			dataLine = append(dataLine, el)
		}
		data = append(data, dataLine)
	}
	return data
}

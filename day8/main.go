package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Point struct {
	X int
	Y int
}

type InputData struct {
	Stations   map[string][]Point
	GridWidth  int
	GridHeight int
}

func main() {
	data := read()
	fmt.Println(data)
	fmt.Println(part1(data.Stations, data.GridWidth, data.GridHeight))
	fmt.Println(part2(data.Stations, data.GridWidth, data.GridHeight))
}

func part1(stations map[string][]Point, x int, y int) int {
	antinodes := map[Point]bool{}
	for _, positions := range stations {
		for index, station1 := range positions[0 : len(positions)-1] {
			for _, station2 := range positions[index+1:] {
				var antinode1Position Point
				var antinode2Position Point
				if station2.X >= station1.X {
					antinode1Position = Point{station1.X - (station2.X - station1.X), station1.Y - (station2.Y - station1.Y)}
					antinode2Position = Point{station2.X + (station2.X - station1.X), station2.Y + (station2.Y - station1.Y)}
				} else {
					antinode1Position = Point{station1.X + (station1.X - station2.X), station1.Y - (station2.Y - station1.Y)}
					antinode2Position = Point{station2.X - (station1.X - station2.X), station2.Y + (station2.Y - station1.Y)}
				}
				if antinode1Position.X >= 0 && antinode1Position.X < x && antinode1Position.Y >= 0 && antinode1Position.Y < y {
					antinodes[antinode1Position] = true
				}
				if antinode2Position.X >= 0 && antinode2Position.X < x && antinode2Position.Y >= 0 && antinode2Position.Y < y {
					antinodes[antinode2Position] = true
				}
			}
		}
	}
	return len(antinodes)
}

func part2(stations map[string][]Point, x int, y int) int {
	antinodes := map[Point]bool{}
	for _, positions := range stations {
		for index, station1 := range positions[0 : len(positions)-1] {
			for _, station2 := range positions[index+1:] {
				antinodes[station1] = true
				antinodes[station2] = true
				var xShift int
				var yShift int
				xShiftDirection := 1
				var antinode1Position Point
				var antinode2Position Point
				if station2.X >= station1.X {
					xShift = station2.X - station1.X
					yShift = station2.Y - station1.Y
					xShiftDirection = -1
				} else {
					xShift = station1.X - station2.X
					yShift = station2.Y - station1.Y
				}
				for modifier := range x {
					if modifier == 0 {
						continue
					}
					antinode1Position = Point{station1.X + xShift*modifier*xShiftDirection, station1.Y - yShift*modifier}
					antinode2Position = Point{station2.X - xShift*modifier*xShiftDirection, station2.Y + yShift*modifier}
					if antinode1Position.X >= 0 && antinode1Position.X < x && antinode1Position.Y >= 0 && antinode1Position.Y < y {
						antinodes[antinode1Position] = true
					}
					if antinode2Position.X >= 0 && antinode2Position.X < x && antinode2Position.Y >= 0 && antinode2Position.Y < y {
						antinodes[antinode2Position] = true
					}
				}
			}
		}
	}
	print(antinodes, x, y)
	return len(antinodes)
}

func print(antinodes map[Point]bool, x, y int) {
	for i := range x {
		for j := range y {
			_, ok := antinodes[Point{j, i}]
			if ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func read() InputData {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data := InputData{Stations: map[string][]Point{}}
	y := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		for x, el := range strings.Split(strings.TrimSpace(line), "") {
			data.GridWidth = x + 1
			if el == "." {
				continue
			}
			fmt.Println(el)
			data.Stations[el] = append(data.Stations[el], Point{x, y})
		}
		data.GridHeight = y + 1
		y = y + 1
	}
	return data
}

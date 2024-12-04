package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type InputData struct {
	Levels [][]string
}

type Position struct {
	H int
	W int
}

type Direction struct {
	V int
	H int
}

func main() {
	data, err := read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	part1(data)
}

func part1(data InputData) {
	h := len(data.Levels)
	w := len(data.Levels[0])
	startIndexes := []Position{}
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if data.Levels[i][j] == "X" {
				startIndexes = append(startIndexes, Position{i, j})
			}
		}
	}
	fmt.Println(startIndexes)
	totalCount := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if j == 0 && i == 0 {
				continue
			}
			for _, position := range startIndexes {
				if trackWord(&data, position, Direction{i, j}, strings.Split("MAS", "")) {
					totalCount++
				}
			}
		}
	}
	fmt.Println(totalCount)
	startIndexes = []Position{}
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if data.Levels[i][j] == "A" {
				startIndexes = append(startIndexes, Position{i, j})
			}
		}
	}
	fmt.Println(startIndexes)
	totalCount = 0
	for _, position := range startIndexes {
		if part2(&data, position) {
			totalCount++
		}
	}
	fmt.Println("Part2 count:", totalCount)
}

func part2(data *InputData, position Position) bool {
	fmt.Println("position:", position)
	levels := data.Levels
	h := len(levels)
	w := len(levels[0])
	if position.W+1 >= w || position.H+1 >= h || position.W-1 < 0 || position.H-1 < 0 {
		return false
	}
	if ((levels[position.H-1][position.W-1] == "S" && levels[position.H+1][position.W+1] == "M") || (levels[position.H-1][position.W-1] == "M" && levels[position.H+1][position.W+1] == "S")) &&
		((levels[position.H+1][position.W-1] == "S" && levels[position.H-1][position.W+1] == "M") || (levels[position.H+1][position.W-1] == "M" && levels[position.H-1][position.W+1] == "S")) {
		return true
	}
	return false
}

func trackWord(data *InputData, position Position, direction Direction, leftLetters []string) bool {
	fmt.Println("leftLetters:", leftLetters, "position:", position, "direction:", direction, "prevLetter:", data.Levels[position.H][position.W])
	h := len(data.Levels)
	w := len(data.Levels[0])
	newPosition := Position{
		position.H + direction.V,
		position.W + direction.H,
	}
	if 0 <= newPosition.H && newPosition.H < h && 0 <= newPosition.W && newPosition.W < w {
		fmt.Println("Current letter:", data.Levels[newPosition.H][newPosition.W])
		if data.Levels[newPosition.H][newPosition.W] == leftLetters[0] {
			if len(leftLetters) == 1 {
				return true
			} else {
				return trackWord(data, newPosition, direction, leftLetters[1:])
			}
		}
	}
	return false
}

func read() (InputData, error) {
	f, err := os.Open("task.txt")
	if err != nil {
		fmt.Println(err)
		return InputData{}, nil
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	data := InputData{}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return data, err
		}
		level := parseLevel(line)
		data.Levels = append(data.Levels, level)
	}
	return data, nil
}

func parseLevel(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	return strings.Split(trimmed, "")
}

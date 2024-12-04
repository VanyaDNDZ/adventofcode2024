package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type InputData struct {
	Levels [][]int64
}

func main() {
	data, err := read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	safeLevels := 0
	for _, level := range data.Levels {
		fmt.Println("Checing level:", level)
		if isSafeLevel(level) {
			safeLevels++
		} else {
			for i := 0; i < len(level); i++ {
				levelCopy := make([]int64, len(level))
				copy(levelCopy, level)
				levelCopy = slices.Delete(levelCopy, i, i+1)
				fmt.Println("Copy:", levelCopy)
				if isSafeLevel(levelCopy) {
					safeLevels++
					break
				}

			}
		}
	}
	fmt.Println("SafelLevels count:", safeLevels)
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
		level, err := parseLevel(line)
		if err != nil {
			return InputData{}, err
		}
		data.Levels = append(data.Levels, level)
	}
	return data, nil
}

func parseLevel(raw string) ([]int64, error) {
	trimmed := strings.TrimSpace(raw)
	data := []int64{}
	for _, split := range strings.Split(trimmed, " ") {
		val, err := strconv.ParseInt(split, 10, 0)
		if err != nil {
			return data, err
		}
		data = append(data, val)
	}
	return data, nil
}

func isSafeLevel(level []int64) bool {
	if len(level) < 2 {
		return false
	}
	var isIncreasing bool
	if level[0] > level[1] {
		isIncreasing = false
	} else {
		isIncreasing = true
	}

	for i := 1; i < len(level); i++ {
		if level[i-1] == level[i] {
			return false
		} else if level[i-1] > level[i] && isIncreasing {
			return false
		} else if level[i-1] < level[i] && !isIncreasing {
			return false
		}
		if isIncreasing && level[i]-level[i-1] > 3 {
			return false
		} else if !isIncreasing && level[i-1]-level[i] > 3 {
			return false
		}
	}
	return true
}

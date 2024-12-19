package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Input struct {
	patterns map[string]bool
	designs  []string
}

func main() {
	data := read()
	possible := 0
	totalCount := 0
	for _, design := range data.designs {
		possibleVariants := variants(design, data.patterns)
		if possibleVariants > 0 {
			possible++
		}
		totalCount += possibleVariants
	}
	fmt.Println(possible)
	fmt.Println(totalCount)
}

func variants(design string, patterns map[string]bool) int {
	memo := make(map[string]int)
	var inner func(subDesign string) int
	inner = func(subDesign string) int {
		if count, ok := memo[subDesign]; ok {
			return count
		}
		if len(subDesign) == 0 {
			return 1
		}
		count := 0
		for pattern := range patterns {
			if left, ok := strings.CutPrefix(subDesign, pattern); ok {
				count += inner(left)
			}
		}
		memo[subDesign] = count
		return count
	}
	return inner(design)
}

func read() Input {
	file, err := os.Open("/Users/ivansalamakha/work/repos/golang/adventofcode/day19/task.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	patterns := map[string]bool{}
	designs := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, ",") {
			for _, pattern := range strings.Split(line, ",") {
				patterns[strings.TrimSpace(pattern)] = true
			}
		} else {
			designs = append(designs, line)
		}
	}
	return Input{
		patterns,
		designs,
	}
}

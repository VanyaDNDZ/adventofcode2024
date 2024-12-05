package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type InputData struct {
	Rules map[int][]int
	Pages [][]int
}

type Rule struct {
	Head int
	Tail int
}

func main() {
	data, err := read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("InputData:", data)
	total := 0
	corretedTotal := 0
	for _, pages := range data.Pages {
		if isCorrectPage(data.Rules, pages) {
			total = total + pages[len(pages)/2]
		} else {
			corrected := reorder(data.Rules, pages)
			corretedTotal = corretedTotal + corrected[len(corrected)/2]
		}
	}
	fmt.Println("Sum of middles:", total)
	fmt.Println("Sum of corrected middles:", corretedTotal)
}

func isCorrectPage(rules map[int][]int, pages []int) bool {
	for index, page := range pages {
		val, _ := rules[page]
		for _, tail := range val {
			if slices.Contains(pages[:index], tail) {
				return false
			}
		}
	}
	return true
}

func reorder(rules map[int][]int, pages []int) []int {
	priorities := map[int]int{}
	for _, page := range pages {
		seenPages := 0
		val, ok := rules[page]
		if ok {
			for _, subPage := range pages {
				if slices.Contains(val, subPage) {
					seenPages++
				}
			}
		}
		priorities[page] = seenPages
	}
	sort.Slice(pages, func(i, j int) bool {
		p1, _ := priorities[pages[i]]
		p2, _ := priorities[pages[j]]
		return p1 > p2
	})
	return pages
}

func read() (InputData, error) {
	file, err := os.Open("task.txt")
	if err != nil {
		return InputData{}, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data := InputData{}
	data.Rules = make(map[int][]int)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return InputData{}, err
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, "|") {
			head_and_tail := strings.Split(line, "|")
			head, err := strconv.ParseInt(head_and_tail[0], 10, 0)
			if err != nil {
				panic(err)
			}
			tail, err := strconv.ParseInt(head_and_tail[1], 10, 0)
			if err != nil {
				panic(err)
			}

			value, ok := data.Rules[int(head)]
			if ok {
				data.Rules[int(head)] = append(value, int(tail))
			} else {
				data.Rules[int(head)] = []int{int(tail)}
			}
		} else {
			pages := strings.Split(line, ",")
			parsed := make([]int, len(pages))
			for index, page := range pages {
				page, err := strconv.ParseInt(page, 10, 0)
				if err != nil {
					panic(err)
				}
				parsed[index] = int(page)
			}
			data.Pages = append(data.Pages, parsed)
		}
	}
	return data, nil
}

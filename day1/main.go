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
	Left  []int64
	Right []int64
}

func main() {
	data, err := read()
	if err != nil {
		fmt.Println(err)
	}
	slices.Sort(data.Left)
	slices.Sort(data.Right)
	fmt.Println("data:", data)
	var sum int64 = 0
	for i := 0; i < len(data.Left); i++ {
		var diff int64 = 0

		if data.Left[i] > data.Right[i] {
			diff = data.Left[i] - data.Right[i]
		} else {
			diff = data.Right[i] - data.Left[i]
		}
		fmt.Println("Left:", data.Left[i], "; Right:", data.Right[i], "; diff:", diff)
		sum = sum + diff
	}
	fmt.Println("Final diff is:", sum)
	fmt.Println("Part two")
	counts := map[int64]int{}
	for _, el := range data.Right {
		v, ok := counts[el]
		if ok {
			counts[el] = v + 1
		} else {
			counts[el] = 1
		}
	}
	var similarityScore int64 = 0
	for _, el := range data.Left {
		v, ok := counts[el]
		if ok {
			similarityScore = similarityScore + int64(v)*el
		}
	}

	fmt.Println("Similarity Score:", similarityScore)
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
		splits := strings.Split(line, "   ")
		if len(splits) != 2 {
			return data, err
		}
		left, err := strconv.ParseInt(splits[0], 10, 0)
		if err != nil {
			return data, err
		}
		data.Left = append(data.Left, left)
		right, err := strconv.ParseInt(splits[1][:len(splits[1])-1], 10, 0)
		if err != nil {
			return data, err
		}
		data.Right = append(data.Right, right)

	}
	return data, nil
}

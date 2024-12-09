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

type Op string

const (
	PLUS   Op = "+"
	MUL    Op = "*"
	CONCAT Op = "||"
)

type Eq struct {
	Expected int64
	Numbers  []int64
}

type InputData struct {
	Equations []Eq
}

func main() {
	data := read()
	fmt.Println(data)
	var sum int64 = 0
	for _, eq := range data.Equations {
		exp := isValidEq(eq.Expected, eq.Numbers, PLUS)
		if exp {
			sum = sum + eq.Expected
			continue
		}
		exp = isValidEq(eq.Expected, eq.Numbers, MUL)
		if exp {
			sum = sum + eq.Expected
		}
	}
	fmt.Println(sum)
}

func isValidEq(expected int64, unusedNumbers []int64, op Op) bool {
	stack := []int64{unusedNumbers[0]}
	for _, el := range unusedNumbers[1:] {
		nextStack := []int64{}
		for {
			if len(stack) == 0 {
				break
			}
			x := stack[0]
			stack = stack[1:]
			for _, op := range []Op{PLUS, MUL, CONCAT} {
				switch op {
				case PLUS:
					nextStack = append(nextStack, x+el)
				case MUL:
					nextStack = append(nextStack, x*el)
				case CONCAT:
					{
						val, err := strconv.ParseInt(fmt.Sprintf("%d%d", x, el), 10, 0)
						if err != nil {
							panic(err)
						}
						nextStack = append(nextStack, val)
					}
				}
			}
		}
		stack = nextStack
	}
	return slices.Contains(stack, expected)
}

func read() InputData {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data := InputData{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		head_and_tail := strings.Split(line, ":")
		expected, err := strconv.ParseInt(head_and_tail[0], 10, 0)
		if err != nil {
			panic(err)
		}
		eq := Eq{Expected: expected, Numbers: []int64{}}
		for _, el := range strings.Split(strings.TrimSpace(head_and_tail[1]), " ") {
			parsed, err := strconv.ParseInt(el, 10, 0)
			if err != nil {
				panic(err)
			}
			eq.Numbers = append(eq.Numbers, parsed)
		}
		data.Equations = append(data.Equations, eq)
	}
	return data
}

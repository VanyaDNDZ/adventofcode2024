package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	eq       = `^Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X\=(\d+), Y\=(\d+)`
	MaxTries = 100
	offset   = 10000000000000
)

type Matrix struct {
	coefMatrix  [][]int
	constMatrix []int
}

func main() {
	eqs := parse()
	sum := 0
	for _, eq := range eqs {
		sum += part1(eq)
	}
	fmt.Println("Part1:", sum)
	sum = 0
	for _, eq := range eqs {
		sum += part2(eq)
	}
	fmt.Println("Part1:", sum)
}

func part1(matrix Matrix) int {
	// c1*b2 - c2*b1
	det_x := matrix.constMatrix[0]*matrix.coefMatrix[1][1] - matrix.constMatrix[1]*matrix.coefMatrix[0][1]
	// a1*c2 - a2*c1
	det_y := matrix.constMatrix[1]*matrix.coefMatrix[0][0] - matrix.constMatrix[0]*matrix.coefMatrix[1][0]
	// a1*b2 - a2*b1
	det_matrix := matrix.coefMatrix[0][0]*matrix.coefMatrix[1][1] - matrix.coefMatrix[1][0]*matrix.coefMatrix[0][1]
	x, y := det_x/det_matrix, det_y/det_matrix
	fmt.Println(x, y, det_x%det_matrix, det_y%det_matrix, det_x, det_y, det_matrix, matrix)
	if x <= MaxTries && y <= MaxTries && det_x%det_matrix == 0 && det_y%det_matrix == 0 {
		return x*3 + y
	}
	return 0
}

func part2(matrix Matrix) int {
	// c1*b2 - c2*b1
	det_x := (matrix.constMatrix[0]+offset)*matrix.coefMatrix[1][1] - (matrix.constMatrix[1]+offset)*matrix.coefMatrix[0][1]
	// a1*c2 - a2*c1
	det_y := (matrix.constMatrix[1]+offset)*matrix.coefMatrix[0][0] - (matrix.constMatrix[0]+offset)*matrix.coefMatrix[1][0]
	// a1*b2 - a2*b1
	det_matrix := matrix.coefMatrix[0][0]*matrix.coefMatrix[1][1] - matrix.coefMatrix[1][0]*matrix.coefMatrix[0][1]
	x, y := det_x/det_matrix, det_y/det_matrix
	fmt.Println(x, y, det_x%det_matrix, det_y%det_matrix, det_x, det_y, det_matrix, matrix)
	if det_x%det_matrix == 0 && det_y%det_matrix == 0 {
		return x*3 + y
	}
	return 0
}

func parse() []Matrix {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	exp := regexp.MustCompile(eq)
	defer file.Close()
	reader := bufio.NewReader(file)
	data := []Matrix{}
	fullEq := []string{}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if line != "\n" {
			fullEq = append(fullEq, line)
		}
		if len(fullEq) == 3 {
			matches := exp.FindAllStringSubmatch(strings.Join(fullEq, ""), -1)
			matchesInt := make([]int, 6)
			for i, el := range matches[0][1:] {
				parsed, err := strconv.ParseInt(el, 10, 0)
				if err != nil {
					panic(err)
				}
				matchesInt[i] = int(parsed)
			}
			matrix := Matrix{}
			matrix.coefMatrix = append(matrix.coefMatrix, []int{matchesInt[0], matchesInt[2]})
			matrix.coefMatrix = append(matrix.coefMatrix, []int{matchesInt[1], matchesInt[3]})
			matrix.constMatrix = []int{matchesInt[4], matchesInt[5]}
			fullEq = []string{}
			data = append(data, matrix)
		}

	}
	return data
}

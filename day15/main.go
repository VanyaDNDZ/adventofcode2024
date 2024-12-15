package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type pos struct{ x, y int }

type Input struct {
	robot    pos
	insructs []pos
	field    [][]string
}

type moveBlock struct {
	position pos
	char     string
}

func main() {
	data := read(false)
	processedField := process(&data.robot, data.insructs, &data.field)
	sum := 0
	for i, rows := range *processedField {
		for j, col := range rows {
			if col == "O" {
				sum += i*100 + j
			}
		}
	}
	fmt.Println(sum)
	data = read(true)
	printField(&data.field)
	processedField = part2(&data.robot, data.insructs, &data.field)
	printField(processedField)
	sum = 0
	for i, rows := range *processedField {
		for j, col := range rows {
			if col == "[" {
				sum += i*100 + j
			}
		}
	}
	fmt.Println(sum)
}

func process(robot *pos, instructs []pos, field *[][]string) *[][]string {
	for _, dir := range instructs {
		if (*field)[robot.y+dir.y][robot.x+dir.x] == "O" {
			gap := nextGap(robot, &dir, field)
			if gap.x == -1 {
				continue
			}
			(*field)[robot.y][robot.x] = "."
			(*field)[robot.y+dir.y][robot.x+dir.x] = "@"
			(*field)[gap.y][gap.x] = "O"
			robot.x, robot.y = robot.x+dir.x, robot.y+dir.y

		} else if (*field)[robot.y+dir.y][robot.x+dir.x] == "." {
			(*field)[robot.y][robot.x] = "."
			(*field)[robot.y+dir.y][robot.x+dir.x] = "@"
			robot.x, robot.y = robot.x+dir.x, robot.y+dir.y
		}
	}

	return field
}

func part2(robot *pos, instructs []pos, field *[][]string) *[][]string {
	for _, dir := range instructs {
		if (*field)[robot.y+dir.y][robot.x+dir.x] == "[" || (*field)[robot.y+dir.y][robot.x+dir.x] == "]" {
			moveBlocks := findBulkMove(robot, dir, field)
			if len(moveBlocks) == 0 {
				continue
			}
			for _, el := range moveBlocks {

				(*field)[el.position.y+dir.y][el.position.x+dir.x] = el.char
				(*field)[el.position.y][el.position.x] = "."

			}
			(*field)[robot.y][robot.x] = "."
			(*field)[robot.y+dir.y][robot.x+dir.x] = "@"
			robot.x, robot.y = robot.x+dir.x, robot.y+dir.y

		} else if (*field)[robot.y+dir.y][robot.x+dir.x] == "." {
			(*field)[robot.y][robot.x] = "."
			(*field)[robot.y+dir.y][robot.x+dir.x] = "@"
			robot.x, robot.y = robot.x+dir.x, robot.y+dir.y
		}
	}
	return field
}

func findBulkMove(robot *pos, dir pos, field *[][]string) []moveBlock {
	blocks := []moveBlock{}
	if dir.x != 0 {
		if dir.x == -1 {
			for i := robot.x; i > 0; i-- {
				if (*field)[robot.y][i] == "." {
					newBlocks := []moveBlock{}
					for j := len(blocks) - 1; j >= 0; j-- {
						newBlocks = append(newBlocks, blocks[j])
					}
					return newBlocks
				} else if (*field)[robot.y][i] == "[" || (*field)[robot.y][i] == "]" {
					blocks = append(blocks, moveBlock{pos{i, robot.y}, (*field)[robot.y][i]})
				} else if (*field)[robot.y][i] == "#" {
					return []moveBlock{}
				}
			}
		} else {
			for i := robot.x; i < len((*field)[0]); i++ {
				if (*field)[robot.y][i] == "." {
					newBlocks := []moveBlock{}
					for j := len(blocks) - 1; j >= 0; j-- {
						newBlocks = append(newBlocks, blocks[j])
					}
					return newBlocks

				} else if (*field)[robot.y][i] == "[" || (*field)[robot.y][i] == "]" {
					blocks = append(blocks, moveBlock{pos{i, robot.y}, (*field)[robot.y][i]})
				} else if (*field)[robot.y][i] == "#" {
					return []moveBlock{}
				}
			}
		}
	}
	stack := []moveBlock{}
	switch (*field)[robot.y+dir.y][robot.x+dir.x] {
	case "[":
		stack = append(stack, moveBlock{pos{robot.x + dir.x, robot.y + dir.y}, (*field)[robot.y+dir.y][robot.x+dir.x]}, moveBlock{pos{robot.x + dir.x + 1, robot.y + dir.y}, (*field)[robot.y+dir.y][robot.x+dir.x+1]})
	case "]":
		stack = append(stack, moveBlock{pos{robot.x + dir.x - 1, robot.y + dir.y}, (*field)[robot.y+dir.y][robot.x+dir.x-1]}, moveBlock{pos{robot.x + dir.x, robot.y + dir.y}, (*field)[robot.y+dir.y][robot.x+dir.x]})
	}
	blocks = append(blocks, stack...)
	for {
		if dir.y != 0 {
			nonEmpty := []moveBlock{}
			for _, el := range stack {
				switch (*field)[el.position.y+dir.y][el.position.x] {
				case "[":
					nonEmpty = append(nonEmpty, moveBlock{pos{el.position.x, el.position.y + dir.y}, (*field)[el.position.y+dir.y][el.position.x]}, moveBlock{pos{el.position.x + 1, el.position.y + dir.y}, (*field)[el.position.y+dir.y][el.position.x+1]})

				case "]":
					nonEmpty = append(nonEmpty, moveBlock{pos{el.position.x - 1, el.position.y + dir.y}, (*field)[el.position.y+dir.y][el.position.x-1]}, moveBlock{pos{el.position.x, el.position.y + dir.y}, (*field)[el.position.y+dir.y][el.position.x]})

				case "#":
					return []moveBlock{}
				}
			}
			if len(nonEmpty) != 0 {
				stack = nonEmpty
				blocks = append(blocks, nonEmpty...)
			} else {
				newBlocks := []moveBlock{}
				for j := len(blocks) - 1; j >= 0; j-- {
					newBlocks = append(newBlocks, blocks[j])
				}
				return newBlocks
			}
		}
	}
}

func printField(field *[][]string) {
	for _, rows := range *field {
		for _, col := range rows {
			fmt.Print(col)
		}
		fmt.Println()
	}
}

func nextGap(robot *pos, dir *pos, field *[][]string) pos {
	gap := pos{-1, -1}
	if dir.x == 0 {
		if dir.y == -1 {
			for i := robot.y; i > 0; i-- {
				if (*field)[i][robot.x] == "." {
					return pos{robot.x, i}
				} else if (*field)[i][robot.x] == "#" {
					return gap
				}
			}
		} else {
			for i := robot.y; i < len((*field)); i++ {
				if (*field)[i][robot.x] == "." {
					return pos{robot.x, i}
				} else if (*field)[i][robot.x] == "#" {
					return gap
				}
			}
		}
	} else {
		if dir.x == -1 {
			for i := robot.x; i > 0; i-- {
				if (*field)[robot.y][i] == "." {
					return pos{i, robot.y}
				} else if (*field)[robot.y][i] == "#" {
					return gap
				}
			}
		} else {
			for i := robot.x; i < len((*field)[0]); i++ {
				if (*field)[robot.y][i] == "." {
					return pos{i, robot.y}
				} else if (*field)[robot.y][i] == "#" {
					return gap
				}
			}
		}
	}
	return gap
}

func read(double bool) Input {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	field := [][]string{}
	robot := pos{}
	instructs := []pos{}
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
		if strings.Contains(line, "#") {
			if index := strings.Index(line, "@"); index != -1 {
				robot = pos{index, len(field)}
				if double {
					robot = pos{(index) * 2, len(field)}
				}
			}
			if double {
				newLine := []string{}
				for _, el := range strings.Split(line, "") {
					if el == "." {
						newLine = append(newLine, ".", ".")
					} else if el == "#" {
						newLine = append(newLine, "#", "#")
					} else if el == "O" {
						newLine = append(newLine, "[", "]")
					} else {
						newLine = append(newLine, el, ".")
					}
				}
				field = append(field, newLine)
			} else {
				field = append(field, strings.Split(line, ""))
			}
		} else {
			for _, move := range strings.Split(line, "") {
				dir := pos{}
				switch move {
				case "^":
					dir = pos{0, -1}
				case ">":
					dir = pos{1, 0}
				case "<":
					dir = pos{-1, 0}
				case "v":
					dir = pos{0, 1}

				}
				instructs = append(instructs, dir)
			}
		}
	}
	return Input{
		robot,
		instructs,
		field,
	}
}

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

type (
	pos      struct{ x, y int }
	velocity struct{ dx, dy int }
)

type robot struct {
	p pos
	v velocity
}

const expText = `p\=(-{0,1}\d+),(-{0,1}\d+) v\=(-{0,1}\d+),(-{0,1}\d+)`

func main() {
	w, h := 101, 103
	// robots := read()
	// printField(robots, h, w)
	// moveRobots(robots, h, w, 100)
	// fmt.Println(countSafetyFactor(robots, h, w))
	// part 2
	robots := read()
	for sec := range 100000 {
		moveRobots(robots, h, w, 1)
		if hasManyRobotsInLine(robots, h, w) {
			fmt.Println(sec)
			printField(robots, h, w)
		}

	}
}

func hasManyRobotsInLine(robots []robot, h, w int) bool {
	field := [][]string{}
	for range h {
		field = append(field, make([]string, w))
	}
	for _, robot := range robots {
		field[robot.p.y][robot.p.x] = "*"
	}
	for i, rows := range field {
		for j, col := range rows {
			if col == "" {
				field[i][j] = "."
			} else {
				field[i][j] = "*"
			}
		}
		if strings.Contains(strings.Join(field[i], ""), "*********") {
			return true
		}
	}
	return false
}

func moveRobots(robots []robot, h, w, sec int) []robot {
	for i, robot := range robots {
		robot.p = pos{(robot.p.x + robot.v.dx*sec) % w, (robot.p.y + robot.v.dy*sec) % h}
		if robot.p.x < 0 {
			robot.p.x = w + robot.p.x
		}
		if robot.p.y < 0 {
			robot.p.y = h + robot.p.y
		}
		robots[i] = robot
	}
	return robots
}

func countSafetyFactor(robots []robot, h, w int) int {
	q1, q2, q3, q4 := 0, 0, 0, 0
	for _, robot := range robots {
		if robot.p.x < w/2 && robot.p.y < h/2 {
			q1++
		} else if robot.p.x > w/2 && robot.p.y < h/2 {
			q2++
		} else if robot.p.x < w/2 && robot.p.y > h/2 {
			q3++
		} else if robot.p.x > w/2 && robot.p.y > h/2 {
			q4++
		}
	}
	fmt.Println(q1, q2, q3, q4)
	return q1 * q2 * q3 * q4
}

func printField(robots []robot, h, w int) {
	field := [][]int{}
	for range h {
		field = append(field, make([]int, w))
	}
	for _, robot := range robots {
		field[robot.p.y][robot.p.x]++
	}
	for _, rows := range field {
		for _, col := range rows {
			if col == 0 {
				fmt.Print(".")
			} else {
				fmt.Print(col)
			}
		}
		fmt.Println()
	}
}

func read() []robot {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	exp := regexp.MustCompile(expText)
	defer file.Close()
	reader := bufio.NewReader(file)
	robots := []robot{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		matches := exp.FindStringSubmatch(line)
		parsed := []int{}
		for _, el := range matches[1:] {
			val, err := strconv.ParseInt(el, 10, 0)
			if err != nil {
				panic(err)
			}
			parsed = append(parsed, int(val))
		}
		robots = append(robots, robot{
			pos{parsed[0], parsed[1]},
			velocity{parsed[2], parsed[3]},
		})
	}
	return robots
}

package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
)

type Position struct {
	X int
	Y int
}

type Direction struct {
	X int
	Y int
}

type CellType string

const (
	Empty    CellType = "empty"
	Obstacle CellType = "obstacle"
	Visited  CellType = "visited"
)

var (
	Up    Direction = Direction{0, -1}
	Down  Direction = Direction{0, 1}
	Left  Direction = Direction{-1, 0}
	Right Direction = Direction{1, 0}
)

type Cell struct {
	CellType         CellType
	Position         Position
	VisitedDirection []Direction
}

type Guard struct {
	Position  Position
	Direction Direction
}

type Map struct {
	Grid  [][]Cell
	Guard Guard
}

func main() {
	questMap := read()
	var copiedMap Map
	if err := DeepCopy(&questMap, &copiedMap); err != nil {
		fmt.Println(err)
		return
	}
	visited, err := walk(&questMap)
	if err != nil {
		fmt.Println("has loop")
	}
	fmt.Println("Visited cells:", visited)
	newObstacle := 0
	var wg sync.WaitGroup
	for i := 0; i < len(questMap.Grid); i++ {
		for j := 0; j < len(questMap.Grid[0]); j++ {
			cell := questMap.Grid[i][j]
			if cell.CellType == Visited {
				wg.Add(1)
				go func() {
					var innerMap Map
					if err := DeepCopy(&copiedMap, &innerMap); err != nil {
						fmt.Println(err)
						return
					}
					innerMap.Grid[i][j].CellType = Obstacle
					_, err := walk(&innerMap)
					if err != nil {
						newObstacle++
					}
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
	fmt.Println(newObstacle)
}

func walk(questMap *Map) (int, error) {
	visited := 1
	guard := questMap.Guard
	height := len(questMap.Grid)
	width := len(questMap.Grid[0])
	for {
		nextCellPosition := Position{
			guard.Position.X + guard.Direction.X,
			guard.Position.Y + guard.Direction.Y,
		}
		if nextCellPosition.X < 0 || nextCellPosition.X >= width || nextCellPosition.Y < 0 || nextCellPosition.Y >= height {
			break
		}
		nextCell := &questMap.Grid[nextCellPosition.Y][nextCellPosition.X]
		switch nextCell.CellType {
		case Empty:
			guard.Position = nextCellPosition
			visited++
			nextCell.CellType = Visited
			nextCell.VisitedDirection = append(nextCell.VisitedDirection, guard.Direction)
		case Visited:
			guard.Position = nextCellPosition
			if !slices.Contains(nextCell.VisitedDirection, guard.Direction) {
				nextCell.VisitedDirection = append(nextCell.VisitedDirection, guard.Direction)
			} else {
				return visited, errors.New("Loop")
			}
		case Obstacle:
			switch guard.Direction {
			case Up:
				guard.Direction = Right
			case Down:
				guard.Direction = Left
			case Right:
				guard.Direction = Down
			case Left:
				guard.Direction = Up
			}
		}
	}
	return visited, nil
}

func printMap(questMap *Map) {
	for i := 0; i < len(questMap.Grid); i++ {
		for j := 0; j < len(questMap.Grid[0]); j++ {
			switch questMap.Grid[i][j].CellType {
			case Empty:
				fmt.Print(".")
			case Visited:
				fmt.Print("X")
			case Obstacle:
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func read() Map {
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	grid := [][]Cell{}
	y := 0
	guard := Guard{Position{0, 0}, Up}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		mapLine := make([]Cell, len(line))
		for x, ch := range strings.Split(line, "") {
			switch ch {
			case ".":
				mapLine[x] = Cell{Empty, Position{x, y}, []Direction{}}
			case "#":
				mapLine[x] = Cell{Obstacle, Position{x, y}, []Direction{}}
			case "v":
				mapLine[x] = Cell{Visited, Position{x, y}, []Direction{}}
				guard.Position = Position{x, y}
				guard.Direction = Down
			case ">":
				mapLine[x] = Cell{Visited, Position{x, y}, []Direction{}}
				guard.Position = Position{x, y}
				guard.Direction = Right
			case "<":
				mapLine[x] = Cell{Visited, Position{x, y}, []Direction{}}
				guard.Position = Position{x, y}
				guard.Direction = Left
			case "^":
				mapLine[x] = Cell{Visited, Position{x, y}, []Direction{}}
				guard.Position = Position{x, y}
				guard.Direction = Up
			default:
				panic("Unxpexted value")
			}
		}
		y++
		grid = append(grid, mapLine)
	}
	return Map{grid, guard}
}

func DeepCopy(src, dest *Map) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	// Encode the source object
	if err := enc.Encode(src); err != nil {
		return err
	}

	// Decode into the destination object
	if err := dec.Decode(dest); err != nil {
		return err
	}

	return nil
}

package main

import (
	"fmt"
	"strconv"
	"strings"
)

var input []string = strings.Split("3 386358 86195 85 1267 3752457 0 741", " ")

var (
	cache       map[string][]string
	blinksCache map[int]map[string]int64
)

func main() {
	cache = map[string][]string{}
	blinksCache = map[int]map[string]int64{}
	fmt.Println("Total stones count is:", part1())
}

func mutate(data string, blinks int) int64 {
	var counter int64 = 0
	stones := []string{}
	if cachedBlinks, ok := blinksCache[blinks]; ok {
		if cachedBlinksStone, ok := cachedBlinks[data]; ok {
			return cachedBlinksStone
		}
	}
	if cached, ok := cache[data]; ok {
		stones = cached
	} else {
		if data == "0" {
			stones = append(stones, "1")
		} else if len(data)%2 == 0 {
			left, right := data[:len(data)/2], data[len(data)/2:]
			val, err := strconv.ParseInt(left, 10, 0)
			if err != nil {
				panic(err)
			}
			stones = append(stones, fmt.Sprint(val))
			val, err = strconv.ParseInt(right, 10, 0)
			if err != nil {
				panic(err)
			}
			stones = append(stones, fmt.Sprint(val))
		} else {
			val, err := strconv.ParseInt(data, 10, 0)
			if err != nil {
				panic(err)
			}
			stones = append(stones, fmt.Sprint(val*2024))
		}
		cache[data] = stones
	}
	if blinks > 1 {

		for _, stone := range stones {
			newStonesCount := mutate(stone, blinks-1)
			if cachedBlinks, ok := blinksCache[blinks-1]; ok {
				cachedBlinks[stone] = newStonesCount
			} else {
				blinksCache[blinks-1] = map[string]int64{stone: newStonesCount}
			}
			counter = counter + newStonesCount
		}
		return counter
	} else {
		if cachedBlinks, ok := blinksCache[1]; ok {
			cachedBlinks[data] = int64(len(stones))
		} else {
			blinksCache[1] = map[string]int64{data: int64(len(stones))}
		}
		return int64(len(stones))
	}
}

func part1() int64 {
	var counter int64
	for _, stone := range input {
		counter = counter + mutate(stone, 75)
	}
	return counter
}

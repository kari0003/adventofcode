package main

import (
	"bufio"
	"fmt"
	"os"
)

const INFINITY = 999

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readFile(path string) *bufio.Reader {
	file, err := os.Open(path)
	check(err)
	return bufio.NewReader(file)
}

func readMap(reader *bufio.Reader) [][]int {
	register := [][]int{}
	y := 0
	register = append(register, []int{})
	for {
		c, err := reader.ReadByte()

		if err != nil {
			break
		}
		if c == '\n' {
			y++
			register = append(register, []int{})
			continue
		}
		register[y] = append(register[y], int(c-48))
	}
	return register
}

type Coord struct {
	x int
	y int
}

type HeatLoss struct {
	heatLoss int
	walked   bool
}

type Curly struct {
	east  [3]HeatLoss
	west  [3]HeatLoss
	north [3]HeatLoss
	south [3]HeatLoss
}

type HeatMap struct {
	max         int
	heatMap     [][]int
	heatLossMap [][]int
	walkedMap   [][]bool
}

func (h *HeatMap) init() {
	for x := 0; x < h.max; x++ {
		h.heatLossMap = append(h.heatLossMap, []int{})
		h.walkedMap = append(h.walkedMap, []bool{})
		for y := 0; y < h.max; y++ {
			h.heatLossMap[x] = append(h.heatLossMap[x], INFINITY)
			h.walkedMap[x] = append(h.walkedMap[x], false)
		}
	}
	h.heatLossMap[0][0] = 0
}

func (h *HeatMap) fill(x int, y int, currentHeat int) {
	fmt.Println("fillnig", x, y, currentHeat)
	nextHeat := h.heatMap[x][y] + currentHeat
	if nextHeat < h.heatLossMap[x][y] {
		h.heatLossMap[x][y] = nextHeat
		h.walkedMap[x][y] = false
	}
}

func (h *HeatMap) walkNode(x int, y int) {
	currentHeat := h.heatLossMap[x][y]
	h.walkedMap[x][y] = true
	if x > 0 {
		h.fill(x-1, y, currentHeat)
	}
	if x+1 < h.max {
		h.fill(x+1, y, currentHeat)
	}
	if y > 0 {
		h.fill(x, y-1, currentHeat)
	}
	if y+1 < h.max {

		h.fill(x, y+1, currentHeat)
	}
}

func (h HeatMap) findNextNode() (x int, y int) {
	min := INFINITY
	found := Coord{-1, -1}
	for x := 0; x < h.max; x++ {
		for y := 0; y < h.max; y++ {
			if !h.walkedMap[x][y] && h.heatLossMap[x][y] < min {
				min = h.heatLossMap[x][y]
				found = Coord{x, y}
			}
		}
	}
	return found.x, found.y
}

func travel(heatMap [][]int) int {
	//trail := []Coord{{0,0}}
	h := HeatMap{len(heatMap), heatMap, [][]int{}, [][]bool{}}
	h.init()

	fmt.Println(h.heatLossMap)
	fmt.Println(h.walkedMap)
	// steps := 0
	for {
		// if steps > 100 {
		// 	break
		// }
		// steps++

		x, y := h.findNextNode()
		if x < 0 {
			// No node found
			break
		}
		h.walkNode(x, y)
	}

	fmt.Println(h.heatLossMap)

	return h.heatLossMap[h.max-1][h.max-1]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	heatmap := readMap(reader)

	fmt.Println(heatmap)

	shortestPath := travel(heatmap)

	fmt.Println("v1", shortestPath)
}

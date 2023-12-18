package main

import (
	"bufio"
	"fmt"
	"os"
)

const INFINITY = 178929 // 141*141*9

const EAST = 0
const WEST = 1
const NORTH = 2
const SOUTH = 3

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
	x         int
	y         int
	direction int
	streak    int
}

type HeatLoss struct {
	heatLoss int
	walked   bool
}

type Node [4][]HeatLoss

type HeatMap struct {
	max         int
	minSpeed    int
	maxSpeed    int
	heatMap     [][]int
	heatLossMap [][]Node
}

func (h *HeatMap) init(heatMap [][]int, minSpeed int, maxSpeed int) {
	h.minSpeed = minSpeed
	h.maxSpeed = maxSpeed
	h.heatMap = heatMap
	h.max = len(heatMap)
	for x := 0; x < h.max; x++ {
		h.heatLossMap = append(h.heatLossMap, []Node{})
		for y := 0; y < h.max; y++ {
			node := Node{}
			for d := 0; d < 4; d++ {
				dirHeatLoss := []HeatLoss{}
				for s := h.minSpeed; s <= h.maxSpeed; s++ {
					dirHeatLoss = append(dirHeatLoss, HeatLoss{INFINITY, false})
				}
				node[d] = dirHeatLoss
			}
			h.heatLossMap[x] = append(h.heatLossMap[x], node)
		}
	}
	h.heatLossMap[0][0][2][0] = HeatLoss{0, false}
	h.heatLossMap[0][0][0][0] = HeatLoss{0, false}
}

func (h *HeatMap) fill(x int, y int, direction int, streak int, currentHeat int) {
	nextHeat := h.heatMap[x][y] + currentHeat
	if nextHeat < h.heatLossMap[x][y][direction][streak].heatLoss {
		h.heatLossMap[x][y][direction][streak] = HeatLoss{nextHeat, false}
	}
}

func offset(x int, y int, offset int, direction int) (xout int, yout int) {
	switch direction {
	case EAST:
		xout = x + offset
		yout = y
		break
	case WEST:
		xout = x - offset
		yout = y
		break
	case SOUTH:
		xout = x
		yout = y + offset
		break
	case NORTH:
		xout = x
		yout = y - offset
		break
	}
	return xout, yout
}

func (h HeatMap) inBounds(x int, y int) bool {
	return x >= 0 && y >= 0 && x < h.max && y < h.max
}

func (h *HeatMap) fillDirection(x int, y int, toDirection int, currentHeat int) {
	streakHeat := currentHeat
	for s := 1; s <= h.maxSpeed; s++ {
		toX, toY := offset(x, y, s, toDirection)
		if h.inBounds(toX, toY) {
			if s >= h.minSpeed {
				h.fill(toX, toY, toDirection, s-h.minSpeed, streakHeat)
			}
			streakHeat += h.heatMap[toX][toY]
		}
	}
}

func (h *HeatMap) walkNode(x int, y int, direction int, streak int) {
	currentHeat := h.heatLossMap[x][y][direction][streak].heatLoss
	h.heatLossMap[x][y][direction][streak].walked = true
	if direction > 1 { // can go east or west
		h.fillDirection(x, y, WEST, currentHeat)
		h.fillDirection(x, y, EAST, currentHeat)
	}
	if direction < 2 { // can go north or south
		h.fillDirection(x, y, SOUTH, currentHeat)
		h.fillDirection(x, y, NORTH, currentHeat)
	}
}

func (h HeatMap) findNextNode() (x int, y int, direction int, streak int) {
	min := INFINITY
	found := Coord{-1, -1, 0, 0}
	for x := 0; x < h.max; x++ {
		for y := 0; y < h.max; y++ {
			for d := 0; d < 4; d++ {
				for s := 0; s <= h.maxSpeed-h.minSpeed; s++ {
					if !h.heatLossMap[x][y][d][s].walked && h.heatLossMap[x][y][d][s].heatLoss < min {
						min = h.heatLossMap[x][y][d][s].heatLoss
						found = Coord{x, y, d, s}
					}
				}
			}
		}
	}
	return found.x, found.y, found.direction, found.streak
}

func (n *Node) findMin() (min int, dir int, str int) {
	min = INFINITY
	streakLen := len(n[0])
	for d := 0; d < 4; d++ {
		for s := 0; s < streakLen; s++ {
			if n[d][s].heatLoss < min {
				min = n[d][s].heatLoss
				dir = d
				str = s
			}
		}
	}
	return min, dir, str
}

func (h *HeatMap) travel() [][]Node {
	// steps := 0
	for {
		// if steps > 10 {
		// 	break
		// }
		// steps++

		x, y, d, s := h.findNextNode()
		if x < 0 {
			// No node found
			break
		}
		h.walkNode(x, y, d, s)
	}

	return h.heatLossMap
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	heatmap := readMap(reader)

	h := HeatMap{}

	// h.init(heatmap, 1, 3) // v1
	h.init(heatmap, 4, 10) // v2

	lossMap := h.travel()

	for y := 0; y < len(lossMap); y++ {
		for x := 0; x < len(lossMap); x++ {
			shortest, _, _ := lossMap[x][y].findMin()
			fmt.Print(" ", shortest)
		}
		fmt.Print("\n")
	}

}

/* v1 solution
0 4 5 8 14 17 23 27 29 32 38 40 46
3 5 6 11 15 20 25 30 32 37 43 42 45
6 7 11 16 17 21 26 32 37 41 43 47 49
13 11 15 21 22 29 31 39 41 46 47 52 51
17 16 19 25 28 33 38 46 48 53 52 55 57
18 21 22 30 34 42 46 53 58 61 56 60 63
23 26 27 34 42 49 54 62 66 68 66 66 70
28 32 31 38 46 53 61 70 76 79 75 71 73
33 38 37 41 50 56 64 75 83 86 85 80 80
41 44 45 45 51 58 67 81 89 94 89 89 84
44 46 48 49 55 63 69 79 86 92 94 93 87
49 51 52 55 60 64 72 80 91 99 103 96 98
53 54 54 56 62 69 75 81 86 92 101 99 102
*/

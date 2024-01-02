package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Brick [6]int
type Bricks []Brick
type FallNode struct {
	over  []int
	under []int
}
type FallTree []FallNode

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

func readBricks(reader *bufio.Reader) Bricks {
	bricks := Bricks{}
	for {
		wholeLine, err := reader.ReadString('\n')
		line := strings.TrimRight(wholeLine, "\n")
		ends := strings.Split(line, "~")
		brick := Brick{}
		for i, end := range ends {
			coords := strings.Split(end, ",")
			for j, coord := range coords {
				c, err := strconv.Atoi(coord)
				check(err)
				brick[i*3+j] = c
			}
		}
		bricks = append(bricks, brick)
		if err != nil {
			break
		}
	}
	return bricks
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (brick Brick) maxZ() int {
	return max(brick[2], brick[5])
}
func (brick Brick) minZ() int {
	return min(brick[2], brick[5])
}

func (brick Brick) overlap(with Brick) bool {
	return max(brick[0], with[0]) <= min(brick[3], with[3]) &&
		max(brick[1], with[1]) <= min(brick[4], with[4])
}

func (brick *Brick) setMinZ(z int) {
	diff := brick[2] - z       //6 - 4
	brick[2] = z               // 4
	brick[5] = brick[5] - diff //8 - 2
}

func (bricks Bricks) sortByZ() {
	zCmp := func(a, b Brick) int {
		return min(a[2], a[5]) - min(b[2], b[5])
	}
	slices.SortFunc(bricks, zCmp)
}

func (bricks Bricks) fall() FallTree {
	tree := FallTree{}
	for i := 0; i < len(bricks); i++ {
		tree = append(tree, FallNode{[]int{}, []int{}})
		minZ := 1
		for fallen := 0; fallen < i; fallen++ {
			if bricks[i].overlap(bricks[fallen]) {
				newZ := bricks[fallen].maxZ() + 1
				if minZ <= newZ {
					minZ = newZ
				}
			}
		}
		bricks[i].setMinZ(minZ)
		for fallen := 0; fallen < i; fallen++ {
			if bricks[i].overlap(bricks[fallen]) && minZ == bricks[fallen].maxZ()+1 {
				tree[fallen].under = append(tree[fallen].under, i)
				tree[i].over = append(tree[i].over, fallen)
			}
		}
	}
	return tree
}

func (tree FallTree) canRemove(node int) bool {
	if len(tree[node].under) == 0 {
		fmt.Println("is leaf", node)
		return true
	}
	for _, over := range tree[node].under {
		if len(tree[over].over) == 1 {
			return false
		}
	}
	//fmt.Println("all unders have multiple overs", node)
	return true
}

func (tree FallTree) countCanRemove() int {
	canRemove := 0
	for i := range tree {
		if tree.canRemove(i) {
			canRemove++
		}
	}
	return canRemove
}

func (tree FallTree) sumOver(fallFlags []bool, node int) []bool {
	allHasFallen := true
	for _, i := range tree[node].over {
		allHasFallen = fallFlags[i] && allHasFallen
	}
	if !allHasFallen {
		return fallFlags
	}
	fallFlags[node] = true
	for _, over := range tree[node].under {
		newFlags := tree.sumOver(fallFlags, over)
		fallFlags = newFlags
	}
	//fmt.Println("all unders have multiple overs", node)
	return fallFlags
}

func (tree FallTree) testFallSum() int {
	fallFlags := []bool{}
	for range tree {
		fallFlags = append(fallFlags, false)
	}
	sumSum := 0
	for node := range tree {
		flags := slices.Clone[[]bool, bool](fallFlags)
		flags[node] = true
		for _, over := range tree[node].under {
			newFlags := tree.sumOver(flags, over)
			flags = newFlags
		}
		count := 0
		for _, flag := range flags {
			if flag {
				count++
			}
		}
		fmt.Println("count for", node, count)
		if count > 1 {
			sumSum += count - 1
		}
	}
	return sumSum
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	bricks := readBricks(reader)

	bricks.sortByZ()
	// fmt.Println(bricks)

	deps := bricks.fall()
	fmt.Println(bricks)
	fmt.Println(deps)

	fmt.Println("removable", deps.countCanRemove())

	falls := deps.testFallSum()
	fmt.Println("sum", falls)
}

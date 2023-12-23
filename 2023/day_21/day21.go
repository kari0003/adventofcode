package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Garden [][]byte

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

func readGarden(reader *bufio.Reader) (Garden, error) {
	garden := Garden{}
	newLine := true
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			return garden, err
		}
		if readByte == '\n' {
			if newLine {
				return garden, nil
			}
			newLine = true
		} else {
			if newLine {
				garden = append(garden, []byte{})
				newLine = false
			}
			garden[len(garden)-1] = append(garden[len(garden)-1], readByte)
		}
	}
}

type loopFn func(int, int, byte)

func (garden Garden) loop(callback loopFn) {
	for y, line := range garden {
		for x, c := range line {
			callback(x, y, c)
		}
	}
}

func (garden Garden) lenX() int {
	return len(garden[0])
}
func (garden Garden) lenY() int {
	return len(garden)
}
func (garden Garden) center() int {
	return garden.lenX() / 2
}

func (garden Garden) copy() Garden {
	copy := Garden{}
	garden.loop(func(x int, y int, c byte) {
		if x == 0 {
			copy = append(copy, []byte{})
		}
		copy[y] = append(copy[y], c)
	})
	return copy
}

func (g Garden) stepInto(x, y int) {
	if x >= 0 && y >= 0 && x < g.lenX() && y < g.lenY() && g[y][x] == '.' {
		g[y][x] = '0'
	}
}

func (g Garden) checkReachableFrom(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.lenX() && y < g.lenY() && g[y][x] == '0'
}

func (g Garden) step(into Garden) {
	g.loop(func(x, y int, c byte) {
		if c == '0' {
			into.stepInto(x, y+1)
			into.stepInto(x, y-1)
			into.stepInto(x+1, y)
			into.stepInto(x-1, y)
		}
	})
}

func (g Garden) takeSteps(steps int, empty Garden) Garden {
	even := g.copy()
	odd := empty.copy()
	for i := 0; i < steps; i++ {
		if i%2 == 0 {
			even.step(odd)
		} else {
			odd.step(even)
		}
	}
	if steps%2 == 0 {
		return even
	} else {
		return odd
	}

}

func (g Garden) countBeds() int {
	beds := 0
	g.loop(func(x, y int, c byte) {
		if c == '0' {
			beds++
		}
	})
	return beds
}
func (g Garden) countBedsSeparated() [6]int {
	beds := [6]int{0, 0, 0, 0, 0, 0}
	midX := g.lenX() / 2
	midY := g.lenY() / 2
	g.loop(func(x, y int, c byte) {
		if c == '0' {
			if x == midX {
				beds[4]++
			}
			if y == midX {
				beds[5]++
			}
			if x < midX {
				if y < midY {
					beds[0]++
				}
				if y > midY {
					beds[2]++
				}
			}
			if x > midX {
				if y < midY {
					beds[1]++
				}
				if y > midY {
					beds[3]++
				}
			}
		}
	})
	return beds
}

func (g Garden) print() {
	g.loop(func(x, y int, c byte) {
		fmt.Print(string(c))
		if x == g.lenX()-1 {
			fmt.Print("\n")
		}
	})
}

func walkEdges(edges [4]Garden, steps int, empty Garden) [4]Garden {
	walked := [4]Garden{}
	for i, garden := range edges {
		walked[i] = garden.takeSteps(steps, empty)
	}
	return walked
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("No File or steps argument given!")
		return
	}
	path := os.Args[1]
	steps, err := strconv.Atoi(os.Args[2])
	check(err)
	reader := readFile(path)
	garden, _ := readGarden(reader)

	empty := garden.copy()
	garden.loop(func(x, y int, c byte) {
		if c == 'S' {
			garden[y][x] = '0'
			empty[y][x] = '.'
		}
	})

	if steps > garden.lenX()/2 {
		fullGardenSteps := steps / garden.lenX()
		//remainder := steps - (garden.lenX()/2)%garden.lenX()

		edges := [4]Garden{
			empty.copy(),
			empty.copy(),
			empty.copy(),
			empty.copy(),
		}

		//  .....
		//  .0.1.
		//  .....
		//  .2.3.
		//  .....
		edges[0][0][0] = '0'
		edges[1][0][garden.lenX()-1] = '0'
		edges[2][garden.lenY()-1][0] = '0'
		edges[3][garden.lenY()-1][garden.lenX()-1] = '0'

		corners := [4]Garden{
			empty.copy(),
			empty.copy(),
			empty.copy(),
			empty.copy(),
		}
		corners[0][garden.center()][0] = '0'
		corners[1][garden.center()][garden.lenX()-1] = '0'
		corners[2][0][garden.center()] = '0'
		corners[3][garden.lenY()-1][garden.center()] = '0'

		fullEven := garden.takeSteps(garden.lenX()+1, empty)
		fullOdd := garden.takeSteps(garden.lenX(), empty)
		edgesWalked := walkEdges(edges, garden.center(), empty)
		edgesHighWalked := walkEdges(edges, garden.center()+garden.lenY(), empty)
		//centerWalkedOdd := garden.takeSteps(garden.center()+1, empty)
		cornersWalked := walkEdges(corners, garden.lenX(), empty)

		fullOddCount := (fullGardenSteps - 1) * (fullGardenSteps - 1)
		fullEvenCount := (fullGardenSteps) * (fullGardenSteps)
		edgeHighCount := fullGardenSteps - 1
		edgeLowCount := fullGardenSteps

		fullOdd.print()
		edgesHighWalked[1].print()
		edgesWalked[1].print()

		fmt.Println(" fullOdd", fullOddCount, "fullEven", fullEvenCount, "edgeHigh", edgeHighCount, "edgeLow", edgeLowCount)

		bedsFullOdd := fullOdd.countBeds()
		bedsFullEven := fullEven.countBeds()
		edgeLowBeds := 0
		edgeHighBeds := 0
		cornerBeds := 0
		for i := 0; i < 4; i++ {
			// g.print()
			edgeLowBeds += edgesWalked[i].countBeds()
			cornerBeds += cornersWalked[i].countBeds()
			edgeHighBeds += edgesHighWalked[i].countBeds()
		}

		allBeds := fullEvenCount*bedsFullEven + fullOddCount*bedsFullOdd + edgeHighCount*edgeHighBeds + edgeLowCount*edgeLowBeds + cornerBeds + 4

		fmt.Println("allEdgeBeds", allBeds)
	} else {

		walked := garden.takeSteps(steps, empty)
		beds := walked.countBeds()
		fmt.Println("can walk to ", beds, "beds")
	}

	//walked.print()
}

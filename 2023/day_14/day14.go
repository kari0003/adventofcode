package main

import (
	"bufio"
	"fmt"
	"os"
)

type Parabola [][]byte

var rollCache = map[string]string{}

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

func readParabola(reader *bufio.Reader) (Parabola, error) {
	garden := Parabola{}
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

func (garden Parabola) loop(callback loopFn) {
	for y, line := range garden {
		for x, c := range line {
			callback(x, y, c)
		}
	}
}

func (garden Parabola) lenX() int {
	return len(garden[0])
}
func (garden Parabola) lenY() int {
	return len(garden)
}
func (garden Parabola) center() int {
	return garden.lenX() / 2
}

func (garden Parabola) copy() Parabola {
	copy := Parabola{}
	garden.loop(func(x int, y int, c byte) {
		if x == 0 {
			copy = append(copy, []byte{})
		}
		copy[y] = append(copy[y], c)
	})
	return copy
}

func (g Parabola) print() {
	g.loop(func(x, y int, c byte) {
		fmt.Print(string(c))
		if x == g.lenX()-1 {
			fmt.Print("\n")
		}
	})
}

func (p Parabola) rollNorth() {
	for x := 0; x < p.lenX(); x++ {
		for y := 0; y < p.lenY(); y++ {
			if p[y][x] == 'O' {
				toY := y
				for {
					if toY-1 >= 0 && p[toY-1][x] == '.' {
						toY -= 1
					} else {
						break
					}
				}
				p[y][x] = '.'
				p[toY][x] = 'O'
			}
		}
	}
}
func (p Parabola) rollSouth() {
	for x := 0; x < p.lenX(); x++ {
		for y := p.lenY() - 1; y >= 0; y-- {
			if p[y][x] == 'O' {
				toY := y
				for {
					if toY+1 < p.lenY() && p[toY+1][x] == '.' {
						toY += 1
					} else {
						break
					}
				}
				p[y][x] = '.'
				p[toY][x] = 'O'
			}
		}
	}
}

func (p Parabola) rollWest() {
	for y := 0; y < p.lenY(); y++ {
		for x := 0; x < p.lenX(); x++ {
			if p[y][x] == 'O' {
				toX := x
				for {
					if toX-1 >= 0 && p[y][toX-1] == '.' {
						toX -= 1
					} else {
						break
					}
				}
				p[y][x] = '.'
				p[y][toX] = 'O'
			}
		}
	}
}
func (p Parabola) rollEast() {
	for y := 0; y < p.lenY(); y++ {
		for x := p.lenX() - 1; x >= 0; x-- {
			if p[y][x] == 'O' {
				toX := x
				for {
					if toX+1 < p.lenX() && p[y][toX+1] == '.' {
						toX += 1
					} else {
						break
					}
				}
				p[y][x] = '.'
				p[y][toX] = 'O'
			}
		}
	}
}
func (p Parabola) weighNorth() int {
	weight := 0
	p.loop(func(x, y int, c byte) {
		if c == 'O' {
			weight += (p.lenY() - y)
		}
	})
	return weight
}

func (parabola Parabola) cycle() {
	parabola.rollNorth()
	parabola.rollWest()
	parabola.rollSouth()
	parabola.rollEast()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	parabola, _ := readParabola(reader)

	v1 := parabola.copy()
	v1.rollNorth()
	fmt.Println(v1.weighNorth())

	v2 := parabola.copy()

	for i := 0; i < 3; i++ {
		v2.cycle()
		// if i%10000000 == 0 {
		// 	v2.print()
		// }
	}

	v2.print()
	fmt.Println(v2.weighNorth())
}

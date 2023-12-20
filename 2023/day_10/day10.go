package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const START = 'S'

type Pipe byte
type PipeMap [][]Pipe
type Mask [][]bool

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

func readMap(reader *bufio.Reader) PipeMap {
	register := PipeMap{}
	y := 0
	register = append(register, []Pipe{})
	for {
		c, err := reader.ReadByte()

		if err != nil {
			break
		}
		if c == '\n' {
			y++
			register = append(register, []Pipe{})
			continue
		}
		register[y] = append(register[y], Pipe(c))
	}
	return register
}

func (p PipeMap) findStart() (x int, y int, err error) {
	for y, line := range p {
		for x, c := range line {
			if c == START {
				return x, y, nil
			}
		}
	}
	return 0, 0, errors.New("Start not found")
}

func (p PipeMap) initMask() Mask {
	mask := [][]bool{}
	for y, line := range p {
		mask = append(mask, []bool{})
		for range line {
			mask[y] = append(mask[y], false)
		}
	}
	return mask
}

// Dir
// 0 - east
// 1 - west
// 2 - north
// 3 - south
func (p PipeMap) goNext(inX int, inY int, inDir int) (outX int, outY int, outDir int) {
	c := p[inY][inX]
	switch c {
	case '-':
		outDir = inDir
	case '|':
		outDir = inDir
	case 'L':
		if inDir == 3 {
			outDir = 1
		} else if inDir == 0 {
			outDir = 2
		}
	case 'J':
		if inDir == 3 {
			outDir = 0
		} else if inDir == 1 {
			outDir = 2
		}
	case 'F':
		if inDir == 2 {
			outDir = 1
		} else if inDir == 0 {
			outDir = 3
		}
	case '7':
		if inDir == 2 {
			outDir = 0
		} else if inDir == 1 {
			outDir = 3
		}
	}
	switch outDir {
	case 0:
		outX = inX - 1
		outY = inY
		break
	case 1:
		outX = inX + 1
		outY = inY
		break
	case 2:
		outX = inX
		outY = inY - 1
		break
	case 3:
		outX = inX
		outY = inY + 1
		break
	}

	return outX, outY, outDir
}

func (p PipeMap) travel() (len int, mask Mask) {
	mask = p.initMask()
	x, y, err := p.findStart()
	check(err)
	len = 1
	x += 1
	dir := 1
	for {
		mask[y][x] = true
		x, y, dir = p.goNext(x, y, dir)
		len += 1
		if p[y][x] == START {
			mask[y][x] = true
			break
		}
	}
	return len, mask
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	pipes := readMap(reader)

	len, mask := pipes.travel()
	fmt.Println("pipe deepness", len/2)

	all := 0
	for x, line := range mask {
		countLine := 0
		crosses := 0
		lastEntry := Pipe('|')
		crossing := false
		for y, flag := range line {
			if flag {
				current := pipes[x][y]

				if current == START {
					current = '-'
				}
				if current == '|' {
					crosses++
					lastEntry = current
				} else if current == 'L' || current == 'F' {
					lastEntry = current
					crossing = true
				} else if current == 'J' {
					if lastEntry == 'F' {
						crosses++
					}
					crossing = false
					lastEntry = '|'
				} else if current == '7' {
					if lastEntry == 'L' {
						crosses++
					}
					crossing = false
					lastEntry = '|'
				}

				fmt.Print(string(pipes[x][y]))
			} else {
				if !crossing && crosses%2 == 1 {
					countLine++

					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			//			fmt.Print(crosses)
		}
		all += countLine
		fmt.Print(" ", countLine, " \n")
	}

	fmt.Println("All", all)
}

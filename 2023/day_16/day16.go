package main

import (
	"bufio"
	"fmt"
	"os"
)

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

func processMap(reader *bufio.Reader) []string {
	space := []string{}
	for {
		line, err := reader.ReadString('\n')
		space = append(space, line)
		if err != nil {
			break
		}
	}
	return space
}

func createRegister(space []string) [][]int {
	register := [][]int{}
	maxX := len(space)
	maxY := len(space[0]) - 1

	for x := 0; x < maxX; x++ {
		register = append(register, []int{})
		for y := 0; y < maxY; y++ {
			register[x] = append(register[x], 0)
		}
	}
	return register
}

type Light struct {
	x  int
	y  int
	vX int
	vY int
}

func dirToReg(x int, y int) int {
	if x < 0 {
		return 1
	} else if x > 0 {
		return 2
	} else if y < 0 {
		return 4
	} else {
		return 8
	}
}

func bounceLights(space []string, l Light) [][]int {
	lights := []Light{l}
	register := createRegister(space)

	maxX := len(register)
	maxY := len(register[0])

	for {
		if len(lights) == 0 {
			break
		}
		newLights := []Light{}
		for _, l := range lights {
			if l.x < 0 || l.x >= maxX || l.y < 0 || l.y >= maxY {
				continue
			}
			dir := dirToReg(l.vX, l.vY)
			if register[l.x][l.y]|dir == register[l.x][l.y] {
				continue
			}
			register[l.x][l.y] |= dir
			switch space[l.x][l.y] {
			case '.':
				l.x += l.vX
				l.y += l.vY
			case '\\':
				d := l.vX
				l.vX = l.vY
				l.vY = d

				l.x += l.vX
				l.y += l.vY
			case '/':
				d := l.vX
				l.vX = -l.vY
				l.vY = -d

				l.x += l.vX
				l.y += l.vY
			case '-':
				if l.vX == 0 {
					l.x += l.vX
					l.y += l.vY
				} else {
					newLights = append(newLights, Light{l.x, l.y - 1, 0, -1})
					l.vX = 0
					l.vY = 1
					l.x += l.vX
					l.y += l.vY
				}
			case '|':
				if l.vY == 0 {
					l.x += l.vX
					l.y += l.vY
				} else {
					newLights = append(newLights, Light{l.x - 1, l.y, -1, 0})
					l.vX = 1
					l.vY = 0
					l.x += l.vX
					l.y += l.vY
				}
			}
			newLights = append(newLights, l)
		}
		// fmt.Println(newLights)
		lights = newLights
	}
	return register
}

func countLights(register [][]int) int {
	sum := 0
	for _, r := range register {
		for _, v := range r {
			if v > 0 {
				sum++
			}
		}
	}
	return sum
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	space := processMap(reader)
	v1Reg := bounceLights(space, Light{0, 0, 0, 1})
	for _, r := range v1Reg {
		for _, v := range r {
			if v > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	v1 := countLights(v1Reg)
	fmt.Println("v1 solution", v1)

	cMax := v1
	max := len(space)
	for i := 0; i < max; i++ {
		reg := bounceLights(space, Light{i, 0, 0, 1})
		count := countLights(reg)
		if count > cMax {
			cMax = count
		}

		reg = bounceLights(space, Light{i, max - 1, 0, -1})
		count = countLights(reg)
		if count > cMax {
			cMax = count
		}

		reg = bounceLights(space, Light{0, i, 1, 0})
		count = countLights(reg)
		if count > cMax {
			cMax = count
		}

		reg = bounceLights(space, Light{max - 1, i, -1, 0})
		count = countLights(reg)
		if count > cMax {
			cMax = count
		}
	}
	fmt.Println("v2 solution", cMax)
}

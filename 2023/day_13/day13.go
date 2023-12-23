package main

import (
	"bufio"
	"fmt"
	"os"
)

type Reading struct {
	len int
	r   [][]int
}

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

func readMap(reader *bufio.Reader) (Reading, error) {
	reading := Reading{
		len: 0,
		r:   [][]int{},
	}
	newLine := true
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			return reading, err
		}
		if readByte == '\n' {
			if newLine {
				return reading, nil
			}
			newLine = true
		} else {
			if newLine {
				reading.r = append(reading.r, []int{})
				newLine = false
			}
			if len(reading.r) == 1 {
				reading.len++
			}
			currentReading := 0
			if readByte == '#' {
				currentReading = 1
			}
			reading.r[len(reading.r)-1] = append(reading.r[len(reading.r)-1], currentReading)
		}
	}
}

func (reading Reading) checkVerticalMirror() (int, bool) {
	mirrored := [][]int{}
	for range reading.r {
		mirrored = append(mirrored, []int{})
	}
	for i := 0; i < reading.len; i++ {
		matches := true
		for y, r := range reading.r {
			mirrored[y] = append(mirrored[y], r[i])
			if mirrored[y] == r[] {

			}
		}
		if matches {
			return i + 1, true
		}
	}
	return 0, false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)

	readings := []Reading{}
	for {
		reading, err := readMap(reader)
		readings = append(readings, reading)

		if err != nil {
			break
		}
	}

	//readings[0].checkVerticalMirror()

	for _, reading := range readings {
		for _, line := range reading.r {
			for _, c := range line {
				if c > 0 {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Print("\n")
		}
		fmt.Print("\n")
	}
}

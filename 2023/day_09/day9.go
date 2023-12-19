package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func readLines(reader *bufio.Reader) [][]int {
	lines := [][]int{}
	for {
		line, err := reader.ReadString('\n')
		newLine := []int{}
		if line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		numStrings := strings.Split(line, " ")
		for _, num := range numStrings {
			parsed, err := strconv.ParseInt(num, 10, 64)
			check(err)
			newLine = append(newLine, int(parsed))
		}
		lines = append(lines, newLine)

		if err != nil {
			break
		}
	}
	return lines
}

type Readings []int

func (r Readings) differentiate(depth int) Readings {
	diff := []int{}
	for i := 0; i < len(r); i++ {
		if i <= depth {
			diff = append(diff, 0)
		} else {
			diff = append(diff, r[i]-r[i-1])
		}
	}
	return diff
}

func (r Readings) isItAllZero() bool {
	for _, diff := range r {
		if diff != 0 {
			return false
		}
	}
	return true
}

func solve(readings Readings, plus int, depth int) int {
	diff := readings.differentiate(depth)
	// fmt.Println("started from the bottom now we are:", depth, diff)
	// if depth > 10 {
	// 	return plus
	// }
	if diff.isItAllZero() {
		return plus + readings[len(readings)-1]
	}
	return solve(diff, plus+readings[len(readings)-1], depth+1)
}

func solveReversed(r Readings) int {
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return solve(r, 0, 0)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	lines := readLines(reader)

	solved := 0
	for _, l := range lines {
		solved += solve(Readings(l), 0, 0)
	}

	fmt.Println("solved", solved)

	solved = 0
	for _, l := range lines {
		solved += solveReversed(Readings(l))
	}

	fmt.Println("solve reversed", solved)
}

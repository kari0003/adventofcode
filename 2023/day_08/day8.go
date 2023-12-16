package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	R string
	L string
}

const instructions = "RL"
const instructions2 = "LLR"
const realInstructions = "LRRRLRLLLLLLLRLRLRRLRRRLRRLRRRLRRLRRRLLRRRLRRLRLRRRLRRLRRRLLRLLRRRLRRRLRLLRLRLRRRLRRLRRLRRLRLRRRLRRLRRRLLRLLRLLRRLRLLRLRRLRLRLRRLRRRLLLRRLRRRLLRRLRLRLRRRLRLRRRLLRLLLRRRLLLRRLLRLLRRLLRLRRRLRLRRLRRLLRRLRLLRLRRRLRRRLRLRRRLRLRLRRLRLRRRLRRRLRRRLRRLRRLRRRLLRLRLLRLLRRRR"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(path string) *bufio.Reader {
	file, err := os.Open(path)
	check(err)
	return bufio.NewReader(file)
}

func parseMap(reader *bufio.Reader) map[string]Node {
	nodes := make(map[string]Node)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		nodes[line[:3]] = Node{
			L: line[7:10],
			R: line[12:15],
		}
	}
	return nodes
}

func navigateMap(nodes map[string]Node, instructions string) int {
	allInst := len(instructions)
	steps := 0
	currentNode := nodes["AAA"]
	for {
		currentOperation := instructions[steps%allInst]
		// fmt.Println(instructions, steps%allInst, steps)
		steps++
		var next string
		if currentOperation == 76 {
			next = currentNode.L
		} else {
			next = currentNode.R
		}
		// fmt.Println(currentOperation, next, next == "ZZZ")
		if next == "ZZZ" {
			break
		}
		currentNode = nodes[next]
	}
	return steps
}

func main() {
	reader := readFile("seed/8nodes")
	nodes := parseMap(reader)
	fmt.Println(nodes)

	steps := navigateMap(nodes, realInstructions)
	fmt.Println(steps)

}

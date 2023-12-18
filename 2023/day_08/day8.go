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
const instructionsv2 = "LR"

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
		nodes[line[:3]] = Node{
			L: line[7:10],
			R: line[12:15],
		}
		if err != nil {
			break
		}
	}
	return nodes
}

/*
AAA 15517
DNA 19199
SHA 12361
DLA 20777
JVA 13939
XLA 17621

# Finding the least common multiple

15517 19199 12361 20777 13939 17621
*/

func getStartNodes(nodes map[string]Node) []string {
	startNodes := []string{}
	for name := range nodes {
		if name[2] == 'A' {
			startNodes = append(startNodes, name)
		}
	}
	return startNodes
}

func checkEnd(name string) bool {
	return name[2] == 'Z'
}

func navigateMap(nodes map[string]Node, instructions string, nodeName string) int {
	allInst := len(instructions)
	steps := 0
	currentNode := nodeName
	for {
		// if steps > 10 {
		// 	break
		// }
		currentOperation := instructions[steps%allInst]
		// fmt.Println(instructions, steps%allInst, steps)
		steps++

		if currentOperation == 'L' {
			currentNode = nodes[currentNode].L
		} else {
			currentNode = nodes[currentNode].R
		}

		if checkEnd(currentNode) {
			break
		}
	}
	return steps
}

func allEqual(multiples []int) bool {
	for _, value := range multiples {
		if value != multiples[0] {
			return false
		}
	}
	return true
}

func leastCommonMultiple(steps []int) int {
	lcm := steps[0]
	for i := 1; i < len(steps); i++ {
		gcd := greatestCommonDenominator(lcm, steps[i])
		lcm = (steps[i] / gcd) * lcm
	}
	return lcm
}

func greatestCommonDenominator(a int, b int) int {
	var bigger int
	var smaller int
	if a > b {
		bigger = a
		smaller = b
	} else {
		smaller = a
		bigger = b
	}
	next := bigger % smaller
	if next == 0 {
		return b
	} else {
		return greatestCommonDenominator(smaller, next)
	}
}

func v1(nodes map[string]Node) {
	steps := navigateMap(nodes, realInstructions, "AAA")
	fmt.Println(steps)
}

func main() {
	reader := readFile("seed/8nodes")
	nodes := parseMap(reader)

	// v1(nodes)

	startNodes := getStartNodes(nodes)

	steps := []int{}
	for _, nodeName := range startNodes {
		step := navigateMap(nodes, realInstructions, nodeName)
		steps = append(steps, step)
	}
	fmt.Println("finishers", steps)

	fmt.Println("leastCommon", leastCommonMultiple(steps))
}

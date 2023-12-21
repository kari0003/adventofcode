package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const LOW = false
const HIGH = true

type Pulse struct {
	from  string
	to    string
	value bool
}

type Machine struct {
	pulses      *list.List
	rx          bool
	pulsesHigh  int
	pulsesLow   int
	flips       map[string]bool
	conjunctors map[string]map[string]bool
	outputs     map[string][]string
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

func (m *Machine) send(from string, val bool) {
	for _, output := range m.outputs[from] {
		m.pulses.PushBack(Pulse{from, output, val})
	}
}

func (m *Machine) run() {
	for m.pulses.Len() > 0 {

		first := m.pulses.Front()
		m.pulses.Remove(first)
		pulse := first.Value.(Pulse)
		//fmt.Println(pulse.from, "-", pulse.value, "->", pulse.to)

		if pulse.to == "rx" && pulse.value == LOW {
			m.rx = true
			fmt.Println("found RX", m.pulsesHigh+m.pulsesLow)
			break
		}

		if pulse.value {
			m.pulsesHigh++
		} else {
			m.pulsesLow++
		}
		flip, isFlip := m.flips[pulse.to]
		conjunctor, isConjunctor := m.conjunctors[pulse.to]

		if isFlip {
			if pulse.value == LOW {
				m.flips[pulse.to] = !flip
				m.send(pulse.to, m.flips[pulse.to])
			}
		} else if isConjunctor {
			//fmt.Println("Conj", pulse.to, conjunctor)
			conjunctor[pulse.from] = pulse.value
			allHigh := true
			for _, val := range conjunctor {
				allHigh = allHigh && val
			}
			m.send(pulse.to, !allHigh)
		} else {
			m.send(pulse.to, pulse.value)
		}
	}
	//fmt.Println("all pulses done", "high:", m.pulsesHigh, "low:", m.pulsesLow)
}

var secondPartSeparator = regexp.MustCompile(` -> `)
var outputSeparator = regexp.MustCompile(`, `)

func parseMachine(reader *bufio.Reader) Machine {
	machine := Machine{
		rx:          false,
		pulses:      list.New(),
		pulsesHigh:  0,
		pulsesLow:   0,
		flips:       map[string]bool{},
		conjunctors: map[string]map[string]bool{},
		outputs:     map[string][]string{},
	}

	for {
		read, err := reader.ReadString('\n')
		line := strings.TrimSuffix(read, "\n")
		outputsString := secondPartSeparator.Split(line, 2)
		outputs := outputSeparator.Split(outputsString[1], -1)
		switch line[0] {
		case '%':
			name := outputsString[0][1:]
			machine.flips[name] = false
			machine.outputs[name] = outputs
			// flipper
		case '&':
			name := outputsString[0][1:]
			machine.conjunctors[name] = map[string]bool{}
			machine.outputs[name] = outputs
			// conjunction
		default:
			name := outputsString[0]
			machine.outputs[name] = outputs
		}
		if err != nil {
			break
		}
	}

	for from, outputs := range machine.outputs {
		for _, to := range outputs {
			conj, isConj := machine.conjunctors[to]
			if isConj {
				conj[from] = false
			}
		}
	}
	return machine
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	machine := parseMachine(reader)
	fmt.Println(machine)

	i := 0
	for !machine.rx {
		machine.pulses.PushBack(Pulse{"button", "broadcaster", false})

		machine.run()

		i++
	}

	fmt.Println("solution after ", i, ":", machine.pulsesHigh*machine.pulsesLow)
}

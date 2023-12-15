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

func hash(reader *bufio.Reader) (byte, error) {
	value := uint8(0)
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return value, err
		}
		if char == 44 {
			return value, nil
		} else {
			value += char
			value *= 17
		}
	}
}

func readOperation(reader *bufio.Reader) (byte, string, bool, byte, error) {
	hash := uint8(0)
	var label string
	isAdd := false
	labelFinished := false
	focalLength := uint8(0)
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return hash, label, isAdd, focalLength, err
		}
		if char == 45 { // -
			labelFinished = true
			isAdd = false
		}
		if char == 61 { // =
			labelFinished = true
			isAdd = true
		}
		if char == 44 { // ,
			return hash, label, isAdd, focalLength, nil
		} else {
			if labelFinished {
				focalLength = char - 48
			} else {
				label += string(char)
				hash += char
				hash *= 17
			}
		}
	}
}

func processV1(reader *bufio.Reader) {
	all := int64(0)
	for {
		currrentValue, err := hash(reader)
		all += int64(currrentValue)

		if err != nil {
			break
		}
	}
	fmt.Println("finished", all)
}

type Lense struct {
	label       string
	focalLength byte
}

func processV2(reader *bufio.Reader) {
	boxes := [256][]Lense{}
	for i := range boxes {
		boxes[i] = []Lense{}
	}
	for {
		hash, label, isAdd, focalLength, err := readOperation(reader)
		if isAdd {
			fmt.Println("Add", label, hash, focalLength)
			replaced := false
			for i, lense := range boxes[hash] {
				if lense.label == label {
					boxes[hash][i].focalLength = focalLength
					replaced = true
				}
			}
			if !replaced {
				var a Lense
				a.label = label
				a.focalLength = focalLength
				boxes[hash] = append(boxes[hash], a)
			}
		} else {
			fmt.Println("Remove", label)
			removeIndex := -1
			for i, lense := range boxes[hash] {
				if lense.label == label {
					removeIndex = i
				}
			}
			if removeIndex >= 0 {
				boxes[hash] = append(boxes[hash][:removeIndex], boxes[hash][removeIndex+1:]...)
			}
		}

		if err != nil {
			break
		}
	}
	focalPower := 0
	for hash, box := range boxes {
		for i, lense := range box {
			currentPower := hash + 1
			currentPower *= i + 1
			currentPower *= int(lense.focalLength)
			focalPower += currentPower

		}
	}
	fmt.Println("finished", focalPower, boxes)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File argument given!")
		return
	}
	path := os.Args[1]
	reader := readFile(path)
	processV1(reader)
	reader2 := readFile(path)
	processV2(reader2)
}

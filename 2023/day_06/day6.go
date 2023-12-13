package main

import (
	"fmt"
	"math"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getInput(test bool, v2 bool) ([4]int, [4]int) {
	if test {
		if v2 {
			times := [4]int{71530, 0, 0, 0}
			dists := [4]int{940200, 0, 0, 0}
			return times, dists
		}
		times := [4]int{7, 15, 30, 0}
		dists := [4]int{9, 40, 200, 0}
		return times, dists
	}
	if v2 {

		times := [4]int{48938595, 0, 0, 0}
		dists := [4]int{296192812361391, 0, 0, 0}
		return times, dists
	}
	times := [4]int{48, 93, 85, 95}
	dists := [4]int{296, 1928, 1236, 1391}
	return times, dists
}

// x * (7-x) = 9
// -x2 + 7x - 9 = 0

// b +- sqrt(b2 + 4c)/2

func solutions(time int, dist int) (int, int) {
	det := math.Sqrt(float64(time*time - 4*dist))
	fmt.Println(time, dist, det)
	x1 := int(math.Floor((float64(time)+det)/2 - 0.0001))
	x2 := int(math.Ceil((float64(time)-det)/2 + 0.0001))
	return x1, x2
}

func main() {
	pow := 1
	times, dists := getInput(false, true)
	for i, time := range times {
		x1, x2 := solutions(time, dists[i])
		fmt.Println("solution", x1, x2, x1-x2+1)
		pow *= x1 - x2 + 1
	}
	fmt.Println(pow)
}

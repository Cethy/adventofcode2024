package main

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	raw, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	var list1, list2 []int
	// make lists
	lines := strings.Split(string(raw), "\n")
	for it, line := range lines {
		items := strings.Split(line, "   ")
		if len(items) != 2 {
			panic("something's wrong on line " + strconv.Itoa(it))
		}

		item1, err := strconv.Atoi(strings.Trim(items[0], " "))
		if err != nil {
			panic(err)
		}
		list1 = append(list1, item1)

		item2, err := strconv.Atoi(strings.Trim(items[1], " "))
		if err != nil {
			panic(err)
		}
		list2 = append(list2, item2)
	}

	// order lists
	sort.Ints(list1)
	sort.Ints(list2)

	distances := make([]int, len(list1))
	// compute distances
	for i := 0; i < len(list1); i++ {
		d := list1[i] - list2[i]
		if d < 0 {
			d = -d
		}
		distances[i] = d
	}

	// return sum of all distances
	var distSum int
	for _, d := range distances {
		distSum += d
	}

	log.Println("sum of all distances: ", distSum)
}

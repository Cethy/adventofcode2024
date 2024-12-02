package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func removeItem(report []int, position int) []int {
	ret := make([]int, 0)
	ret = append(ret, report[:position]...)
	ret = append(ret, report[position+1:]...)
	return ret
}

// generateCombinations from s by removing once every possible item in the list
func generateCombinations(s []int) [][]int {
	combinations := make([][]int, 0)
	for i := 0; i < len(s); i++ {
		combinations = append(combinations, removeItem(s, i))
	}
	return combinations
}

func main() {
	raw, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}

	safeCount := 0
	for _, line := range strings.Split(string(raw), "\n") {
		lineExpl := strings.Split(line, " ")

		report := make([]int, len(lineExpl))
		for i, e := range lineExpl {
			report[i], _ = strconv.Atoi(e)
		}
		log.Println("original:", line)
		if isSafe(report) {
			log.Println("safe")
			safeCount++
		} else {
			log.Println("check combinations:")
			for _, combination := range generateCombinations(report) {
				if isSafe(combination) {
					log.Println(combination, "safe")
					safeCount++
					break
				}
			}
		}
		fmt.Println("")
	}

	log.Printf("Safe: %d\n", safeCount)
}

func isSafe(report []int) bool {
	if len(report) < 2 {
		return true
	}

	direction := ""
	if report[0] < report[1] {
		direction = "incr"
	} else if report[0] > report[1] {
		direction = "decr"
	}
	if direction == "" {
		return false
	}

	prevVal := report[0]
	for i := 1; i < len(report); i++ {
		diff := math.Abs(float64(prevVal - report[i]))

		if diff < 1 || diff > 3 {
			return false
		}
		if direction == "incr" && prevVal > report[i] {
			return false
		}
		if direction == "decr" && prevVal < report[i] {
			return false
		}

		prevVal = report[i]
	}
	return true
}

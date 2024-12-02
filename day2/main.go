package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

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

		if isSafe(report) {
			log.Println(line)
			safeCount++
		}
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
		// 2 first number are identical
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
